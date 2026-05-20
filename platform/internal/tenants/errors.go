package tenants

import "errors"

var (
	ErrTenantNotFound       = errors.New("tenant not found")
	ErrTenantInactive       = errors.New("tenant inactive")
	ErrTenantMembership     = errors.New("tenant membership required")
	ErrTenantMemberInactive = errors.New("tenant member inactive")
)
