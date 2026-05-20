package tenants

import "context"

type SelectTenant struct {
	Tenants Repository
}

type SelectTenantInput struct {
	UserID   UserID
	TenantID TenantID
}

type SelectTenantResult struct {
	Tenant Tenant
	Member TenantMember
}

func (uc SelectTenant) Execute(ctx context.Context, input SelectTenantInput) (SelectTenantResult, error) {
	tenant, err := uc.Tenants.FindTenant(ctx, input.TenantID)
	if err != nil {
		return SelectTenantResult{}, err
	}
	if !tenant.IsActive() {
		return SelectTenantResult{}, ErrTenantInactive
	}
	member, err := uc.Tenants.FindMember(ctx, input.UserID, input.TenantID)
	if err != nil {
		return SelectTenantResult{}, ErrTenantMembership
	}
	if !member.IsActive() {
		return SelectTenantResult{}, ErrTenantMemberInactive
	}
	return SelectTenantResult{Tenant: tenant, Member: member}, nil
}
