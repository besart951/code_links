package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/besart951/code-links/packages/adminaccess"
	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
	UserStatusLocked   UserStatus = "locked"
)

type AdminRole string

const (
	AdminRoleAdmin   AdminRole = adminaccess.AdminRoleAdmin
	AdminRoleSupport AdminRole = adminaccess.AdminRoleSupport
	AdminRoleAuditor AdminRole = adminaccess.AdminRoleAuditor
	AdminRoleUser    AdminRole = adminaccess.RoleUser
)

type AdminPermission string

const (
	PermissionDashboardRead      AdminPermission = adminaccess.PermissionDashboardRead
	PermissionUsersRead          AdminPermission = adminaccess.PermissionUsersRead
	PermissionUsersUpdate        AdminPermission = adminaccess.PermissionUsersUpdate
	PermissionUsersLock          AdminPermission = adminaccess.PermissionUsersLock
	PermissionUsersChangeRole    AdminPermission = adminaccess.PermissionUsersChangeRole
	PermissionAuthLogsRead       AdminPermission = adminaccess.PermissionAuthLogsRead
	PermissionSecurityEventsRead AdminPermission = adminaccess.PermissionSecurityEventsRead
	PermissionSMTPSettingsRead   AdminPermission = adminaccess.PermissionSMTPSettingsRead
	PermissionSMTPSettingsUpdate AdminPermission = adminaccess.PermissionSMTPSettingsUpdate
	PermissionNotificationsRead  AdminPermission = adminaccess.PermissionNotificationsRead
	PermissionAuditEntriesRead   AdminPermission = adminaccess.PermissionAuditEntriesRead
)

type User struct {
	ID                   uuid.UUID  `json:"id"`
	Email                string     `json:"email"`
	Name                 string     `json:"name"`
	PasswordHash         string     `json:"-"`
	Status               UserStatus `json:"status"`
	EmailVerifiedAt      *time.Time `json:"emailVerifiedAt"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	LastLoginAt          *time.Time `json:"lastLoginAt"`
	LastLoginIP          string     `json:"lastLoginIpAddress"`
	LastLoginCountryCode string     `json:"lastLoginCountryCode"`
}

type AdminUserListQuery struct {
	Query     string
	Role      string
	Status    UserStatus
	Page      int
	PageSize  int
	Sort      string
	Direction string
}

type AdminUserListResult struct {
	Items    []AdminUserListItem `json:"items"`
	Total    int                 `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}

type AdminUserListItem struct {
	ID                   string     `json:"id"`
	Name                 string     `json:"name"`
	Email                string     `json:"email"`
	PrimaryRole          AdminRole  `json:"primaryRole"`
	Status               UserStatus `json:"status"`
	EmailVerified        bool       `json:"emailVerified"`
	CreatedAt            time.Time  `json:"createdAt"`
	LastLoginAt          *time.Time `json:"lastLoginAt"`
	SuccessfulLoginCount int        `json:"successfulLoginCount"`
	FailedLoginCount     int        `json:"failedLoginCount"`
	LastKnownIPAddress   *string    `json:"lastKnownIpAddress"`
	LastLoginCountryCode *string    `json:"lastLoginCountryCode"`
}

type UserPermissionGrant struct {
	Role      AdminRole `json:"role"`
	GrantedAt time.Time `json:"grantedAt"`
	GrantedBy string    `json:"grantedBy"`
}

type ManagedUserDetail struct {
	AdminUserListItem
	Roles            []UserPermissionGrant `json:"roles"`
	ProductLicenses  []string              `json:"productLicenses"`
	KnownIPAddresses []string              `json:"knownIpAddresses"`
	LoginCountries   []string              `json:"loginCountries"`
	UsedDevices      []string              `json:"usedDevices"`
}

type AdminActor struct {
	ID          string            `json:"id"`
	Email       string            `json:"email"`
	Name        string            `json:"name"`
	Roles       []AdminRole       `json:"roles"`
	Permissions []AdminPermission `json:"permissions"`
}

type UserSnapshot struct {
	ID            string            `json:"id"`
	Email         string            `json:"email"`
	Name          string            `json:"name"`
	Status        UserStatus        `json:"status"`
	EmailVerified bool              `json:"emailVerified"`
	Licenses      []string          `json:"licenses"`
	Roles         []AdminRole       `json:"roles"`
	Permissions   []AdminPermission `json:"permissions"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

type UserCard struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LoginFailureReason string

const (
	LoginFailureWrongPassword     LoginFailureReason = "wrong_password"
	LoginFailureUnknownEmail      LoginFailureReason = "unknown_email"
	LoginFailureAccountLocked     LoginFailureReason = "account_locked"
	LoginFailureTooManyAttempts   LoginFailureReason = "too_many_attempts"
	LoginFailureInvalidToken      LoginFailureReason = "invalid_token"
	LoginFailureEmailNotConfirmed LoginFailureReason = "email_not_confirmed"
)

type LoginAttempt struct {
	ID              uuid.UUID           `json:"id"`
	UserID          *uuid.UUID          `json:"userId"`
	EmailAttempted  string              `json:"emailAttempted"`
	OccurredAt      time.Time           `json:"occurredAt"`
	IPAddress       string              `json:"ipAddress"`
	IPHash          string              `json:"-"`
	CountryCode     string              `json:"countryCode"`
	City            string              `json:"city"`
	UserAgent       string              `json:"userAgent"`
	Browser         string              `json:"browser"`
	OperatingSystem string              `json:"operatingSystem"`
	Success         bool                `json:"success"`
	FailureReason   *LoginFailureReason `json:"failureReason"`
	AuthMethod      string              `json:"authMethod"`
	RiskScore       int                 `json:"riskScore"`
	CorrelationID   string              `json:"correlationId"`
}

type LoginAttemptListQuery struct {
	UserID   *uuid.UUID
	Success  *bool
	Query    string
	Page     int
	PageSize int
}

type LoginAttemptListResult struct {
	Items    []LoginAttempt `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

type SecurityEvent struct {
	ID              uuid.UUID  `json:"id"`
	UserID          *uuid.UUID `json:"userId"`
	Type            string     `json:"type"`
	Severity        string     `json:"severity"`
	Status          string     `json:"status"`
	Summary         string     `json:"summary"`
	DetectedAt      time.Time  `json:"detectedAt"`
	ResolvedAt      *time.Time `json:"resolvedAt"`
	SourceIPAddress string     `json:"sourceIpAddress"`
	CountryCode     string     `json:"countryCode"`
}

type SmtpEncryption string

const (
	SmtpEncryptionNone     SmtpEncryption = "none"
	SmtpEncryptionSSL      SmtpEncryption = "ssl"
	SmtpEncryptionTLS      SmtpEncryption = "tls"
	SmtpEncryptionStartTLS SmtpEncryption = "starttls"
)

type SmtpSettings struct {
	Host              string         `json:"host"`
	Port              int            `json:"port"`
	Username          string         `json:"username"`
	PasswordEncrypted string         `json:"-"`
	HasPassword       bool           `json:"hasPassword"`
	Encryption        SmtpEncryption `json:"encryption"`
	FromEmail         string         `json:"fromEmail"`
	FromName          string         `json:"fromName"`
	ReplyToEmail      string         `json:"replyToEmail"`
	Active            bool           `json:"active"`
	UpdatedAt         time.Time      `json:"updatedAt"`
}

type Notification struct {
	ID        uuid.UUID  `json:"id"`
	Type      string     `json:"type"`
	Channel   string     `json:"channel"`
	Recipient string     `json:"recipient"`
	Subject   string     `json:"subject"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	SentAt    *time.Time `json:"sentAt"`
}

type AdminAuditEntry struct {
	ID          uuid.UUID `json:"id"`
	ActorUserID uuid.UUID `json:"actorUserId"`
	Action      string    `json:"action"`
	TargetType  string    `json:"targetType"`
	TargetID    string    `json:"targetId"`
	Reason      string    `json:"reason"`
	IPAddress   string    `json:"ipAddress"`
	CreatedAt   time.Time `json:"createdAt"`
}

type RuntimeLogEntry struct {
	ID         string    `json:"id"`
	OccurredAt time.Time `json:"occurredAt"`
	Level      string    `json:"level"`
	Source     string    `json:"source"`
	Message    string    `json:"message"`
	Raw        string    `json:"raw"`
}

type DashboardStats struct {
	Users struct {
		Total         int `json:"total"`
		Active        int `json:"active"`
		Locked        int `json:"locked"`
		NewLast7Days  int `json:"newLast7Days"`
		NewLast30Days int `json:"newLast30Days"`
	} `json:"users"`
	LoginAttempts struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"loginAttempts"`
	PasswordResetRequests int         `json:"passwordResetRequests"`
	Notifications         int         `json:"notifications"`
	OpenSecurityEvents    int         `json:"openSecurityEvents"`
	TopCountries          []CountStat `json:"topCountries"`
	TopIPAddresses        []CountStat `json:"topIpAddresses"`
}

type CountStat struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

var (
	ErrNotFound  = errors.New("not found")
	ErrConflict  = errors.New("conflict")
	ErrForbidden = errors.New("forbidden")
)

func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func PermissionsForRoles(roles []AdminRole) []AdminPermission {
	roleValues := make([]string, 0, len(roles))
	for _, role := range roles {
		roleValues = append(roleValues, string(role))
	}
	return adminPermissions(adminaccess.PermissionsForRoles(roleValues))
}

func PermissionsByRole(role AdminRole) []AdminPermission {
	return adminPermissions(adminaccess.PermissionsForRole(string(role)))
}

func adminPermissions(values []string) []AdminPermission {
	permissions := make([]AdminPermission, 0, len(values))
	for _, value := range values {
		permissions = append(permissions, AdminPermission(value))
	}
	return permissions
}

func ValidAdminRole(role AdminRole) bool {
	return adminaccess.IsRole(string(role))
}

func ValidUserStatus(status UserStatus) bool {
	return status == UserStatusActive || status == UserStatusDisabled || status == UserStatusLocked
}

func PrimaryRole(roles []AdminRole) AdminRole {
	for _, preferred := range []AdminRole{AdminRoleAdmin, AdminRoleSupport, AdminRoleAuditor, AdminRoleUser} {
		for _, role := range roles {
			if role == preferred {
				return preferred
			}
		}
	}

	return AdminRoleUser
}

func HasAdminPermission(actor AdminActor, permission AdminPermission) bool {
	for _, grant := range actor.Permissions {
		if grant == permission {
			return true
		}
	}

	return false
}
