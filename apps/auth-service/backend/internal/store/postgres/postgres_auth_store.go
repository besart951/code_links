package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Store) FindUserByEmail(ctx context.Context, email string) (User, []string, error) {
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

func (s *Store) FindUserByID(ctx context.Context, userID uuid.UUID) (User, []string, error) {
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

func (s *Store) CreateUser(ctx context.Context, name string, email string, passwordHash string) (User, error) {
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, errConflict
		}
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

func (s *Store) GrantLicense(ctx context.Context, userID uuid.UUID, productID string) ([]string, error) {
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

func (s *Store) CreateRefreshSession(ctx context.Context, tokenHash string, userID uuid.UUID, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		insert into refresh_sessions (token_hash, user_id, expires_at)
		values ($1, $2, $3)
	`, tokenHash, userID, expiresAt)

	return err
}

func (s *Store) FindRefreshSession(ctx context.Context, tokenHash string) (uuid.UUID, error) {
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

func (s *Store) DeleteRefreshSession(ctx context.Context, tokenHash string) error {
	_, err := s.pool.Exec(ctx, `delete from refresh_sessions where token_hash = $1`, tokenHash)
	return err
}

func (s *Store) CreateEmailVerificationToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		insert into email_verification_tokens (token_hash, user_id, expires_at)
		values ($1, $2, $3)
	`, tokenHash, userID, expiresAt)

	return err
}

func (s *Store) VerifyEmailToken(ctx context.Context, tokenHash string, now time.Time) (User, error) {
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

func (s *Store) CreatePasswordResetToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `
		insert into password_reset_tokens (token_hash, user_id, expires_at)
		values ($1, $2, $3)
	`, tokenHash, userID, expiresAt)

	return err
}

func (s *Store) ResetPasswordByToken(ctx context.Context, tokenHash string, passwordHash string, now time.Time) (User, error) {
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

func (s *Store) RecordLoginAttempt(ctx context.Context, attempt LoginAttempt) error {
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
