package entitlements

import (
	"testing"
	"time"
)

func TestHasFeatureRequiresMatchingEnabledNonExpiredEntitlement(t *testing.T) {
	now := time.Date(2026, 5, 19, 12, 0, 0, 0, time.UTC)
	expires := now.Add(time.Hour)
	items := []Entitlement{
		{TenantID: "tenant-a", ProductKey: "planer_link", FeatureKey: "planer.pdf_export", Enabled: false},
		{TenantID: "tenant-a", ProductKey: "planer_link", FeatureKey: "planer.pdf_export", Enabled: true, ExpiresAt: &expires},
	}

	if !HasFeature(items, now, "tenant-a", "planer_link", "planer.pdf_export") {
		t.Fatal("expected matching entitlement to allow feature")
	}
	if HasFeature(items, now, "tenant-a", "infra_link", "infra.module_bacnet") {
		t.Fatal("expected unrelated entitlement to be denied")
	}
}
