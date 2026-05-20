package access

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCheckFeatureAccessRequiresFeatureAndPermission(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	repo := fakeAccessRepo{
		tenant: TenantAccess{
			TenantID:            "tenant_1",
			TenantType:          TenantTypeCompany,
			ProductKey:          "infra_link",
			EntitlementsVersion: 12,
			Entitlements: []Entitlement{
				{TenantID: "tenant_1", ProductKey: "infra_link", FeatureKey: ProductAccessFeature, Enabled: true},
				{TenantID: "tenant_1", ProductKey: "infra_link", FeatureKey: "infra_link.project.write", Enabled: true},
			},
		},
		member: MemberAccess{
			UserID:      "user_1",
			TenantID:    "tenant_1",
			Permissions: []Permission{{Key: "infra_link.project.read"}},
		},
	}
	uc := CheckFeatureAccess{Repo: repo, Clock: fixedClock{now: now}}

	decision, err := uc.Execute(context.Background(), CheckFeatureAccessInput{
		UserID:             "user_1",
		TenantID:           "tenant_1",
		Product:            "infra_link",
		Feature:            "infra_link.project.write",
		RequiredPermission: "infra_link.project.write",
	})
	if err != nil {
		t.Fatal(err)
	}
	if decision.Allowed || decision.Reason != "permission_required" {
		t.Fatalf("expected permission denial, got %#v", decision)
	}

	repo.member.Permissions = append(repo.member.Permissions, Permission{Key: "infra_link.project.write"})
	uc.Repo = repo
	decision, err = uc.Execute(context.Background(), CheckFeatureAccessInput{
		UserID:             "user_1",
		TenantID:           "tenant_1",
		Product:            "infra_link",
		Feature:            "infra_link.project.write",
		RequiredPermission: "infra_link.project.write",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !decision.Allowed {
		t.Fatalf("expected access, got %#v", decision)
	}
}

func TestCheckProductAccessRequiresMembership(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	uc := CheckProductAccess{
		Repo: fakeAccessRepo{
			tenant: TenantAccess{
				TenantID:            "tenant_1",
				ProductKey:          "infra_link",
				EntitlementsVersion: 12,
				Entitlements:        []Entitlement{{TenantID: "tenant_1", ProductKey: "infra_link", FeatureKey: ProductAccessFeature, Enabled: true}},
			},
			memberErr: ErrTenantMembership,
		},
		Clock: fixedClock{now: now},
	}
	_, err := uc.Execute(context.Background(), CheckProductAccessInput{
		UserID:   "user_1",
		TenantID: "tenant_1",
		Product:  "infra_link",
	})
	if !errors.Is(err, ErrTenantMembership) {
		t.Fatalf("expected membership error, got %v", err)
	}
}

func TestIssueAccessTokenBuildsSnapshotFromServerSideState(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	repo := fakeAccessRepo{
		session: SessionAccess{
			ID:               "session_1",
			UserID:           "user_1",
			TokenVersion:     3,
			UserTokenVersion: 3,
			ExpiresAt:        now.Add(time.Hour),
		},
		tenant: TenantAccess{
			TenantID:            "tenant_1",
			TenantType:          TenantTypeCompany,
			ProductKey:          "infra_link",
			PlanKey:             "business",
			EntitlementsVersion: 12,
			Entitlements: []Entitlement{
				{TenantID: "tenant_1", ProductKey: "infra_link", FeatureKey: ProductAccessFeature, Enabled: true},
			},
		},
		member: MemberAccess{
			UserID:      "user_1",
			TenantID:    "tenant_1",
			Roles:       []Role{{Key: "owner"}},
			Permissions: []Permission{{Key: "infra_link.project.read"}},
		},
	}
	tokens := &capturingTokenIssuer{}
	uc := IssueAccessToken{
		Repo:      repo,
		Tokens:    tokens,
		Clock:     fixedClock{now: now},
		IDs:       fixedIDs{id: "token_1"},
		Issuer:    "https://auth.codelinks.ch",
		AccessTTL: 10 * time.Minute,
	}

	issued, err := uc.Execute(context.Background(), IssueAccessTokenInput{
		SessionID: "session_1",
		TenantID:  "tenant_1",
		Audience:  "infra_link",
	})
	if err != nil {
		t.Fatal(err)
	}
	if issued.Value != "nested-token" || issued.JWTID != "token_1" {
		t.Fatalf("unexpected token result %#v", issued)
	}
	if tokens.snapshot.Subject != "user_1" || tokens.snapshot.TokenVersion != 3 || tokens.snapshot.EntitlementsVersion != 12 {
		t.Fatalf("snapshot did not derive server-side state: %#v", tokens.snapshot)
	}
	if len(tokens.snapshot.Permissions) != 1 || tokens.snapshot.Permissions[0] != "infra_link.project.read" {
		t.Fatalf("snapshot did not include member permissions: %#v", tokens.snapshot)
	}
}

func TestIssueAccessTokenDenials(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	base := fakeAccessRepo{
		session: SessionAccess{ID: "session_1", UserID: "user_1", TokenVersion: 3, UserTokenVersion: 3, ExpiresAt: now.Add(time.Hour)},
		tenant: TenantAccess{
			TenantID:            "tenant_1",
			TenantType:          TenantTypeCompany,
			ProductKey:          "infra_link",
			EntitlementsVersion: 12,
			Entitlements:        []Entitlement{{TenantID: "tenant_1", ProductKey: "infra_link", FeatureKey: ProductAccessFeature, Enabled: true}},
		},
		member: MemberAccess{UserID: "user_1", TenantID: "tenant_1"},
	}
	tests := []struct {
		name string
		repo fakeAccessRepo
		want error
	}{
		{
			name: "inactive session",
			repo: func() fakeAccessRepo {
				r := base
				r.session.ExpiresAt = now.Add(-time.Minute)
				return r
			}(),
			want: ErrSessionInactive,
		},
		{
			name: "stale token version",
			repo: func() fakeAccessRepo {
				r := base
				r.session.UserTokenVersion = 4
				return r
			}(),
			want: ErrStaleTokenVersion,
		},
		{
			name: "no product access",
			repo: func() fakeAccessRepo {
				r := base
				r.tenant.Entitlements = nil
				return r
			}(),
			want: ErrProductAccessDenied,
		},
		{
			name: "stale entitlement version",
			repo: func() fakeAccessRepo {
				r := base
				r.tenant.EntitlementsVersion = 0
				return r
			}(),
			want: ErrStaleEntitlements,
		},
		{
			name: "missing membership",
			repo: func() fakeAccessRepo {
				r := base
				r.memberErr = ErrTenantMembership
				return r
			}(),
			want: ErrTenantMembership,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := IssueAccessToken{
				Repo:      tt.repo,
				Tokens:    &capturingTokenIssuer{},
				Clock:     fixedClock{now: now},
				IDs:       fixedIDs{id: "token_1"},
				Issuer:    "https://auth.codelinks.ch",
				AccessTTL: 10 * time.Minute,
			}
			_, err := uc.Execute(context.Background(), IssueAccessTokenInput{
				SessionID: "session_1",
				TenantID:  "tenant_1",
				Audience:  "infra_link",
			})
			if !errors.Is(err, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}

type fakeAccessRepo struct {
	session   SessionAccess
	tenant    TenantAccess
	member    MemberAccess
	memberErr error
}

func (r fakeAccessRepo) SessionAccess(ctx context.Context, sessionID SessionID) (SessionAccess, error) {
	return r.session, nil
}

func (r fakeAccessRepo) TenantAccess(ctx context.Context, tenantID TenantID, product ProductKey) (TenantAccess, error) {
	return r.tenant, nil
}

func (r fakeAccessRepo) MemberAccess(ctx context.Context, userID UserID, tenantID TenantID, product ProductKey) (MemberAccess, error) {
	if r.memberErr != nil {
		return MemberAccess{}, r.memberErr
	}
	return r.member, nil
}

type capturingTokenIssuer struct {
	snapshot AuthorizationSnapshot
}

func (i *capturingTokenIssuer) IssueAccessToken(ctx context.Context, snapshot AuthorizationSnapshot) (IssuedToken, error) {
	i.snapshot = snapshot
	return IssuedToken{Value: "nested-token", ExpiresAt: snapshot.ExpiresAt, JWTID: snapshot.JWTID}, nil
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}

type fixedIDs struct {
	id string
}

func (g fixedIDs) NewID(prefix string) (string, error) {
	return g.id, nil
}
