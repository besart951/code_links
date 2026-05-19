package entitlements

import "time"

type Entitlement struct {
	TenantID   string     `json:"tenant_id"`
	ProductKey string     `json:"product_key"`
	FeatureKey string     `json:"feature_key"`
	Source     string     `json:"source"`
	Enabled    bool       `json:"enabled"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

type AuthorizeRequest struct {
	UserID     string `json:"user_id"`
	TenantID   string `json:"tenant_id"`
	ProductKey string `json:"product_key"`
	FeatureKey string `json:"feature_key"`
}

type AuthorizeResponse struct {
	Allowed bool    `json:"allowed"`
	Reason  *string `json:"reason,omitempty"`
}
