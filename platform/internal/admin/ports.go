package admin

import (
	"context"

	"github.com/besart951/code_links/platform/internal/auth"
)

type Repository interface {
	HasPermission(ctx context.Context, userID auth.UserID, permission PermissionKey) (bool, error)
	Dashboard(ctx context.Context) (DashboardSummary, error)
	Search(ctx context.Context, query string, limit, offset int) (SearchResponse, error)
	ListTenants(ctx context.Context, limit, offset int) (Page[TenantSummary], error)
	GetTenant(ctx context.Context, tenantID string) (TenantSummary, error)
	ListUsers(ctx context.Context, limit, offset int) (Page[UserSummary], error)
	GetUser(ctx context.Context, userID string) (UserSummary, error)
	ListAdminProducts(ctx context.Context) ([]ProductSummary, error)
	ListSubscriptions(ctx context.Context, limit, offset int) (Page[SubscriptionSummary], error)
	ListAudit(ctx context.Context, limit, offset int) (Page[AuditLogEntry], error)
	ListNotifications(ctx context.Context) (NotificationsSummary, error)
	ListSecurity(ctx context.Context, limit, offset int) (Page[SecurityEventSummary], error)
	ListSettings(ctx context.Context, limit, offset int) (Page[Setting], error)
}
