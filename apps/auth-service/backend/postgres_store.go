package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresStore struct {
	pool *pgxpool.Pool
}

func openStore(ctx context.Context, config config) (Store, func(), error) {
	if config.DatabaseURL == "" {
		store, err := newMemoryStore()
		return store, func() {}, err
	}

	var lastErr error
	for attempt := 0; attempt < 30; attempt++ {
		pool, err := pgxpool.New(ctx, config.DatabaseURL)
		if err == nil {
			if err = pool.Ping(ctx); err == nil {
				return &postgresStore{pool: pool}, pool.Close, nil
			}
			pool.Close()
		}
		lastErr = err
		time.Sleep(time.Second)
	}

	return nil, func() {}, lastErr
}

func (s *postgresStore) FindUserByEmail(ctx context.Context, email string) (User, []string, error) {
	user, err := s.scanUser(ctx, `
		select id, email, name, password_hash, status, email_verified_at, created_at, updated_at,
		       last_login_at, coalesce(last_login_ip::text, ''), coalesce(last_login_country_code::text, '')
		from users
		where lower(email) = lower($1)
	`, email)
	if err != nil {
		return User{}, nil, err
	}

	licenses, err := s.listLicenses(ctx, user.ID)
	return user, licenses, err
}

func (s *postgresStore) FindUserByID(ctx context.Context, userID uuid.UUID) (User, []string, error) {
	user, err := s.scanUser(ctx, `
		select id, email, name, password_hash, status, email_verified_at, created_at, updated_at,
		       last_login_at, coalesce(last_login_ip::text, ''), coalesce(last_login_country_code::text, '')
		from users
		where id = $1
	`, userID)
	if err != nil {
		return User{}, nil, err
	}

	licenses, err := s.listLicenses(ctx, user.ID)
	return user, licenses, err
}

func (s *postgresStore) CreateUser(ctx context.Context, name string, email string, passwordHash string) (User, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, errConflict
		}
		return User{}, err
	}
	defer tx.Rollback(ctx)

	var user User
	err = tx.QueryRow(ctx, `
		insert into users (email, name, password_hash, status)
		values (lower($1), $2, $3, 'active')
		returning id, email, name, password_hash, status, email_verified_at, created_at, updated_at,
		          last_login_at, coalesce(last_login_ip::text, ''), coalesce(last_login_country_code::text, '')
	`, email, name, passwordHash).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.Status,
		&user.EmailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
		&user.LastLoginIP,
		&user.LastLoginCountryCode,
	)
	if err != nil {
		return User{}, err
	}

	if _, err := tx.Exec(ctx, `
		insert into user_roles (user_id, role_id)
		values ($1, 'user')
		on conflict do nothing
	`, user.ID); err != nil {
		return User{}, err
	}

	return user, tx.Commit(ctx)
}

func (s *postgresStore) GrantLicense(ctx context.Context, userID uuid.UUID, productID string) ([]string, error) {
	_, err := s.pool.Exec(ctx, `
		insert into user_licenses (user_id, product_id)
		values ($1, $2)
		on conflict (user_id, product_id) do nothing
	`, userID, productID)
	if err != nil {
		return nil, err
	}

	return s.listLicenses(ctx, userID)
}

func (s *postgresStore) CreateRefreshSession(ctx context.Context, tokenHash string, userID uuid.UUID, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		insert into refresh_sessions (token_hash, user_id, expires_at)
		values ($1, $2, $3)
	`, tokenHash, userID, expiresAt)

	return err
}

func (s *postgresStore) FindRefreshSession(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	var userID uuid.UUID
	err := s.pool.QueryRow(ctx, `
		select user_id
		from refresh_sessions
		where token_hash = $1
		  and expires_at > now()
	`, tokenHash).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, errNotFound
	}

	return userID, err
}

func (s *postgresStore) DeleteRefreshSession(ctx context.Context, tokenHash string) error {
	_, err := s.pool.Exec(ctx, `delete from refresh_sessions where token_hash = $1`, tokenHash)
	return err
}

func (s *postgresStore) CreateEmailVerificationToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		insert into email_verification_tokens (token_hash, user_id, expires_at)
		values ($1, $2, $3)
	`, tokenHash, userID, expiresAt)

	return err
}

func (s *postgresStore) VerifyEmailToken(ctx context.Context, tokenHash string, now time.Time) (User, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback(ctx)

	var userID uuid.UUID
	err = tx.QueryRow(ctx, `
		update email_verification_tokens
		set used_at = $2
		where token_hash = $1
		  and used_at is null
		  and expires_at > $2
		returning user_id
	`, tokenHash, now).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, errNotFound
	}
	if err != nil {
		return User{}, err
	}

	if _, err := tx.Exec(ctx, `update users set email_verified_at = $2, updated_at = $2 where id = $1`, userID, now); err != nil {
		return User{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}

	user, _, err := s.FindUserByID(ctx, userID)
	return user, err
}

func (s *postgresStore) CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		insert into password_reset_tokens (token_hash, user_id, expires_at)
		values ($1, $2, $3)
	`, tokenHash, userID, expiresAt)

	return err
}

func (s *postgresStore) ResetPasswordByToken(ctx context.Context, tokenHash string, passwordHash string, now time.Time) (User, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return User{}, err
	}
	defer tx.Rollback(ctx)

	var userID uuid.UUID
	err = tx.QueryRow(ctx, `
		update password_reset_tokens
		set used_at = $2
		where token_hash = $1
		  and used_at is null
		  and expires_at > $2
		returning user_id
	`, tokenHash, now).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, errNotFound
	}
	if err != nil {
		return User{}, err
	}
	if _, err := tx.Exec(ctx, `update users set password_hash = $2, updated_at = $3 where id = $1`, userID, passwordHash, now); err != nil {
		return User{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}

	user, _, err := s.FindUserByID(ctx, userID)
	return user, err
}

func (s *postgresStore) RecordLoginAttempt(ctx context.Context, attempt LoginAttempt) error {
	if attempt.ID == uuid.Nil {
		attempt.ID = uuid.New()
	}
	if attempt.OccurredAt.IsZero() {
		attempt.OccurredAt = time.Now().UTC()
	}

	_, err := s.pool.Exec(ctx, `
		insert into login_attempts (
			id, user_id, email_attempted, occurred_at, ip_address, ip_hash, country_code, city,
			user_agent, browser, operating_system, success, failure_reason, auth_method, risk_score, correlation_id
		)
		values ($1, $2, $3, $4, $5::inet, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`, attempt.ID, attempt.UserID, attempt.EmailAttempted, attempt.OccurredAt, attempt.IPAddress, attempt.IPHash, attempt.CountryCode, attempt.City, attempt.UserAgent, attempt.Browser, attempt.OperatingSystem, attempt.Success, attempt.FailureReason, attempt.AuthMethod, attempt.RiskScore, attempt.CorrelationID)
	if err != nil {
		return err
	}
	if attempt.Success && attempt.UserID != nil {
		_, err = s.pool.Exec(ctx, `
			update users
			set last_login_at = $2,
			    last_login_ip = $3::inet,
			    last_login_country_code = $4,
			    updated_at = $2
			where id = $1
		`, *attempt.UserID, attempt.OccurredAt, attempt.IPAddress, attempt.CountryCode)
	}

	return err
}

func (s *postgresStore) GetAdminActor(ctx context.Context, userID uuid.UUID) (AdminActor, error) {
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

func (s *postgresStore) ListUserRoles(ctx context.Context, userID uuid.UUID) ([]AdminRole, []AdminPermission, error) {
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

func (s *postgresStore) ListAdminUsers(ctx context.Context, query AdminUserListQuery) (AdminUserListResult, error) {
	rows, err := s.pool.Query(ctx, `
		select id, email, name, password_hash, status, email_verified_at, created_at, updated_at,
		       last_login_at, coalesce(last_login_ip::text, ''), coalesce(last_login_country_code::text, '')
		from users
		order by created_at desc
	`)
	if err != nil {
		return AdminUserListResult{}, err
	}
	defer rows.Close()

	items := []AdminUserListItem{}
	for rows.Next() {
		user, err := scanUserRow(rows)
		if err != nil {
			return AdminUserListResult{}, err
		}
		item, err := s.adminUserListItem(ctx, user)
		if err != nil {
			return AdminUserListResult{}, err
		}
		if matchesUserQuery(item, query) {
			items = append(items, item)
		}
	}
	if err := rows.Err(); err != nil {
		return AdminUserListResult{}, err
	}

	sortAdminUsers(items, query.Sort, query.Direction)
	total := len(items)
	page, pageSize := normalizedPage(query.Page, query.PageSize)
	start := (page - 1) * pageSize
	if start > len(items) {
		start = len(items)
	}
	end := start + pageSize
	if end > len(items) {
		end = len(items)
	}

	return AdminUserListResult{Items: items[start:end], Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *postgresStore) GetManagedUserDetail(ctx context.Context, userID uuid.UUID) (ManagedUserDetail, error) {
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

func (s *postgresStore) SetUserStatus(ctx context.Context, userID uuid.UUID, status UserStatus) (User, error) {
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

func (s *postgresStore) SetUserRole(ctx context.Context, userID uuid.UUID, role AdminRole) (User, error) {
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

func (s *postgresStore) ListLoginAttempts(ctx context.Context, query LoginAttemptListQuery) (LoginAttemptListResult, error) {
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

func (s *postgresStore) ListSecurityEvents(ctx context.Context, limit int) ([]SecurityEvent, error) {
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

func (s *postgresStore) GetDashboardStats(ctx context.Context) (DashboardStats, error) {
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
	_ = s.pool.QueryRow(ctx, `select count(*) from password_reset_tokens where created_at > now() - interval '24 hours'`).Scan(&stats.PasswordResetRequests)
	_ = s.pool.QueryRow(ctx, `select count(*) from notifications`).Scan(&stats.Notifications)
	_ = s.pool.QueryRow(ctx, `select count(*) from security_events where status = 'open'`).Scan(&stats.OpenSecurityEvents)
	stats.TopCountries, _ = s.topCounts(ctx, `select coalesce(country_code::text, ''), count(*) from login_attempts where country_code is not null group by country_code order by count(*) desc limit 5`)
	stats.TopIPAddresses, _ = s.topCounts(ctx, `select ip_address::text, count(*) from login_attempts group by ip_address order by count(*) desc limit 5`)

	return stats, nil
}

func (s *postgresStore) GetSmtpSettings(ctx context.Context) (SmtpSettings, error) {
	var settings SmtpSettings
	err := s.pool.QueryRow(ctx, `
		select host, port, username, coalesce(password_encrypted, ''), encryption, from_email, from_name, reply_to_email, active, updated_at
		from smtp_settings
		where id = 'default'
	`).Scan(&settings.Host, &settings.Port, &settings.Username, &settings.PasswordEncrypted, &settings.Encryption, &settings.FromEmail, &settings.FromName, &settings.ReplyToEmail, &settings.Active, &settings.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		settings.Port = 587
		settings.Encryption = SmtpEncryptionStartTLS
		settings.FromEmail = "no-reply@codelinks.dev"
		settings.FromName = "CodeLinks"
		settings.ReplyToEmail = "support@codelinks.dev"
		settings.UpdatedAt = time.Now().UTC()
		return settings, nil
	}
	settings.HasPassword = settings.PasswordEncrypted != ""

	return settings, err
}

func (s *postgresStore) SaveSmtpSettings(ctx context.Context, settings SmtpSettings) (SmtpSettings, error) {
	current, _ := s.GetSmtpSettings(ctx)
	if settings.PasswordEncrypted == "" {
		settings.PasswordEncrypted = current.PasswordEncrypted
	}
	err := s.pool.QueryRow(ctx, `
		insert into smtp_settings (id, host, port, username, password_encrypted, encryption, from_email, from_name, reply_to_email, active, updated_at)
		values ('default', $1, $2, $3, $4, $5, $6, $7, $8, $9, now())
		on conflict (id) do update set
			host = excluded.host,
			port = excluded.port,
			username = excluded.username,
			password_encrypted = excluded.password_encrypted,
			encryption = excluded.encryption,
			from_email = excluded.from_email,
			from_name = excluded.from_name,
			reply_to_email = excluded.reply_to_email,
			active = excluded.active,
			updated_at = now()
		returning host, port, username, coalesce(password_encrypted, ''), encryption, from_email, from_name, reply_to_email, active, updated_at
	`, settings.Host, settings.Port, settings.Username, settings.PasswordEncrypted, settings.Encryption, settings.FromEmail, settings.FromName, settings.ReplyToEmail, settings.Active).
		Scan(&settings.Host, &settings.Port, &settings.Username, &settings.PasswordEncrypted, &settings.Encryption, &settings.FromEmail, &settings.FromName, &settings.ReplyToEmail, &settings.Active, &settings.UpdatedAt)
	settings.HasPassword = settings.PasswordEncrypted != ""

	return settings, err
}

func (s *postgresStore) ListNotifications(ctx context.Context, limit int) ([]Notification, error) {
	rows, err := s.pool.Query(ctx, `
		select id, type, channel, recipient, subject, status, created_at, sent_at
		from notifications
		order by created_at desc
		limit $1
	`, normalizedLimit(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []Notification{}
	for rows.Next() {
		var notification Notification
		if err := rows.Scan(&notification.ID, &notification.Type, &notification.Channel, &notification.Recipient, &notification.Subject, &notification.Status, &notification.CreatedAt, &notification.SentAt); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, rows.Err()
}

func (s *postgresStore) CreateNotification(ctx context.Context, notification Notification) error {
	if notification.ID == uuid.Nil {
		notification.ID = uuid.New()
	}
	if notification.CreatedAt.IsZero() {
		notification.CreatedAt = time.Now().UTC()
	}
	_, err := s.pool.Exec(ctx, `
		insert into notifications (id, type, channel, recipient, subject, status, created_at, sent_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8)
	`, notification.ID, notification.Type, notification.Channel, notification.Recipient, notification.Subject, notification.Status, notification.CreatedAt, notification.SentAt)

	return err
}

func (s *postgresStore) RecordAdminAuditEntry(ctx context.Context, entry AdminAuditEntry) error {
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	_, err := s.pool.Exec(ctx, `
		insert into admin_audit_entries (id, actor_user_id, action, target_type, target_id, reason, ip_address, created_at)
		values ($1, $2, $3, $4, $5, $6, nullif($7, '')::inet, $8)
	`, entry.ID, entry.ActorUserID, entry.Action, entry.TargetType, entry.TargetID, entry.Reason, entry.IPAddress, entry.CreatedAt)

	return err
}

func (s *postgresStore) scanUser(ctx context.Context, query string, args ...any) (User, error) {
	user, err := scanUserRow(s.pool.QueryRow(ctx, query, args...))
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, errNotFound
	}

	return user, err
}

type scanner interface {
	Scan(dest ...any) error
}

func scanUserRow(row scanner) (User, error) {
	var user User
	var emailVerifiedAt sql.NullTime
	var lastLoginAt sql.NullTime
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.Status,
		&emailVerifiedAt,
		&user.CreatedAt,
		&user.UpdatedAt,
		&lastLoginAt,
		&user.LastLoginIP,
		&user.LastLoginCountryCode,
	); err != nil {
		return User{}, err
	}
	if emailVerifiedAt.Valid {
		user.EmailVerifiedAt = &emailVerifiedAt.Time
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}

	return user, nil
}

func scanLoginAttempt(row scanner) (LoginAttempt, error) {
	var attempt LoginAttempt
	var failure sql.NullString
	if err := row.Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.EmailAttempted,
		&attempt.OccurredAt,
		&attempt.IPAddress,
		&attempt.IPHash,
		&attempt.CountryCode,
		&attempt.City,
		&attempt.UserAgent,
		&attempt.Browser,
		&attempt.OperatingSystem,
		&attempt.Success,
		&failure,
		&attempt.AuthMethod,
		&attempt.RiskScore,
		&attempt.CorrelationID,
	); err != nil {
		return LoginAttempt{}, err
	}
	if failure.Valid {
		reason := LoginFailureReason(failure.String)
		attempt.FailureReason = &reason
	}

	return attempt, nil
}

func (s *postgresStore) listLicenses(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := s.pool.Query(ctx, `
		select product_id
		from user_licenses
		where user_id = $1
		order by product_id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	licenses := []string{}
	for rows.Next() {
		var productID string
		if err := rows.Scan(&productID); err != nil {
			return nil, err
		}
		licenses = append(licenses, productID)
	}

	return licenses, rows.Err()
}

func (s *postgresStore) adminUserListItem(ctx context.Context, user User) (AdminUserListItem, error) {
	roles, _, err := s.ListUserRoles(ctx, user.ID)
	if err != nil {
		return AdminUserListItem{}, err
	}
	successful, failed, err := s.loginCounts(ctx, user.ID)
	if err != nil {
		return AdminUserListItem{}, err
	}

	return AdminUserListItem{
		ID:                   user.ID.String(),
		Name:                 user.Name,
		Email:                user.Email,
		PrimaryRole:          primaryRole(roles),
		Status:               user.Status,
		EmailVerified:        user.EmailVerifiedAt != nil,
		CreatedAt:            user.CreatedAt,
		LastLoginAt:          user.LastLoginAt,
		SuccessfulLoginCount: successful,
		FailedLoginCount:     failed,
		LastKnownIPAddress:   emptyStringToNil(user.LastLoginIP),
		LastLoginCountryCode: emptyStringToNil(user.LastLoginCountryCode),
	}, nil
}

func (s *postgresStore) loginCounts(ctx context.Context, userID uuid.UUID) (int, int, error) {
	var successful int
	var failed int
	err := s.pool.QueryRow(ctx, `
		select count(*) filter (where success),
		       count(*) filter (where not success)
		from login_attempts
		where user_id = $1
	`, userID).Scan(&successful, &failed)

	return successful, failed, err
}

func (s *postgresStore) userPermissionGrants(ctx context.Context, userID uuid.UUID) ([]UserPermissionGrant, error) {
	rows, err := s.pool.Query(ctx, `
		select role_id, granted_at, coalesce(granted_by::text, 'system')
		from user_roles
		where user_id = $1
		order by role_id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	grants := []UserPermissionGrant{}
	for rows.Next() {
		var grant UserPermissionGrant
		if err := rows.Scan(&grant.Role, &grant.GrantedAt, &grant.GrantedBy); err != nil {
			return nil, err
		}
		grants = append(grants, grant)
	}

	return grants, rows.Err()
}

func (s *postgresStore) userLoginFacets(ctx context.Context, userID uuid.UUID) ([]string, []string, []string, error) {
	rows, err := s.pool.Query(ctx, `
		select distinct ip_address::text, coalesce(country_code::text, ''), trim(coalesce(browser, '') || ' / ' || coalesce(operating_system, ''))
		from login_attempts
		where user_id = $1
	`, userID)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()

	ips := []string{}
	countries := []string{}
	devices := []string{}
	for rows.Next() {
		var ip string
		var country string
		var device string
		if err := rows.Scan(&ip, &country, &device); err != nil {
			return nil, nil, nil, err
		}
		if ip != "" {
			ips = append(ips, ip)
		}
		if country != "" {
			countries = append(countries, country)
		}
		if device != "" && device != "/" {
			devices = append(devices, device)
		}
	}

	return ips, countries, devices, rows.Err()
}

func (s *postgresStore) topCounts(ctx context.Context, query string) ([]CountStat, error) {
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := []CountStat{}
	for rows.Next() {
		var stat CountStat
		if err := rows.Scan(&stat.Key, &stat.Count); err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

func normalizedLimit(limit int) int {
	if limit < 1 {
		return 50
	}
	if limit > 100 {
		return 100
	}

	return limit
}

func normalizedOffset(page int, pageSize int) int {
	page, pageSize = normalizedPage(page, pageSize)
	return (page - 1) * pageSize
}
