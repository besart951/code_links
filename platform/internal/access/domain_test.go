package access

import (
	"testing"
	"time"
)

func TestTenantAccessRules(t *testing.T) {
	now := time.Date(2026, 5, 20, 10, 0, 0, 0, time.UTC)
	expired := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tenant := TenantAccess{
		TenantID:   "tenant_1",
		ProductKey: "infra_link",
		Subscriptions: []Subscription{{
			TenantID:         "tenant_1",
			ProductKey:       "infra_link",
			PlanKey:          "business",
			Status:           SubscriptionStatusActive,
			CurrentPeriodEnd: &future,
		}},
		Entitlements: []Entitlement{
			{
				TenantID:   "tenant_1",
				ProductKey: "infra_link",
				FeatureKey: "infra_link.project.read",
				Enabled:    true,
				ExpiresAt:  &future,
			},
			{
				TenantID:   "tenant_1",
				ProductKey: "infra_link",
				FeatureKey: "infra_link.project.write",
				Enabled:    true,
				ExpiresAt:  &expired,
			},
		},
		FeatureLimits: []FeatureLimit{{
			TenantID:   "tenant_1",
			ProductKey: "infra_link",
			FeatureKey: "infra_link.project.read",
			LimitKey:   "max_projects",
			Value:      100,
		}},
	}

	if !tenant.HasActiveSubscription("infra_link", now) {
		t.Fatal("expected active subscription")
	}
	if !tenant.HasProductAccess("infra_link", now) {
		t.Fatal("expected product access from active subscription")
	}
	if !tenant.CanUseFeature("infra_link", "infra_link.project.read", now) {
		t.Fatal("expected active feature entitlement")
	}
	if tenant.CanUseFeature("infra_link", "infra_link.project.write", now) {
		t.Fatal("expired entitlement should not grant feature access")
	}
	limit, ok := tenant.LimitFor("infra_link", "infra_link.project.read", "max_projects")
	if !ok || limit.Value != 100 {
		t.Fatalf("expected max_projects limit 100, got %#v ok=%v", limit, ok)
	}
}

func TestMemberAccessHasPermission(t *testing.T) {
	member := MemberAccess{
		UserID:   "user_1",
		TenantID: "tenant_1",
		Roles: []Role{{
			Key: "admin",
		}},
		Permissions: []Permission{
			{Key: "infra_link.project.read"},
			{Key: "infra_link.project.write"},
		},
	}

	if !member.HasRole("admin") {
		t.Fatal("expected admin role")
	}
	if !member.HasPermission("infra_link.project.write") {
		t.Fatal("expected write permission")
	}
	if member.HasPermission("infra_link.field_device.manage") {
		t.Fatal("unexpected field device permission")
	}
}
