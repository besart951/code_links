package tenants

import "context"

type Repository interface {
	FindTenant(ctx context.Context, tenantID TenantID) (Tenant, error)
	FindMember(ctx context.Context, userID UserID, tenantID TenantID) (TenantMember, error)
	ListTenantsForUser(ctx context.Context, userID UserID) ([]Tenant, error)
}
