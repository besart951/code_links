package postgres

import "github.com/besart951/code-links/apps/auth-service/backend/internal/domain"

type User = domain.User
type UserCard = domain.UserCard
type UserStatus = domain.UserStatus
type AdminRole = domain.AdminRole
type AdminPermission = domain.AdminPermission
type AdminActor = domain.AdminActor
type LoginFailureReason = domain.LoginFailureReason
type LoginAttempt = domain.LoginAttempt
type LoginAttemptListQuery = domain.LoginAttemptListQuery
type LoginAttemptListResult = domain.LoginAttemptListResult
type SecurityEvent = domain.SecurityEvent
type SmtpSettings = domain.SmtpSettings
type Notification = domain.Notification
type AdminAuditEntry = domain.AdminAuditEntry
type AdminUserListQuery = domain.AdminUserListQuery
type AdminUserListResult = domain.AdminUserListResult
type AdminUserListItem = domain.AdminUserListItem
type ManagedUserDetail = domain.ManagedUserDetail
type UserPermissionGrant = domain.UserPermissionGrant
type DashboardStats = domain.DashboardStats
type CountStat = domain.CountStat

const (
	SmtpEncryptionStartTLS = domain.SmtpEncryptionStartTLS
)

var (
	errNotFound  = domain.ErrNotFound
	errConflict  = domain.ErrConflict
	errForbidden = domain.ErrForbidden
)

func permissionsForRoles(roles []AdminRole) []AdminPermission {
	return domain.PermissionsForRoles(roles)
}
func primaryRole(roles []AdminRole) AdminRole { return domain.PrimaryRole(roles) }
func validUserStatus(status UserStatus) bool  { return domain.ValidUserStatus(status) }
func validAdminRole(role AdminRole) bool      { return domain.ValidAdminRole(role) }

func normalizedPage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func emptyStringToNil(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
