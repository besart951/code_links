package tenants

import "time"

type Tenant struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	OwnerUserID  string    `json:"owner_user_id"`
	Status       string    `json:"status"`
	BillingEmail *string   `json:"billing_email,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	RoleKey      string    `json:"role_key,omitempty"`
}

type TenantMember struct {
	TenantID string    `json:"tenant_id"`
	UserID   string    `json:"user_id"`
	RoleID   string    `json:"role_id"`
	Status   string    `json:"status"`
	JoinedAt time.Time `json:"joined_at"`
}
