package billing

import "time"

type Product struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Plan struct {
	ID         string  `json:"id"`
	ProductKey *string `json:"product_key,omitempty"`
	BundleKey  *string `json:"bundle_key,omitempty"`
	Name       string  `json:"name"`
	Interval   string  `json:"interval"`
	Status     string  `json:"status"`
}

type Subscription struct {
	ID                     string     `json:"id"`
	TenantID               string     `json:"tenant_id"`
	PlanID                 string     `json:"plan_id"`
	Provider               string     `json:"provider"`
	ProviderCustomerID     *string    `json:"provider_customer_id,omitempty"`
	ProviderSubscriptionID *string    `json:"provider_subscription_id,omitempty"`
	Status                 string     `json:"status"`
	CurrentPeriodStart     *time.Time `json:"current_period_start,omitempty"`
	CurrentPeriodEnd       *time.Time `json:"current_period_end,omitempty"`
	CancelAt               *time.Time `json:"cancel_at,omitempty"`
}

type FeatureLimit struct {
	TenantID   string     `json:"tenant_id"`
	ProductKey string     `json:"product_key"`
	FeatureKey string     `json:"feature_key"`
	LimitKey   string     `json:"limit_key"`
	Value      int64      `json:"value"`
	Period     string     `json:"period"`
	ResetAt    *time.Time `json:"reset_at"`
}
