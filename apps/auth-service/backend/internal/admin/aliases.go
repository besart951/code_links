package admin

import "github.com/besart951/code-links/apps/auth-service/backend/internal/domain"

type AdminActor = domain.AdminActor
type AdminAuditEntry = domain.AdminAuditEntry
type AdminPermission = domain.AdminPermission
type AdminRole = domain.AdminRole
type AdminUserListItem = domain.AdminUserListItem
type AdminUserListQuery = domain.AdminUserListQuery
type AdminUserListResult = domain.AdminUserListResult
type CountStat = domain.CountStat
type DashboardStats = domain.DashboardStats
type LoginAttempt = domain.LoginAttempt
type LoginAttemptListQuery = domain.LoginAttemptListQuery
type LoginAttemptListResult = domain.LoginAttemptListResult
type ManagedUserDetail = domain.ManagedUserDetail
type Notification = domain.Notification
type RuntimeLogEntry = domain.RuntimeLogEntry
type SecurityEvent = domain.SecurityEvent
type SmtpEncryption = domain.SmtpEncryption
type SmtpSettings = domain.SmtpSettings
type User = domain.User
type UserPermissionGrant = domain.UserPermissionGrant
type UserStatus = domain.UserStatus

const (
	PermissionDashboardRead      = domain.PermissionDashboardRead
	PermissionUsersRead          = domain.PermissionUsersRead
	PermissionUsersUpdate        = domain.PermissionUsersUpdate
	PermissionUsersLock          = domain.PermissionUsersLock
	PermissionUsersChangeRole    = domain.PermissionUsersChangeRole
	PermissionAuthLogsRead       = domain.PermissionAuthLogsRead
	PermissionSecurityEventsRead = domain.PermissionSecurityEventsRead
	PermissionSMTPSettingsRead   = domain.PermissionSMTPSettingsRead
	PermissionSMTPSettingsUpdate = domain.PermissionSMTPSettingsUpdate
	PermissionNotificationsRead  = domain.PermissionNotificationsRead
	PermissionAuditEntriesRead   = domain.PermissionAuditEntriesRead
	UserStatusActive             = domain.UserStatusActive
	UserStatusLocked             = domain.UserStatusLocked
)

var errNotFound = domain.ErrNotFound
