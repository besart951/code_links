package entitlements

import "time"

func HasFeature(items []Entitlement, now time.Time, tenantID, productKey, featureKey string) bool {
	for _, item := range items {
		if item.TenantID != tenantID || item.ProductKey != productKey || item.FeatureKey != featureKey {
			continue
		}
		if !item.Enabled {
			continue
		}
		if item.ExpiresAt != nil && !item.ExpiresAt.After(now) {
			continue
		}
		return true
	}
	return false
}
