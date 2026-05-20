package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/besart951/code_links/platform/internal/admin"
	"github.com/besart951/code_links/platform/internal/auth"
	"github.com/jackc/pgx/v5"
)

func (s *Store) HasPermission(ctx context.Context, userID auth.UserID, permission admin.PermissionKey) (bool, error) {
	row := s.db.QueryRow(ctx, `
		select exists (
			select 1
			from platform_admin_permissions
			where user_id = $1::uuid and permission_key = $2 and revoked_at is null
		)
	`, string(userID), string(permission))
	var allowed bool
	return allowed, row.Scan(&allowed)
}

func (s *Store) Dashboard(ctx context.Context) (admin.DashboardSummary, error) {
	now := time.Now().UTC()
	tenants, err := s.count(ctx, `select count(*) from tenants`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	activeCompanies, err := s.count(ctx, `select count(*) from tenants where type = 'company' and status = 'active'`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	activeUsers, err := s.count(ctx, `select count(*) from users where status = 'active'`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	lockedUsers, err := s.count(ctx, `select count(*) from users where locked_until is not null and locked_until > now()`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	products, err := s.count(ctx, `select count(*) from products`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	activeSubscriptions, err := s.count(ctx, `select count(*) from subscriptions where status in ('active', 'trialing')`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	expiringSubscriptions, err := s.count(ctx, `select count(*) from subscriptions where status in ('active', 'trialing') and current_period_end between now() and now() + interval '30 days'`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	expiredSubscriptions, err := s.count(ctx, `select count(*) from subscriptions where status = 'expired'`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	activeSessions, err := s.count(ctx, `select count(*) from sessions where revoked_at is null and expires_at > now()`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	securityWarnings, err := s.count(ctx, `select count(*) from security_events where severity in ('warning', 'critical') and resolved_at is null`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	systemMessages, err := s.count(ctx, `select count(*) from security_events where event_type = 'system.health_warning' and resolved_at is null`)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	productSummaries, err := s.ListAdminProducts(ctx)
	if err != nil {
		return admin.DashboardSummary{}, err
	}
	return admin.DashboardSummary{
		Metrics: []admin.DashboardMetric{
			{Key: "tenants", Label: "Tenants", Value: tenants, Tone: "neutral"},
			{Key: "active_companies", Label: "Active companies", Value: activeCompanies, Tone: "neutral"},
			{Key: "active_users", Label: "Active users", Value: activeUsers, Tone: "success"},
			{Key: "locked_users", Label: "Locked users", Value: lockedUsers, Tone: toneForCount(lockedUsers)},
			{Key: "products", Label: "Products", Value: products, Tone: "neutral"},
			{Key: "active_subscriptions", Label: "Active subscriptions", Value: activeSubscriptions, Tone: "success"},
			{Key: "expiring_subscriptions", Label: "Expiring subscriptions", Value: expiringSubscriptions, Tone: toneForCount(expiringSubscriptions)},
			{Key: "expired_subscriptions", Label: "Expired subscriptions", Value: expiredSubscriptions, Tone: toneForCount(expiredSubscriptions)},
			{Key: "active_sessions", Label: "Active sessions", Value: activeSessions, Tone: "info"},
		},
		Products:           productSummaries,
		SecurityWarnings:   securityWarnings,
		OpenSystemMessages: systemMessages,
		GeneratedAt:        now,
	}, nil
}

func (s *Store) Search(ctx context.Context, query string, limit, offset int) (admin.SearchResponse, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return admin.SearchResponse{
			Query:   query,
			Results: []admin.SearchResult{},
			Facets:  map[string]int{},
			Page:    admin.PageMeta{Limit: limit, Offset: offset, Total: 0},
		}, nil
	}
	like := "%" + query + "%"
	rows, err := s.db.Query(ctx, `
		select 'tenant', id::text, name, type, array['name']::text[], 0.8
		from tenants
		where name ilike $1 or slug ilike $1
		union all
		select 'user', id::text, display_name, email, array['email', 'display_name']::text[], 0.7
		from users
		where email ilike $1 or display_name ilike $1
		order by 6 desc, 3 asc
		limit $2 offset $3
	`, like, limit, offset)
	if err != nil {
		return admin.SearchResponse{}, err
	}
	defer rows.Close()
	var results []admin.SearchResult
	facets := map[string]int{}
	for rows.Next() {
		var item admin.SearchResult
		if err := rows.Scan(&item.Type, &item.ID, &item.Title, &item.Subtitle, &item.MatchedFields, &item.Rank); err != nil {
			return admin.SearchResponse{}, err
		}
		results = append(results, item)
		facets[item.Type]++
	}
	return admin.SearchResponse{
		Query:   query,
		Results: results,
		Facets:  facets,
		Page:    admin.PageMeta{Limit: limit, Offset: offset, Total: len(results)},
	}, rows.Err()
}

func (s *Store) ListTenants(ctx context.Context, limit, offset int) (admin.Page[admin.TenantSummary], error) {
	rows, err := s.db.Query(ctx, tenantSummarySQL()+` group by t.id order by t.created_at desc limit $1 offset $2`, limit, offset)
	if err != nil {
		return admin.Page[admin.TenantSummary]{}, err
	}
	defer rows.Close()
	items, err := scanTenantSummaries(rows)
	return admin.Page[admin.TenantSummary]{Items: items, Page: admin.PageMeta{Limit: limit, Offset: offset, Total: len(items)}}, err
}

func (s *Store) GetTenant(ctx context.Context, tenantID string) (admin.TenantSummary, error) {
	row := s.db.QueryRow(ctx, tenantSummarySQL()+` and t.id = $1::uuid group by t.id limit 1`, tenantID)
	return scanTenantSummary(row)
}

func (s *Store) ListUsers(ctx context.Context, limit, offset int) (admin.Page[admin.UserSummary], error) {
	rows, err := s.db.Query(ctx, userSummarySQL()+` group by u.id order by u.created_at desc limit $1 offset $2`, limit, offset)
	if err != nil {
		return admin.Page[admin.UserSummary]{}, err
	}
	defer rows.Close()
	items, err := scanUserSummaries(rows)
	return admin.Page[admin.UserSummary]{Items: items, Page: admin.PageMeta{Limit: limit, Offset: offset, Total: len(items)}}, err
}

func (s *Store) GetUser(ctx context.Context, userID string) (admin.UserSummary, error) {
	row := s.db.QueryRow(ctx, userSummarySQL()+` and u.id = $1::uuid group by u.id limit 1`, userID)
	return scanUserSummary(row)
}

func (s *Store) ListAdminProducts(ctx context.Context) ([]admin.ProductSummary, error) {
	rows, err := s.db.Query(ctx, `
		select p.key, p.name,
			count(distinct s.tenant_id) filter (where s.status in ('active', 'trialing')) as active_tenants,
			count(distinct tm.user_id) filter (where tm.status = 'active' and s.status in ('active', 'trialing')) as active_users,
			count(distinct s.id) filter (where s.status in ('active', 'trialing')) as active_subscriptions,
			0 as warning_count,
			null::timestamptz as last_access_at
		from products p
		left join plans pl on pl.product_key = p.key
		left join subscriptions s on s.plan_id = pl.id
		left join tenant_members tm on tm.tenant_id = s.tenant_id
		group by p.key, p.name
		order by p.key
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []admin.ProductSummary
	for rows.Next() {
		var item admin.ProductSummary
		if err := rows.Scan(&item.ProductKey, &item.Name, &item.ActiveTenants, &item.ActiveUsers, &item.ActiveSubscriptions, &item.WarningCount, &item.LastAccessAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) ListSubscriptions(ctx context.Context, limit, offset int) (admin.Page[admin.SubscriptionSummary], error) {
	rows, err := s.db.Query(ctx, `
		select s.id::text, s.tenant_id::text, t.name, coalesce(pl.product_key, pl.bundle_key, ''), pl.name, s.status, s.current_period_end
		from subscriptions s
		join tenants t on t.id = s.tenant_id
		join plans pl on pl.id = s.plan_id
		order by s.created_at desc
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return admin.Page[admin.SubscriptionSummary]{}, err
	}
	defer rows.Close()
	var items []admin.SubscriptionSummary
	for rows.Next() {
		var item admin.SubscriptionSummary
		if err := rows.Scan(&item.ID, &item.TenantID, &item.TenantName, &item.ProductKey, &item.PlanName, &item.Status, &item.CurrentPeriodEnd); err != nil {
			return admin.Page[admin.SubscriptionSummary]{}, err
		}
		items = append(items, item)
	}
	return admin.Page[admin.SubscriptionSummary]{Items: items, Page: admin.PageMeta{Limit: limit, Offset: offset, Total: len(items)}}, rows.Err()
}

func (s *Store) ListAudit(ctx context.Context, limit, offset int) (admin.Page[admin.AuditLogEntry], error) {
	rows, err := s.db.Query(ctx, `
		select id::text, tenant_id::text, actor_user_id::text, target_type, target_id, action, reason, ip_address, user_agent, created_at, metadata
		from audit_logs
		order by created_at desc
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return admin.Page[admin.AuditLogEntry]{}, err
	}
	defer rows.Close()
	var items []admin.AuditLogEntry
	for rows.Next() {
		var item admin.AuditLogEntry
		var metadata []byte
		if err := rows.Scan(&item.ID, &item.TenantID, &item.ActorUserID, &item.TargetType, &item.TargetID, &item.Action, &item.Reason, &item.IPAddress, &item.UserAgent, &item.CreatedAt, &metadata); err != nil {
			return admin.Page[admin.AuditLogEntry]{}, err
		}
		item.Metadata = decodeMetadata(metadata)
		items = append(items, item)
	}
	return admin.Page[admin.AuditLogEntry]{Items: items, Page: admin.PageMeta{Limit: limit, Offset: offset, Total: len(items)}}, rows.Err()
}

func (s *Store) ListNotifications(ctx context.Context) (admin.NotificationsSummary, error) {
	templates, err := s.notificationTemplates(ctx)
	if err != nil {
		return admin.NotificationsSummary{}, err
	}
	deliveries, err := s.notificationDeliveries(ctx)
	if err != nil {
		return admin.NotificationsSummary{}, err
	}
	return admin.NotificationsSummary{Templates: templates, Deliveries: deliveries}, nil
}

func (s *Store) ListSecurity(ctx context.Context, limit, offset int) (admin.Page[admin.SecurityEventSummary], error) {
	rows, err := s.db.Query(ctx, `
		select id::text, event_type, severity, user_id::text, tenant_id::text, ip_address, created_at, summary
		from security_events
		order by created_at desc
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return admin.Page[admin.SecurityEventSummary]{}, err
	}
	defer rows.Close()
	var items []admin.SecurityEventSummary
	for rows.Next() {
		var item admin.SecurityEventSummary
		if err := rows.Scan(&item.ID, &item.EventType, &item.Severity, &item.UserID, &item.TenantID, &item.IPAddress, &item.CreatedAt, &item.Summary); err != nil {
			return admin.Page[admin.SecurityEventSummary]{}, err
		}
		items = append(items, item)
	}
	return admin.Page[admin.SecurityEventSummary]{Items: items, Page: admin.PageMeta{Limit: limit, Offset: offset, Total: len(items)}}, rows.Err()
}

func (s *Store) ListSettings(ctx context.Context, limit, offset int) (admin.Page[admin.Setting], error) {
	rows, err := s.db.Query(ctx, `
		select key, label, value_json, value_type, sensitive, requires_reason, updated_at
		from admin_settings
		order by key
		limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return admin.Page[admin.Setting]{}, err
	}
	defer rows.Close()
	var items []admin.Setting
	for rows.Next() {
		var item admin.Setting
		var value []byte
		if err := rows.Scan(&item.Key, &item.Label, &value, &item.ValueType, &item.Sensitive, &item.RequiresReason, &item.UpdatedAt); err != nil {
			return admin.Page[admin.Setting]{}, err
		}
		item.Value = decodeAny(value)
		items = append(items, item)
	}
	return admin.Page[admin.Setting]{Items: items, Page: admin.PageMeta{Limit: limit, Offset: offset, Total: len(items)}}, rows.Err()
}

func (s *Store) notificationTemplates(ctx context.Context) ([]admin.NotificationTemplateSummary, error) {
	rows, err := s.db.Query(ctx, `
		select id::text, key, channel, subject, enabled, updated_at
		from notification_templates
		order by key, channel
		limit 50
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []admin.NotificationTemplateSummary
	for rows.Next() {
		var item admin.NotificationTemplateSummary
		if err := rows.Scan(&item.ID, &item.Key, &item.Channel, &item.Subject, &item.Enabled, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) notificationDeliveries(ctx context.Context) ([]admin.NotificationDeliverySummary, error) {
	rows, err := s.db.Query(ctx, `
		select id::text, event_key, channel, status, recipient, created_at, last_attempt_at
		from notification_deliveries
		order by created_at desc
		limit 50
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []admin.NotificationDeliverySummary
	for rows.Next() {
		var item admin.NotificationDeliverySummary
		if err := rows.Scan(&item.ID, &item.EventKey, &item.Channel, &item.Status, &item.Recipient, &item.CreatedAt, &item.LastAttemptAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) count(ctx context.Context, query string) (int, error) {
	var count int
	err := s.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}

func tenantSummarySQL() string {
	return `
		select t.id::text, t.name, t.type, t.status, t.created_at, t.updated_at, t.owner_user_id::text,
			t.billing_email, t.country, t.locale, t.timezone,
			coalesce(string_agg(distinct pl.product_key, ',') filter (where s.status in ('active', 'trialing') and pl.product_key is not null), '') as active_products,
			coalesce(max(s.status), 'none') as subscription_status,
			'normal' as risk_status
		from tenants t
		left join subscriptions s on s.tenant_id = t.id
		left join plans pl on pl.id = s.plan_id
		where true`
}

func scanTenantSummaries(rows pgx.Rows) ([]admin.TenantSummary, error) {
	var items []admin.TenantSummary
	for rows.Next() {
		item, err := scanTenantSummary(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanTenantSummary(row pgx.Row) (admin.TenantSummary, error) {
	var item admin.TenantSummary
	var products string
	err := row.Scan(&item.ID, &item.Name, &item.TenantType, &item.Status, &item.CreatedAt, &item.UpdatedAt, &item.OwnerUserID, &item.BillingEmail, &item.Country, &item.Locale, &item.Timezone, &products, &item.SubscriptionStatus, &item.RiskStatus)
	if err != nil {
		return admin.TenantSummary{}, mapAdminNoRows(err)
	}
	item.ActiveProducts = splitCSV(products)
	return item, nil
}

func userSummarySQL() string {
	return `
		select u.id::text, u.email, u.display_name, u.status, u.email_verified, u.mfa_enabled, u.last_login_at,
			u.failed_login_count, u.locked_until, u.created_at,
			count(distinct tm.tenant_id) filter (where tm.status = 'active') as tenant_count,
			count(distinct s.id) filter (where s.revoked_at is null and s.expires_at > now()) as active_sessions
		from users u
		left join tenant_members tm on tm.user_id = u.id
		left join sessions s on s.user_id = u.id
		where true`
}

func scanUserSummaries(rows pgx.Rows) ([]admin.UserSummary, error) {
	var items []admin.UserSummary
	for rows.Next() {
		item, err := scanUserSummary(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanUserSummary(row pgx.Row) (admin.UserSummary, error) {
	var item admin.UserSummary
	err := row.Scan(&item.ID, &item.Email, &item.DisplayName, &item.Status, &item.EmailVerified, &item.MFAEnabled, &item.LastLoginAt, &item.FailedLoginCount, &item.LockedUntil, &item.CreatedAt, &item.TenantCount, &item.ActiveSessions)
	if err != nil {
		return admin.UserSummary{}, mapAdminNoRows(err)
	}
	return item, nil
}

func mapAdminNoRows(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return admin.ErrNotFound
	}
	return err
}

func splitCSV(value string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func decodeMetadata(data []byte) map[string]any {
	var result map[string]any
	if len(data) == 0 || json.Unmarshal(data, &result) != nil {
		return map[string]any{}
	}
	return result
}

func decodeAny(data []byte) any {
	var result any
	if len(data) == 0 || json.Unmarshal(data, &result) != nil {
		return nil
	}
	return result
}

func toneForCount(count int) string {
	if count > 0 {
		return "warning"
	}
	return "success"
}

var _ admin.Repository = (*Store)(nil)
