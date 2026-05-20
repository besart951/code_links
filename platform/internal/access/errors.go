package access

import "errors"

var (
	ErrProductAccessDenied = errors.New("product access denied")
	ErrFeatureAccessDenied = errors.New("feature access denied")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrAudienceMismatch    = errors.New("audience mismatch")
	ErrSessionInactive     = errors.New("session inactive")
	ErrTenantMembership    = errors.New("tenant membership required")
	ErrStaleTokenVersion   = errors.New("stale token version")
	ErrStaleEntitlements   = errors.New("stale entitlements version")
)
