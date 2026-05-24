package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *Store) scanUser(ctx context.Context, query string, args ...any) (User, error) {
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

func (s *Store) listLicenses(ctx context.Context, userID uuid.UUID) ([]string, error) {
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

func (s *Store) adminUserListItem(ctx context.Context, user User) (AdminUserListItem, error) {
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

func (s *Store) loginCounts(ctx context.Context, userID uuid.UUID) (int, int, error) {
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

func (s *Store) userPermissionGrants(ctx context.Context, userID uuid.UUID) ([]UserPermissionGrant, error) {
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

func (s *Store) userLoginFacets(ctx context.Context, userID uuid.UUID) ([]string, []string, []string, error) {
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

func (s *Store) topCounts(ctx context.Context, query string) ([]CountStat, error) {
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
