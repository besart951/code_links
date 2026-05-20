package access

import (
	"context"
	"time"
)

type CheckProductAccess struct {
	Repo  Repository
	Clock Clock
}

type CheckProductAccessInput struct {
	UserID   UserID
	TenantID TenantID
	Product  ProductKey
}

type AccessDecision struct {
	Allowed bool
	Reason  string
}

func (uc CheckProductAccess) Execute(ctx context.Context, input CheckProductAccessInput) (AccessDecision, error) {
	member, err := uc.Repo.MemberAccess(ctx, input.UserID, input.TenantID, input.Product)
	if err != nil || member.UserID == "" {
		return AccessDecision{}, ErrTenantMembership
	}
	access, err := uc.Repo.TenantAccess(ctx, input.TenantID, input.Product)
	if err != nil {
		return AccessDecision{}, err
	}
	if !access.HasProductAccess(input.Product, uc.Clock.Now().UTC()) {
		return AccessDecision{Allowed: false, Reason: "product_access_required"}, nil
	}
	return AccessDecision{Allowed: true}, nil
}

type CheckFeatureAccess struct {
	Repo  Repository
	Clock Clock
}

type CheckFeatureAccessInput struct {
	UserID             UserID
	TenantID           TenantID
	Product            ProductKey
	Feature            FeatureKey
	RequiredPermission PermissionKey
}

func (uc CheckFeatureAccess) Execute(ctx context.Context, input CheckFeatureAccessInput) (AccessDecision, error) {
	now := uc.Clock.Now().UTC()
	tenantAccess, err := uc.Repo.TenantAccess(ctx, input.TenantID, input.Product)
	if err != nil {
		return AccessDecision{}, err
	}
	if !tenantAccess.HasProductAccess(input.Product, now) {
		return AccessDecision{Allowed: false, Reason: "product_access_required"}, nil
	}
	if !tenantAccess.CanUseFeature(input.Product, input.Feature, now) {
		return AccessDecision{Allowed: false, Reason: "feature_access_required"}, nil
	}
	if input.RequiredPermission == "" {
		return AccessDecision{Allowed: true}, nil
	}
	member, err := uc.Repo.MemberAccess(ctx, input.UserID, input.TenantID, input.Product)
	if err != nil {
		return AccessDecision{}, err
	}
	if !member.HasPermission(input.RequiredPermission) {
		return AccessDecision{Allowed: false, Reason: "permission_required"}, nil
	}
	return AccessDecision{Allowed: true}, nil
}

type IssueAccessToken struct {
	Repo      Repository
	Tokens    TokenIssuer
	Clock     Clock
	IDs       IDGenerator
	Issuer    string
	AccessTTL time.Duration
}

type IssueAccessTokenInput struct {
	SessionID SessionID
	TenantID  TenantID
	Audience  ProductKey
}

func (uc IssueAccessToken) Execute(ctx context.Context, input IssueAccessTokenInput) (IssuedToken, error) {
	now := uc.Clock.Now().UTC()
	session, err := uc.Repo.SessionAccess(ctx, input.SessionID)
	if err != nil {
		return IssuedToken{}, ErrSessionInactive
	}
	if !session.IsActive(now) {
		return IssuedToken{}, ErrSessionInactive
	}
	if !session.HasCurrentTokenVersion() {
		return IssuedToken{}, ErrStaleTokenVersion
	}

	tenantAccess, err := uc.Repo.TenantAccess(ctx, input.TenantID, input.Audience)
	if err != nil {
		return IssuedToken{}, err
	}
	if tenantAccess.EntitlementsVersion <= 0 {
		return IssuedToken{}, ErrStaleEntitlements
	}
	if !tenantAccess.HasProductAccess(input.Audience, now) {
		return IssuedToken{}, ErrProductAccessDenied
	}

	memberAccess, err := uc.Repo.MemberAccess(ctx, session.UserID, input.TenantID, input.Audience)
	if err != nil {
		return IssuedToken{}, ErrTenantMembership
	}

	jti, err := uc.IDs.NewID("token")
	if err != nil {
		return IssuedToken{}, err
	}
	expiresAt := now.Add(uc.AccessTTL)
	roles := make([]RoleKey, 0, len(memberAccess.Roles))
	for _, role := range memberAccess.Roles {
		roles = append(roles, role.Key)
	}
	permissions := make([]PermissionKey, 0, len(memberAccess.Permissions))
	for _, permission := range memberAccess.Permissions {
		permissions = append(permissions, permission.Key)
	}
	snapshot := AuthorizationSnapshot{
		Issuer:              uc.Issuer,
		Subject:             session.UserID,
		Audience:            input.Audience,
		TenantID:            input.TenantID,
		TenantType:          tenantAccess.TenantType,
		SessionID:           input.SessionID,
		TokenVersion:        session.TokenVersion,
		EntitlementsVersion: tenantAccess.EntitlementsVersion,
		Roles:               roles,
		Permissions:         permissions,
		Product: ProductSnapshot{
			ProductKey:   input.Audience,
			PlanKey:      tenantAccess.PlanKey,
			Access:       true,
			Entitlements: tenantAccess.Entitlements,
			Limits:       tenantAccess.FeatureLimits,
		},
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: expiresAt,
		JWTID:     jti,
	}
	return uc.Tokens.IssueAccessToken(ctx, snapshot)
}
