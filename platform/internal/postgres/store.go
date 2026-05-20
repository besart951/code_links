package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/besart951/code_links/platform/internal/access"
	"github.com/besart951/code_links/platform/internal/auth"
	"github.com/besart951/code_links/platform/internal/tenants"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func (s *Store) FindByEmail(ctx context.Context, email string) (auth.User, error) {
	row := s.db.QueryRow(ctx, `
		select id::text, email, display_name, password_hash, status, token_version, created_at, last_login_at
		from users
		where lower(email) = lower($1)
		limit 1
	`, email)
	return scanAuthUser(row)
}

func (s *Store) FindByID(ctx context.Context, userID auth.UserID) (auth.User, error) {
	row := s.db.QueryRow(ctx, `
		select id::text, email, display_name, password_hash, status, token_version, created_at, last_login_at
		from users
		where id = $1::uuid
		limit 1
	`, string(userID))
	return scanAuthUser(row)
}

func (s *Store) TouchLastLogin(ctx context.Context, userID auth.UserID, at time.Time) error {
	_, err := s.db.Exec(ctx, `update users set last_login_at = $2 where id = $1::uuid`, string(userID), at)
	return err
}

func (s *Store) CreateSession(ctx context.Context, session auth.Session, refresh auth.RefreshToken) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		insert into sessions (id, user_id, token_hash, token_version, user_agent, ip, created_at, expires_at)
		values ($1, $2::uuid, $3, $4, $5, $6, $7, $8)
	`, string(session.ID), string(session.UserID), session.TokenHash, int(session.TokenVersion), nullableString(session.UserAgent), nullableString(session.IP), session.CreatedAt, session.ExpiresAt); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		insert into refresh_tokens (id, session_id, user_id, token_hash, user_agent, ip, expires_at, created_at)
		values ($1, $2, $3::uuid, $4, $5, $6, $7, $8)
	`, string(refresh.ID), string(refresh.SessionID), string(refresh.UserID), refresh.TokenHash, nullableString(refresh.UserAgent), nullableString(refresh.IP), refresh.ExpiresAt, refresh.CreatedAt); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) FindBySessionTokenHash(ctx context.Context, tokenHash string) (auth.AuthenticatedSession, error) {
	row := s.db.QueryRow(ctx, `
		select
			u.id::text, u.email, u.display_name, u.password_hash, u.status, u.token_version, u.created_at, u.last_login_at,
			s.id, s.user_id::text, s.token_hash, s.token_version, s.user_agent, s.ip, s.created_at, s.last_seen_at, s.revoked_at, s.expires_at
		from sessions s
		join users u on u.id = s.user_id
		where s.token_hash = $1
		limit 1
	`, tokenHash)
	var result auth.AuthenticatedSession
	var userAgent, ip *string
	err := row.Scan(
		&result.User.ID,
		&result.User.Email,
		&result.User.DisplayName,
		&result.User.PasswordHash,
		&result.User.Status,
		&result.User.TokenVersion,
		&result.User.CreatedAt,
		&result.User.LastLoginAt,
		&result.Session.ID,
		&result.Session.UserID,
		&result.Session.TokenHash,
		&result.Session.TokenVersion,
		&userAgent,
		&ip,
		&result.Session.CreatedAt,
		&result.Session.LastSeenAt,
		&result.Session.RevokedAt,
		&result.Session.ExpiresAt,
	)
	if err != nil {
		return auth.AuthenticatedSession{}, mapNoRows(err, auth.ErrUnauthorized)
	}
	result.Session.UserAgent = valueOrEmpty(userAgent)
	result.Session.IP = valueOrEmpty(ip)
	return result, nil
}

func (s *Store) FindByRefreshTokenHash(ctx context.Context, tokenHash string) (auth.Session, auth.RefreshToken, error) {
	row := s.db.QueryRow(ctx, `
		select
			s.id, s.user_id::text, s.token_hash, s.token_version, s.user_agent, s.ip, s.created_at, s.last_seen_at, s.revoked_at, s.expires_at,
			rt.id, rt.session_id, rt.user_id::text, rt.token_hash, rt.user_agent, rt.ip, rt.expires_at, rt.created_at, rt.revoked_at
		from refresh_tokens rt
		join sessions s on s.id = rt.session_id
		where rt.token_hash = $1
		limit 1
	`, tokenHash)
	var session auth.Session
	var refresh auth.RefreshToken
	var sessionUserAgent, sessionIP, refreshUserAgent, refreshIP *string
	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.TokenVersion,
		&sessionUserAgent,
		&sessionIP,
		&session.CreatedAt,
		&session.LastSeenAt,
		&session.RevokedAt,
		&session.ExpiresAt,
		&refresh.ID,
		&refresh.SessionID,
		&refresh.UserID,
		&refresh.TokenHash,
		&refreshUserAgent,
		&refreshIP,
		&refresh.ExpiresAt,
		&refresh.CreatedAt,
		&refresh.RevokedAt,
	)
	if err != nil {
		return auth.Session{}, auth.RefreshToken{}, mapNoRows(err, auth.ErrUnauthorized)
	}
	session.UserAgent = valueOrEmpty(sessionUserAgent)
	session.IP = valueOrEmpty(sessionIP)
	refresh.UserAgent = valueOrEmpty(refreshUserAgent)
	refresh.IP = valueOrEmpty(refreshIP)
	return session, refresh, nil
}

func (s *Store) RotateRefreshToken(ctx context.Context, oldTokenHash string, next auth.RefreshToken, revokedAt time.Time) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		update refresh_tokens
		set revoked_at = coalesce(revoked_at, $2)
		where token_hash = $1
	`, oldTokenHash, revokedAt); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `
		insert into refresh_tokens (id, session_id, user_id, token_hash, user_agent, ip, expires_at, created_at)
		values ($1, $2, $3::uuid, $4, $5, $6, $7, $8)
	`, string(next.ID), string(next.SessionID), string(next.UserID), next.TokenHash, nullableString(next.UserAgent), nullableString(next.IP), next.ExpiresAt, next.CreatedAt); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `update sessions set last_seen_at = $2 where id = $1`, string(next.SessionID), revokedAt); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *Store) RevokeSession(ctx context.Context, sessionID auth.SessionID, revokedAt time.Time) error {
	_, err := s.db.Exec(ctx, `
		update sessions
		set revoked_at = coalesce(revoked_at, $2)
		where id = $1
	`, string(sessionID), revokedAt)
	return err
}

func (s *Store) SessionAccess(ctx context.Context, sessionID access.SessionID) (access.SessionAccess, error) {
	row := s.db.QueryRow(ctx, `
		select s.id, s.user_id::text, s.token_version, u.token_version, s.revoked_at, s.expires_at
		from sessions s
		join users u on u.id = s.user_id
		where s.id = $1
		limit 1
	`, string(sessionID))
	var result access.SessionAccess
	err := row.Scan(
		&result.ID,
		&result.UserID,
		&result.TokenVersion,
		&result.UserTokenVersion,
		&result.RevokedAt,
		&result.ExpiresAt,
	)
	return result, mapNoRows(err, access.ErrSessionInactive)
}

func (s *Store) FindTenant(ctx context.Context, tenantID tenants.TenantID) (tenants.Tenant, error) {
	row := s.db.QueryRow(ctx, `
		select id::text, type, name, slug, owner_user_id::text, status, billing_email, created_at
		from tenants
		where id = $1::uuid
		limit 1
	`, string(tenantID))
	var tenant tenants.Tenant
	var billingEmail *string
	err := row.Scan(&tenant.ID, &tenant.Type, &tenant.Name, &tenant.Slug, &tenant.OwnerUserID, &tenant.Status, &billingEmail, &tenant.CreatedAt)
	if err != nil {
		return tenants.Tenant{}, mapNoRows(err, tenants.ErrTenantNotFound)
	}
	tenant.BillingEmail = valueOrEmpty(billingEmail)
	return tenant, nil
}

func (s *Store) FindMember(ctx context.Context, userID tenants.UserID, tenantID tenants.TenantID) (tenants.TenantMember, error) {
	row := s.db.QueryRow(ctx, `
		select tm.tenant_id::text, tm.user_id::text, tm.role_id::text, tm.status, tm.joined_at
		from tenant_members tm
		join tenants t on t.id = tm.tenant_id
		where tm.user_id = $1::uuid
		  and tm.tenant_id = $2::uuid
		  and t.status = 'active'
		limit 1
	`, string(userID), string(tenantID))
	var member tenants.TenantMember
	err := row.Scan(&member.TenantID, &member.UserID, &member.RoleID, &member.Status, &member.JoinedAt)
	return member, mapNoRows(err, tenants.ErrTenantMembership)
}

func (s *Store) ListTenantsForUser(ctx context.Context, userID tenants.UserID) ([]tenants.Tenant, error) {
	rows, err := s.db.Query(ctx, `
		select t.id::text, t.type, t.name, t.slug, t.owner_user_id::text, t.status, t.billing_email, t.created_at
		from tenants t
		join tenant_members tm on tm.tenant_id = t.id
		where tm.user_id = $1::uuid and tm.status = 'active' and t.status = 'active'
		order by t.name asc
	`, string(userID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []tenants.Tenant
	for rows.Next() {
		var tenant tenants.Tenant
		var billingEmail *string
		if err := rows.Scan(&tenant.ID, &tenant.Type, &tenant.Name, &tenant.Slug, &tenant.OwnerUserID, &tenant.Status, &billingEmail, &tenant.CreatedAt); err != nil {
			return nil, err
		}
		tenant.BillingEmail = valueOrEmpty(billingEmail)
		result = append(result, tenant)
	}
	return result, rows.Err()
}

func (s *Store) TenantAccess(ctx context.Context, tenantID access.TenantID, product access.ProductKey) (access.TenantAccess, error) {
	var result access.TenantAccess
	row := s.db.QueryRow(ctx, `
		select t.id::text, t.type, coalesce(tpv.entitlements_version, 1)
		from tenants t
		left join tenant_product_versions tpv on tpv.tenant_id = t.id and tpv.product_key = $2
		where t.id = $1::uuid and t.status = 'active'
		limit 1
	`, string(tenantID), string(product))
	if err := row.Scan(&result.TenantID, &result.TenantType, &result.EntitlementsVersion); err != nil {
		return access.TenantAccess{}, mapNoRows(err, tenants.ErrTenantNotFound)
	}
	result.ProductKey = product

	subscriptions, planKey, err := s.subscriptionsForProduct(ctx, tenantID, product)
	if err != nil {
		return access.TenantAccess{}, err
	}
	result.Subscriptions = subscriptions
	result.PlanKey = planKey

	entitlements, err := s.entitlementsForProduct(ctx, tenantID, product)
	if err != nil {
		return access.TenantAccess{}, err
	}
	result.Entitlements = entitlements

	limits, err := s.featureLimitsForProduct(ctx, tenantID, product)
	if err != nil {
		return access.TenantAccess{}, err
	}
	result.FeatureLimits = limits
	return result, nil
}

func (s *Store) MemberAccess(ctx context.Context, userID access.UserID, tenantID access.TenantID, product access.ProductKey) (access.MemberAccess, error) {
	member, err := s.FindMember(ctx, tenants.UserID(userID), tenants.TenantID(tenantID))
	if err != nil {
		return access.MemberAccess{}, access.ErrTenantMembership
	}
	if !member.IsActive() {
		return access.MemberAccess{}, access.ErrTenantMembership
	}
	result := access.MemberAccess{UserID: userID, TenantID: tenantID}

	roleRows, err := s.db.Query(ctx, `
		select distinct r.key, r.name, coalesce(r.product_key, '')
		from tenant_members tm
		join roles r on r.id = tm.role_id
		where tm.user_id = $1::uuid
		  and tm.tenant_id = $2::uuid
		  and tm.status = 'active'
		  and (r.product_key is null or r.product_key = $3)
		order by r.key
	`, string(userID), string(tenantID), string(product))
	if err != nil {
		return access.MemberAccess{}, err
	}
	defer roleRows.Close()
	for roleRows.Next() {
		var role access.Role
		if err := roleRows.Scan(&role.Key, &role.Name, &role.ProductKey); err != nil {
			return access.MemberAccess{}, err
		}
		result.Roles = append(result.Roles, role)
	}
	if err := roleRows.Err(); err != nil {
		return access.MemberAccess{}, err
	}

	permissionRows, err := s.db.Query(ctx, `
		select distinct p.key, p.description, coalesce(p.product_key, '')
		from tenant_members tm
		join role_permissions rp on rp.role_id = tm.role_id
		join permissions p on p.id = rp.permission_id
		where tm.user_id = $1::uuid
		  and tm.tenant_id = $2::uuid
		  and tm.status = 'active'
		  and (p.product_key is null or p.product_key = $3)
		order by p.key
	`, string(userID), string(tenantID), string(product))
	if err != nil {
		return access.MemberAccess{}, err
	}
	defer permissionRows.Close()
	for permissionRows.Next() {
		var permission access.Permission
		if err := permissionRows.Scan(&permission.Key, &permission.Description, &permission.ProductKey); err != nil {
			return access.MemberAccess{}, err
		}
		result.Permissions = append(result.Permissions, permission)
	}
	return result, permissionRows.Err()
}

func (s *Store) ListProducts(ctx context.Context) ([]access.Product, error) {
	rows, err := s.db.Query(ctx, `select key, name, status from products order by key`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []access.Product
	for rows.Next() {
		var product access.Product
		if err := rows.Scan(&product.Key, &product.Name, &product.Status); err != nil {
			return nil, err
		}
		result = append(result, product)
	}
	return result, rows.Err()
}

func (s *Store) ListEntitlements(ctx context.Context, tenantID access.TenantID) ([]access.Entitlement, error) {
	rows, err := s.db.Query(ctx, `
		select tenant_id::text, product_key, feature_key, source, enabled, expires_at
		from entitlements
		where tenant_id = $1::uuid
		order by product_key, feature_key
	`, string(tenantID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []access.Entitlement
	for rows.Next() {
		var item access.Entitlement
		if err := rows.Scan(&item.TenantID, &item.ProductKey, &item.FeatureKey, &item.Source, &item.Enabled, &item.ExpiresAt); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) ListFeatureLimits(ctx context.Context, tenantID access.TenantID) ([]access.FeatureLimit, error) {
	rows, err := s.db.Query(ctx, `
		select tenant_id::text, product_key, feature_key, limit_key, value, period, reset_at
		from feature_limits
		where tenant_id = $1::uuid
		order by product_key, feature_key, limit_key
	`, string(tenantID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []access.FeatureLimit
	for rows.Next() {
		var item access.FeatureLimit
		if err := rows.Scan(&item.TenantID, &item.ProductKey, &item.FeatureKey, &item.LimitKey, &item.Value, &item.Period, &item.ResetAt); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) subscriptionsForProduct(ctx context.Context, tenantID access.TenantID, product access.ProductKey) ([]access.Subscription, access.PlanKey, error) {
	rows, err := s.db.Query(ctx, `
		select s.id::text, s.tenant_id::text, coalesce(p.product_key, $2), p.name, s.status,
		       s.current_period_start, s.current_period_end, s.cancel_at
		from subscriptions s
		join plans p on p.id = s.plan_id
		where s.tenant_id = $1::uuid
		  and (p.product_key = $2 or p.bundle_key is not null)
		order by s.created_at desc
	`, string(tenantID), string(product))
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()
	var result []access.Subscription
	var planKey access.PlanKey
	now := time.Now().UTC()
	for rows.Next() {
		var item access.Subscription
		if err := rows.Scan(&item.ID, &item.TenantID, &item.ProductKey, &item.PlanKey, &item.Status, &item.CurrentPeriodStart, &item.CurrentPeriodEnd, &item.CancelAt); err != nil {
			return nil, "", err
		}
		if planKey == "" && item.IsActiveAt(now) {
			planKey = item.PlanKey
		}
		result = append(result, item)
	}
	return result, planKey, rows.Err()
}

func (s *Store) entitlementsForProduct(ctx context.Context, tenantID access.TenantID, product access.ProductKey) ([]access.Entitlement, error) {
	rows, err := s.db.Query(ctx, `
		select tenant_id::text, product_key, feature_key, source, enabled, expires_at
		from entitlements
		where tenant_id = $1::uuid and product_key = $2
		order by feature_key
	`, string(tenantID), string(product))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []access.Entitlement
	for rows.Next() {
		var item access.Entitlement
		if err := rows.Scan(&item.TenantID, &item.ProductKey, &item.FeatureKey, &item.Source, &item.Enabled, &item.ExpiresAt); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func (s *Store) featureLimitsForProduct(ctx context.Context, tenantID access.TenantID, product access.ProductKey) ([]access.FeatureLimit, error) {
	rows, err := s.db.Query(ctx, `
		select tenant_id::text, product_key, feature_key, limit_key, value, period, reset_at
		from feature_limits
		where tenant_id = $1::uuid and product_key = $2
		order by feature_key, limit_key
	`, string(tenantID), string(product))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []access.FeatureLimit
	for rows.Next() {
		var item access.FeatureLimit
		if err := rows.Scan(&item.TenantID, &item.ProductKey, &item.FeatureKey, &item.LimitKey, &item.Value, &item.Period, &item.ResetAt); err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, rows.Err()
}

func scanAuthUser(row pgx.Row) (auth.User, error) {
	var user auth.User
	err := row.Scan(&user.ID, &user.Email, &user.DisplayName, &user.PasswordHash, &user.Status, &user.TokenVersion, &user.CreatedAt, &user.LastLoginAt)
	return user, mapNoRows(err, auth.ErrUnauthorized)
}

func mapNoRows(err error, target error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return target
	}
	return err
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

var _ auth.UserRepository = (*Store)(nil)
var _ auth.SessionRepository = (*Store)(nil)
var _ tenants.Repository = (*Store)(nil)
var _ access.Repository = (*Store)(nil)
