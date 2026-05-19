package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/besart951/code_links/platform/auth"
	"github.com/besart951/code_links/platform/billing"
	"github.com/besart951/code_links/platform/entitlements"
	"github.com/besart951/code_links/platform/tenants"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) FindUserByEmail(ctx context.Context, email string) (auth.User, error) {
	row := s.db.QueryRow(ctx, `
		select id::text, email, password_hash, display_name, status, created_at, last_login_at
		from users
		where lower(email) = lower($1)
		limit 1
	`, email)
	return scanUser(row)
}

func (s *PostgresStore) FindUserByID(ctx context.Context, userID string) (auth.User, error) {
	row := s.db.QueryRow(ctx, `
		select id::text, email, password_hash, display_name, status, created_at, last_login_at
		from users
		where id = $1::uuid
		limit 1
	`, userID)
	return scanUser(row)
}

func (s *PostgresStore) TouchLastLogin(ctx context.Context, userID string, at time.Time) error {
	_, err := s.db.Exec(ctx, `update users set last_login_at = $2 where id = $1::uuid`, userID, at)
	return err
}

func (s *PostgresStore) StoreRefreshToken(ctx context.Context, token auth.RefreshToken) error {
	_, err := s.db.Exec(ctx, `
		insert into refresh_tokens (id, user_id, token_hash, user_agent, ip, expires_at, created_at)
		values ($1, $2::uuid, $3, $4, $5, $6, $7)
	`, token.ID, token.UserID, token.TokenHash, token.UserAgent, token.IP, token.ExpiresAt, token.CreatedAt)
	return err
}

func (s *PostgresStore) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (auth.RefreshToken, error) {
	row := s.db.QueryRow(ctx, `
		select id, user_id::text, token_hash, user_agent, ip, expires_at, created_at, revoked_at
		from refresh_tokens
		where token_hash = $1
		limit 1
	`, tokenHash)
	var token auth.RefreshToken
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.UserAgent,
		&token.IP,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.RevokedAt,
	)
	return token, err
}

func (s *PostgresStore) RevokeRefreshTokenByHash(ctx context.Context, tokenHash string, at time.Time) error {
	_, err := s.db.Exec(ctx, `
		update refresh_tokens
		set revoked_at = coalesce(revoked_at, $2)
		where token_hash = $1
	`, tokenHash, at)
	return err
}

func (s *PostgresStore) ListTenantsForUser(ctx context.Context, userID string) ([]tenants.Tenant, error) {
	rows, err := s.db.Query(ctx, `
		select t.id::text, t.type, t.name, t.slug, t.owner_user_id::text, t.status, t.billing_email, t.created_at, r.key
		from tenants t
		join tenant_members tm on tm.tenant_id = t.id
		join roles r on r.id = tm.role_id
		where tm.user_id = $1::uuid and tm.status = 'active' and t.status = 'active'
		order by t.name asc
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []tenants.Tenant
	for rows.Next() {
		var tenant tenants.Tenant
		if err := rows.Scan(
			&tenant.ID,
			&tenant.Type,
			&tenant.Name,
			&tenant.Slug,
			&tenant.OwnerUserID,
			&tenant.Status,
			&tenant.BillingEmail,
			&tenant.CreatedAt,
			&tenant.RoleKey,
		); err != nil {
			return nil, err
		}
		result = append(result, tenant)
	}
	return result, rows.Err()
}

func (s *PostgresStore) IsTenantMember(ctx context.Context, userID, tenantID string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx, `
		select exists (
			select 1
			from tenant_members tm
			join tenants t on t.id = tm.tenant_id
			where tm.user_id = $1::uuid
			  and tm.tenant_id = $2::uuid
			  and tm.status = 'active'
			  and t.status = 'active'
		)
	`, userID, tenantID).Scan(&exists)
	return exists, err
}

func (s *PostgresStore) ListEntitlements(ctx context.Context, tenantID string) ([]entitlements.Entitlement, error) {
	rows, err := s.db.Query(ctx, `
		select tenant_id::text, product_key, feature_key, source, enabled, expires_at
		from entitlements
		where tenant_id = $1::uuid
		order by product_key, feature_key
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []entitlements.Entitlement
	for rows.Next() {
		var entitlement entitlements.Entitlement
		if err := rows.Scan(
			&entitlement.TenantID,
			&entitlement.ProductKey,
			&entitlement.FeatureKey,
			&entitlement.Source,
			&entitlement.Enabled,
			&entitlement.ExpiresAt,
		); err != nil {
			return nil, err
		}
		result = append(result, entitlement)
	}
	return result, rows.Err()
}

func (s *PostgresStore) ListFeatureLimits(ctx context.Context, tenantID string) ([]billing.FeatureLimit, error) {
	rows, err := s.db.Query(ctx, `
		select tenant_id::text, product_key, feature_key, limit_key, value, period, reset_at
		from feature_limits
		where tenant_id = $1::uuid
		order by product_key, feature_key, limit_key
	`, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []billing.FeatureLimit
	for rows.Next() {
		var limit billing.FeatureLimit
		if err := rows.Scan(
			&limit.TenantID,
			&limit.ProductKey,
			&limit.FeatureKey,
			&limit.LimitKey,
			&limit.Value,
			&limit.Period,
			&limit.ResetAt,
		); err != nil {
			return nil, err
		}
		result = append(result, limit)
	}
	return result, rows.Err()
}

func (s *PostgresStore) HasEntitlement(ctx context.Context, tenantID, productKey, featureKey string) (bool, error) {
	var exists bool
	err := s.db.QueryRow(ctx, `
		select exists (
			select 1
			from entitlements
			where tenant_id = $1::uuid
			  and product_key = $2
			  and feature_key = $3
			  and enabled = true
			  and (expires_at is null or expires_at > now())
		)
	`, tenantID, productKey, featureKey).Scan(&exists)
	return exists, err
}

func scanUser(row pgx.Row) (auth.User, error) {
	var user auth.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.DisplayName,
		&user.Status,
		&user.CreatedAt,
		&user.LastLoginAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return auth.User{}, auth.ErrUnauthorized
	}
	return user, err
}
