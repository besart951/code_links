package admin

import "github.com/besart951/code-links/apps/auth-service/backend/internal/domain"

func permissionsForRoles(roles []AdminRole) []AdminPermission {
	return PermissionsForRoles(roles)
}

func permissionsByRole(role AdminRole) []AdminPermission {
	return PermissionsByRole(role)
}

func validAdminRole(role AdminRole) bool {
	return ValidAdminRole(role)
}

func primaryRole(roles []AdminRole) AdminRole {
	return PrimaryRole(roles)
}

func hasAdminPermission(actor AdminActor, permission AdminPermission) bool {
	return HasPermission(actor, permission)
}

func PermissionsForRoles(roles []AdminRole) []AdminPermission {
	return domain.PermissionsForRoles(roles)
}
func PermissionsByRole(role AdminRole) []AdminPermission { return domain.PermissionsByRole(role) }
func ValidAdminRole(role AdminRole) bool                 { return domain.ValidAdminRole(role) }
func PrimaryRole(roles []AdminRole) AdminRole            { return domain.PrimaryRole(roles) }
func HasPermission(actor AdminActor, permission AdminPermission) bool {
	return domain.HasAdminPermission(actor, permission)
}
