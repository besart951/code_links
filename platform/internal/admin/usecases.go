package admin

import (
	"context"
	"time"

	"github.com/besart951/code_links/platform/internal/auth"
)

const (
	defaultLimit = 50
	maxLimit     = 100
)

type ReadService struct {
	Repo  Repository
	Clock interface{ Now() time.Time }
}

func (s ReadService) Me(ctx context.Context, user auth.User) (Me, error) {
	allowed, err := s.Repo.HasPermission(ctx, user.ID, PermissionRead)
	if err != nil {
		return Me{}, err
	}
	if !allowed {
		return Me{}, ErrSuperadminRequired
	}
	return Me{
		User:        UserToSummary(user),
		Permissions: []string{string(PermissionRead)},
		Superadmin:  true,
	}, nil
}

func (s ReadService) Dashboard(ctx context.Context) (DashboardSummary, error) {
	summary, err := s.Repo.Dashboard(ctx)
	if err != nil {
		return DashboardSummary{}, err
	}
	if summary.GeneratedAt.IsZero() && s.Clock != nil {
		summary.GeneratedAt = s.Clock.Now().UTC()
	}
	return summary, nil
}

func (s ReadService) Search(ctx context.Context, query string, limit, offset int) (SearchResponse, error) {
	return s.Repo.Search(ctx, query, normalizeLimit(limit), normalizeOffset(offset))
}

func (s ReadService) ListTenants(ctx context.Context, limit, offset int) (Page[TenantSummary], error) {
	return s.Repo.ListTenants(ctx, normalizeLimit(limit), normalizeOffset(offset))
}

func (s ReadService) GetTenant(ctx context.Context, tenantID string) (TenantSummary, error) {
	return s.Repo.GetTenant(ctx, tenantID)
}

func (s ReadService) ListUsers(ctx context.Context, limit, offset int) (Page[UserSummary], error) {
	return s.Repo.ListUsers(ctx, normalizeLimit(limit), normalizeOffset(offset))
}

func (s ReadService) GetUser(ctx context.Context, userID string) (UserSummary, error) {
	return s.Repo.GetUser(ctx, userID)
}

func (s ReadService) ListProducts(ctx context.Context) ([]ProductSummary, error) {
	return s.Repo.ListAdminProducts(ctx)
}

func (s ReadService) ListSubscriptions(ctx context.Context, limit, offset int) (Page[SubscriptionSummary], error) {
	return s.Repo.ListSubscriptions(ctx, normalizeLimit(limit), normalizeOffset(offset))
}

func (s ReadService) ListAudit(ctx context.Context, limit, offset int) (Page[AuditLogEntry], error) {
	return s.Repo.ListAudit(ctx, normalizeLimit(limit), normalizeOffset(offset))
}

func (s ReadService) ListNotifications(ctx context.Context) (NotificationsSummary, error) {
	return s.Repo.ListNotifications(ctx)
}

func (s ReadService) ListSecurity(ctx context.Context, limit, offset int) (Page[SecurityEventSummary], error) {
	return s.Repo.ListSecurity(ctx, normalizeLimit(limit), normalizeOffset(offset))
}

func (s ReadService) ListSettings(ctx context.Context, limit, offset int) (Page[Setting], error) {
	return s.Repo.ListSettings(ctx, normalizeLimit(limit), normalizeOffset(offset))
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
