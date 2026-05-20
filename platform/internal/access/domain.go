package access

import "time"

type UserID string
type TenantID string
type SessionID string
type RoleKey string
type PermissionKey string
type ProductKey string
type PlanKey string
type FeatureKey string
type LimitKey string
type TokenVersion int
type EntitlementsVersion int

const ProductAccessFeature FeatureKey = "product.access"

type SessionAccess struct {
	ID               SessionID
	UserID           UserID
	TokenVersion     TokenVersion
	UserTokenVersion TokenVersion
	RevokedAt        *time.Time
	ExpiresAt        time.Time
}

func (s SessionAccess) IsActive(now time.Time) bool {
	return s.RevokedAt == nil && s.ExpiresAt.After(now)
}

func (s SessionAccess) HasCurrentTokenVersion() bool {
	return s.TokenVersion > 0 && s.TokenVersion == s.UserTokenVersion
}

type TenantType string

const (
	TenantTypePersonal TenantType = "personal"
	TenantTypeTeam     TenantType = "team"
	TenantTypeCompany  TenantType = "company"
	TenantTypeMandate  TenantType = "mandate"
)

type Role struct {
	Key        RoleKey
	Name       string
	ProductKey ProductKey
}

type Permission struct {
	Key         PermissionKey
	Description string
	ProductKey  ProductKey
}

type MemberAccess struct {
	UserID      UserID
	TenantID    TenantID
	Roles       []Role
	Permissions []Permission
}

func (m MemberAccess) HasRole(role RoleKey) bool {
	for _, item := range m.Roles {
		if item.Key == role {
			return true
		}
	}
	return false
}

func (m MemberAccess) HasPermission(permission PermissionKey) bool {
	for _, item := range m.Permissions {
		if item.Key == permission {
			return true
		}
	}
	return false
}

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
)

type Product struct {
	Key    ProductKey
	Name   string
	Status ProductStatus
}

func (p Product) IsActive() bool {
	return p.Status == ProductStatusActive
}

type PlanStatus string

const (
	PlanStatusActive   PlanStatus = "active"
	PlanStatusInactive PlanStatus = "inactive"
)

type Plan struct {
	Key        PlanKey
	ProductKey ProductKey
	BundleKey  string
	Name       string
	Status     PlanStatus
}

func (p Plan) IsActive() bool {
	return p.Status == PlanStatusActive
}

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusTrialing SubscriptionStatus = "trialing"
	SubscriptionStatusCanceled SubscriptionStatus = "canceled"
	SubscriptionStatusExpired  SubscriptionStatus = "expired"
)

type Subscription struct {
	ID                 string
	TenantID           TenantID
	ProductKey         ProductKey
	PlanKey            PlanKey
	Status             SubscriptionStatus
	CurrentPeriodStart *time.Time
	CurrentPeriodEnd   *time.Time
	CancelAt           *time.Time
}

func (s Subscription) IsActiveAt(now time.Time) bool {
	if s.Status != SubscriptionStatusActive && s.Status != SubscriptionStatusTrialing {
		return false
	}
	if s.CurrentPeriodStart != nil && s.CurrentPeriodStart.After(now) {
		return false
	}
	if s.CurrentPeriodEnd != nil && !s.CurrentPeriodEnd.After(now) {
		return false
	}
	return true
}

func (s Subscription) IsActiveFor(product ProductKey, now time.Time) bool {
	return s.ProductKey == product && s.IsActiveAt(now)
}

type EntitlementSource string

const (
	EntitlementSourceSubscription EntitlementSource = "subscription"
	EntitlementSourceManual       EntitlementSource = "manual"
	EntitlementSourceTrial        EntitlementSource = "trial"
)

type Entitlement struct {
	TenantID   TenantID
	ProductKey ProductKey
	FeatureKey FeatureKey
	Source     EntitlementSource
	Enabled    bool
	ExpiresAt  *time.Time
}

func (e Entitlement) IsActiveAt(now time.Time) bool {
	return e.Enabled && (e.ExpiresAt == nil || e.ExpiresAt.After(now))
}

func (e Entitlement) GrantsFeature(product ProductKey, feature FeatureKey, now time.Time) bool {
	return e.ProductKey == product && e.FeatureKey == feature && e.IsActiveAt(now)
}

func (e Entitlement) GrantsProductAccess(product ProductKey, now time.Time) bool {
	return e.ProductKey == product && e.FeatureKey == ProductAccessFeature && e.IsActiveAt(now)
}

type FeatureLimit struct {
	TenantID   TenantID
	ProductKey ProductKey
	FeatureKey FeatureKey
	LimitKey   LimitKey
	Value      int64
	Period     string
	ResetAt    *time.Time
}

type TenantAccess struct {
	TenantID            TenantID
	TenantType          TenantType
	ProductKey          ProductKey
	PlanKey             PlanKey
	Subscriptions       []Subscription
	Entitlements        []Entitlement
	FeatureLimits       []FeatureLimit
	EntitlementsVersion EntitlementsVersion
}

func (t TenantAccess) HasActiveSubscription(product ProductKey, now time.Time) bool {
	for _, item := range t.Subscriptions {
		if item.IsActiveFor(product, now) {
			return true
		}
	}
	return false
}

func (t TenantAccess) HasProductAccess(product ProductKey, now time.Time) bool {
	if t.HasActiveSubscription(product, now) {
		return true
	}
	for _, item := range t.Entitlements {
		if item.ProductKey == product && item.IsActiveAt(now) {
			return true
		}
	}
	return false
}

func (t TenantAccess) CanUseFeature(product ProductKey, feature FeatureKey, now time.Time) bool {
	for _, item := range t.Entitlements {
		if item.GrantsFeature(product, feature, now) {
			return true
		}
	}
	return false
}

func (t TenantAccess) LimitFor(product ProductKey, feature FeatureKey, limit LimitKey) (FeatureLimit, bool) {
	for _, item := range t.FeatureLimits {
		if item.ProductKey == product && item.FeatureKey == feature && item.LimitKey == limit {
			return item, true
		}
	}
	return FeatureLimit{}, false
}

type ProductSnapshot struct {
	ProductKey   ProductKey
	PlanKey      PlanKey
	Access       bool
	Entitlements []Entitlement
	Limits       []FeatureLimit
}

type AuthorizationSnapshot struct {
	Issuer              string
	Subject             UserID
	Audience            ProductKey
	TenantID            TenantID
	TenantType          TenantType
	SessionID           SessionID
	TokenVersion        TokenVersion
	EntitlementsVersion EntitlementsVersion
	Roles               []RoleKey
	Permissions         []PermissionKey
	Product             ProductSnapshot
	IssuedAt            time.Time
	NotBefore           time.Time
	ExpiresAt           time.Time
	JWTID               string
}
