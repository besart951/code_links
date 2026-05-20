package admin

import (
	"time"

	"github.com/besart951/code_links/platform/internal/auth"
)

type PermissionKey string

const PermissionRead PermissionKey = "platform.admin.read"

type PageMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

type Page[T any] struct {
	Items []T      `json:"items"`
	Page  PageMeta `json:"page"`
}

type Me struct {
	User        UserSummary `json:"user"`
	Permissions []string    `json:"permissions"`
	Superadmin  bool        `json:"superadmin"`
}

type UserSummary struct {
	ID               string     `json:"id"`
	Email            string     `json:"email"`
	DisplayName      string     `json:"display_name"`
	Status           string     `json:"status"`
	EmailVerified    bool       `json:"email_verified"`
	MFAEnabled       bool       `json:"mfa_enabled"`
	LastLoginAt      *time.Time `json:"last_login_at"`
	FailedLoginCount int        `json:"failed_login_count"`
	LockedUntil      *time.Time `json:"locked_until"`
	CreatedAt        time.Time  `json:"created_at"`
	TenantCount      int        `json:"tenant_count"`
	ActiveSessions   int        `json:"active_sessions"`
}

type TenantSummary struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	TenantType         string     `json:"tenant_type"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	OwnerUserID        string     `json:"owner_user_id"`
	BillingEmail       *string    `json:"billing_email"`
	Country            *string    `json:"country"`
	Locale             *string    `json:"locale"`
	Timezone           *string    `json:"timezone"`
	ActiveProducts     []string   `json:"active_products"`
	SubscriptionStatus string     `json:"subscription_status"`
	RiskStatus         string     `json:"risk_status"`
}

type DashboardMetric struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value int    `json:"value"`
	Tone  string `json:"tone"`
}

type ProductSummary struct {
	ProductKey          string     `json:"product_key"`
	Name                string     `json:"name"`
	ActiveTenants       int        `json:"active_tenants"`
	ActiveUsers         int        `json:"active_users"`
	ActiveSubscriptions int        `json:"active_subscriptions"`
	WarningCount        int        `json:"warning_count"`
	LastAccessAt        *time.Time `json:"last_access_at"`
}

type DashboardSummary struct {
	Metrics            []DashboardMetric `json:"metrics"`
	Products           []ProductSummary  `json:"products"`
	SecurityWarnings   int               `json:"security_warnings"`
	OpenSystemMessages int               `json:"open_system_messages"`
	GeneratedAt        time.Time         `json:"generated_at"`
}

type SearchResult struct {
	Type          string   `json:"type"`
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Subtitle      string   `json:"subtitle"`
	MatchedFields []string `json:"matched_fields"`
	Rank          float64  `json:"rank"`
}

type SearchResponse struct {
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
	Facets  map[string]int `json:"facets"`
	Page    PageMeta       `json:"page"`
}

type AuditLogEntry struct {
	ID          string         `json:"id"`
	TenantID    *string        `json:"tenant_id"`
	ActorUserID string         `json:"actor_user_id"`
	TargetType  string         `json:"target_type"`
	TargetID    string         `json:"target_id"`
	Action      string         `json:"action"`
	Reason      *string        `json:"reason"`
	IPAddress   *string        `json:"ip_address"`
	UserAgent   *string        `json:"user_agent"`
	CreatedAt   time.Time      `json:"created_at"`
	Metadata    map[string]any `json:"metadata"`
}

type NotificationTemplateSummary struct {
	ID        string    `json:"id"`
	Key       string    `json:"key"`
	Channel   string    `json:"channel"`
	Subject   string    `json:"subject"`
	Enabled   bool      `json:"enabled"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotificationDeliverySummary struct {
	ID            string     `json:"id"`
	EventKey      string     `json:"event_key"`
	Channel       string     `json:"channel"`
	Status        string     `json:"status"`
	Recipient     string     `json:"recipient"`
	CreatedAt     time.Time  `json:"created_at"`
	LastAttemptAt *time.Time `json:"last_attempt_at"`
}

type NotificationsSummary struct {
	Templates  []NotificationTemplateSummary `json:"templates"`
	Deliveries []NotificationDeliverySummary `json:"deliveries"`
}

type SecurityEventSummary struct {
	ID        string    `json:"id"`
	EventType string    `json:"event_type"`
	Severity  string    `json:"severity"`
	UserID    *string   `json:"user_id"`
	TenantID  *string   `json:"tenant_id"`
	IPAddress *string   `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
	Summary   string    `json:"summary"`
}

type Setting struct {
	Key            string    `json:"key"`
	Label          string    `json:"label"`
	Value          any       `json:"value"`
	ValueType      string    `json:"value_type"`
	Sensitive      bool      `json:"sensitive"`
	RequiresReason bool      `json:"requires_reason"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type SubscriptionSummary struct {
	ID               string     `json:"id"`
	TenantID         string     `json:"tenant_id"`
	TenantName       string     `json:"tenant_name"`
	ProductKey       string     `json:"product_key"`
	PlanName         string     `json:"plan_name"`
	Status           string     `json:"status"`
	CurrentPeriodEnd *time.Time `json:"current_period_end"`
}

func UserToSummary(user auth.User) UserSummary {
	return UserSummary{
		ID:          string(user.ID),
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Status:      string(user.Status),
		CreatedAt:   user.CreatedAt,
		LastLoginAt: user.LastLoginAt,
	}
}
