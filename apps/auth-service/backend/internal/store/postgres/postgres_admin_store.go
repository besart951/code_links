package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func (s *Store) GetAdminActor(ctx context.Context, userID uuid.UUID) (AdminActor, error) {
	user, _, err := s.FindUserByID(ctx, userID)
	if err != nil {
		return AdminActor{}, err
	}

	roles, permissions, err := s.ListUserRoles(ctx, userID)
	if err != nil {
		return AdminActor{}, err
	}
	if len(permissions) == 0 {
		return AdminActor{}, errForbidden
	}

	return AdminActor{ID: user.ID.String(), Email: user.Email, Name: user.Name, Roles: roles, Permissions: permissions}, nil
}

func (s *Store) ListUserRoles(ctx context.Context, userID uuid.UUID) ([]AdminRole, []AdminPermission, error) {
	rows, err := s.pool.Query(ctx, `
		select role_id
		from user_roles
		where user_id = $1
		order by role_id
	`, userID)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	roles := []AdminRole{}
	for rows.Next() {
		var role AdminRole
		if err := rows.Scan(&role); err != nil {
			return nil, nil, err
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return roles, permissionsForRoles(roles), nil
}

func (s *Store) ListAdminUsers(ctx context.Context, query AdminUserListQuery) (AdminUserListResult, error) {
	page, pageSize := normalizedPage(query.Page, query.PageSize)
	offset := (page - 1) * pageSize
	orderBy := adminUserListOrderBy(query.Sort, query.Direction)
	rows, err := s.pool.Query(ctx, `
		with role_summary as (
			select user_id,
			       bool_or(role_id = 'admin') as is_admin,
			       bool_or(role_id = 'support') as is_support,
			       bool_or(role_id = 'auditor') as is_auditor
			from user_roles
			group by user_id
		),
		login_summary as (
			select user_id,
			       count(*) filter (where success) as successful_login_count,
			       count(*) filter (where not success) as failed_login_count
			from login_attempts
			where user_id is not null
			group by user_id
		),
		filtered as (
			select u.id,
			       u.email,
			       u.name,
			       u.status,
			       u.email_verified_at,
			       u.created_at,
			       u.last_login_at,
			       coalesce(u.last_login_ip::text, '') as last_login_ip,
			       coalesce(u.last_login_country_code::text, '') as last_login_country_code,
			       case
			       when coalesce(rs.is_admin, false) then 'admin'
			       when coalesce(rs.is_support, false) then 'support'
			       when coalesce(rs.is_auditor, false) then 'auditor'
			       else 'user'
			       end as primary_role,
			       coalesce(ls.successful_login_count, 0) as successful_login_count,
			       coalesce(ls.failed_login_count, 0) as failed_login_count
			from users u
			left join role_summary rs on rs.user_id = u.id
			left join login_summary ls on ls.user_id = u.id
			where ($1 = '' or lower(u.name) like '%' || lower($1) || '%' or lower(u.email) like '%' || lower($1) || '%')
			  and ($2 = '' or u.status = $2)
		),
		counted as (
			select *, count(*) over() as total_count
			from filtered
			where ($3 = '' or primary_role = $3)
		)
		select id,
		       email,
		       name,
		       status,
		       email_verified_at,
		       created_at,
		       last_login_at,
		       last_login_ip,
		       last_login_country_code,
		       primary_role,
		       successful_login_count,
		       failed_login_count,
		       total_count
		from counted
		order by `+orderBy+`
		limit $4 offset $5
	`, query.Query, string(query.Status), query.Role, pageSize, offset)
	if err != nil {
		return AdminUserListResult{}, err
	}
	defer rows.Close()

	items := []AdminUserListItem{}
	total := 0
	for rows.Next() {
		var id uuid.UUID
		var email string
		var name string
		var status UserStatus
		var emailVerifiedAt sql.NullTime
		var createdAt time.Time
		var lastLoginAt sql.NullTime
		var lastLoginIP string
		var lastLoginCountryCode string
		var primaryRole AdminRole
		var successfulLoginCount int
		var failedLoginCount int
		if err := rows.Scan(
			&id,
			&email,
			&name,
			&status,
			&emailVerifiedAt,
			&createdAt,
			&lastLoginAt,
			&lastLoginIP,
			&lastLoginCountryCode,
			&primaryRole,
			&successfulLoginCount,
			&failedLoginCount,
			&total,
		); err != nil {
			return AdminUserListResult{}, err
		}
		var lastLogin *time.Time
		if lastLoginAt.Valid {
			lastLogin = &lastLoginAt.Time
		}
		items = append(items, AdminUserListItem{
			ID:                   id.String(),
			Name:                 name,
			Email:                email,
			PrimaryRole:          primaryRole,
			Status:               status,
			EmailVerified:        emailVerifiedAt.Valid,
			CreatedAt:            createdAt,
			LastLoginAt:          lastLogin,
			SuccessfulLoginCount: successfulLoginCount,
			FailedLoginCount:     failedLoginCount,
			LastKnownIPAddress:   emptyStringToNil(lastLoginIP),
			LastLoginCountryCode: emptyStringToNil(lastLoginCountryCode),
		})
	}
	if err := rows.Err(); err != nil {
		return AdminUserListResult{}, err
	}

	return AdminUserListResult{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func adminUserListOrderBy(field string, direction string) string {
	column := "created_at"
	switch field {
	case "name":
		column = "lower(name)"
	case "email":
		column = "lower(email)"
	case "primaryRole":
		column = "primary_role"
	case "status":
		column = "status"
	case "lastLoginAt":
		column = "last_login_at"
	case "createdAt", "":
		column = "created_at"
	}

	dir := "desc"
	if direction == "asc" {
		dir = "asc"
	}

	return column + " " + dir + " nulls last, created_at desc, id asc"
}

func (s *Store) GetManagedUserDetail(ctx context.Context, userID uuid.UUID) (ManagedUserDetail, error) {
	user, _, err := s.FindUserByID(ctx, userID)
	if err != nil {
		return ManagedUserDetail{}, err
	}
	item, err := s.adminUserListItem(ctx, user)
	if err != nil {
		return ManagedUserDetail{}, err
	}
	licenses, err := s.listLicenses(ctx, userID)
	if err != nil {
		return ManagedUserDetail{}, err
	}
	roles, err := s.userPermissionGrants(ctx, userID)
	if err != nil {
		return ManagedUserDetail{}, err
	}
	knownIPs, countries, devices, err := s.userLoginFacets(ctx, userID)
	if err != nil {
		return ManagedUserDetail{}, err
	}

	return ManagedUserDetail{
		AdminUserListItem: item,
		Roles:             roles,
		ProductLicenses:   licenses,
		KnownIPAddresses:  knownIPs,
		LoginCountries:    countries,
		UsedDevices:       devices,
	}, nil
}

func (s *Store) SetUserStatus(ctx context.Context, userID uuid.UUID, status UserStatus) (User, error) {
	if !validUserStatus(status) {
		return User{}, errNotFound
	}
	tag, err := s.pool.Exec(ctx, `
		update users
		set status = $2, updated_at = now()
		where id = $1
	`, userID, status)
	if err != nil {
		return User{}, err
	}
	if tag.RowsAffected() == 0 {
		return User{}, errNotFound
	}

	return s.scanUser(ctx, `
		select id, email, name, password_hash, status, email_verified_at, created_at, updated_at,
		       last_login_at, coalesce(last_login_ip::text, ''), coalesce(last_login_country_code::text, '')
		from users
		where id = $1
	`, userID)
}

func (s *Store) SetUserRole(ctx context.Context, userID uuid.UUID, role AdminRole) (User, error) {
	if !validAdminRole(role) {
		return User{}, errNotFound
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `delete from user_roles where user_id = $1`, userID); err != nil {
		return User{}, err
	}
	if _, err := tx.Exec(ctx, `insert into user_roles (user_id, role_id) values ($1, $2)`, userID, role); err != nil {
		return User{}, err
	}
	if _, err := tx.Exec(ctx, `update users set updated_at = now() where id = $1`, userID); err != nil {
		return User{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}

	user, _, err := s.FindUserByID(ctx, userID)
	return user, err
}

func (s *Store) ListLoginAttempts(ctx context.Context, query LoginAttemptListQuery) (LoginAttemptListResult, error) {
	rows, err := s.pool.Query(ctx, `
		select id, user_id, email_attempted, occurred_at, ip_address::text, ip_hash, coalesce(country_code::text, ''),
		       coalesce(city, ''), coalesce(user_agent, ''), coalesce(browser, ''), coalesce(operating_system, ''),
		       success, failure_reason, auth_method, risk_score, coalesce(correlation_id, '')
		from login_attempts
		where ($1::uuid is null or user_id = $1)
		  and ($2::boolean is null or success = $2)
		  and ($3 = '' or lower(email_attempted) like '%' || lower($3) || '%' or ip_address::text like '%' || $3 || '%')
		order by occurred_at desc
		limit $4 offset $5
	`, query.UserID, query.Success, query.Query, normalizedLimit(query.PageSize), normalizedOffset(query.Page, query.PageSize))
	if err != nil {
		return LoginAttemptListResult{}, err
	}
	defer rows.Close()

	items := []LoginAttempt{}
	for rows.Next() {
		attempt, err := scanLoginAttempt(rows)
		if err != nil {
			return LoginAttemptListResult{}, err
		}
		items = append(items, attempt)
	}
	if err := rows.Err(); err != nil {
		return LoginAttemptListResult{}, err
	}

	var total int
	if err := s.pool.QueryRow(ctx, `
		select count(*)
		from login_attempts
		where ($1::uuid is null or user_id = $1)
		  and ($2::boolean is null or success = $2)
		  and ($3 = '' or lower(email_attempted) like '%' || lower($3) || '%' or ip_address::text like '%' || $3 || '%')
	`, query.UserID, query.Success, query.Query).Scan(&total); err != nil {
		return LoginAttemptListResult{}, err
	}
	page, pageSize := normalizedPage(query.Page, query.PageSize)

	return LoginAttemptListResult{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Store) ListSecurityEvents(ctx context.Context, limit int) ([]SecurityEvent, error) {
	rows, err := s.pool.Query(ctx, `
		select id, user_id, type, severity, status, summary, detected_at, resolved_at,
		       coalesce(metadata->>'sourceIpAddress', ''), coalesce(metadata->>'countryCode', '')
		from security_events
		order by detected_at desc
		limit $1
	`, normalizedLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []SecurityEvent{}
	for rows.Next() {
		var event SecurityEvent
		if err := rows.Scan(&event.ID, &event.UserID, &event.Type, &event.Severity, &event.Status, &event.Summary, &event.DetectedAt, &event.ResolvedAt, &event.SourceIPAddress, &event.CountryCode); err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, rows.Err()
}

func (s *Store) GetDashboardStats(ctx context.Context) (DashboardStats, error) {
	var stats DashboardStats
	if err := s.pool.QueryRow(ctx, `
		select count(*),
		       count(*) filter (where status = 'active'),
		       count(*) filter (where status = 'locked'),
		       count(*) filter (where created_at > now() - interval '7 days'),
		       count(*) filter (where created_at > now() - interval '30 days')
		from users
	`).Scan(&stats.Users.Total, &stats.Users.Active, &stats.Users.Locked, &stats.Users.NewLast7Days, &stats.Users.NewLast30Days); err != nil {
		return DashboardStats{}, err
	}
	if err := s.pool.QueryRow(ctx, `
		select count(*),
		       count(*) filter (where success),
		       count(*) filter (where not success)
		from login_attempts
	`).Scan(&stats.LoginAttempts.Total, &stats.LoginAttempts.Successful, &stats.LoginAttempts.Failed); err != nil {
		return DashboardStats{}, err
	}
	if err := s.pool.QueryRow(ctx, `select count(*) from password_reset_tokens where created_at > now() - interval '24 hours'`).Scan(&stats.PasswordResetRequests); err != nil {
		return DashboardStats{}, err
	}
	if err := s.pool.QueryRow(ctx, `select count(*) from notifications`).Scan(&stats.Notifications); err != nil {
		return DashboardStats{}, err
	}
	if err := s.pool.QueryRow(ctx, `select count(*) from security_events where status = 'open'`).Scan(&stats.OpenSecurityEvents); err != nil {
		return DashboardStats{}, err
	}
	var err error
	stats.TopCountries, err = s.topCounts(ctx, `select coalesce(country_code::text, ''), count(*) from login_attempts where country_code is not null group by country_code order by count(*) desc limit 5`)
	if err != nil {
		return DashboardStats{}, err
	}
	stats.TopIPAddresses, err = s.topCounts(ctx, `select ip_address::text, count(*) from login_attempts group by ip_address order by count(*) desc limit 5`)
	if err != nil {
		return DashboardStats{}, err
	}

	return stats, nil
}
