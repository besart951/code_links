package main

import "github.com/besart951/code-links/packages/adminaccess"

func permissionsForRoles(roles []AdminRole) []AdminPermission {
	roleValues := make([]string, 0, len(roles))
	for _, role := range roles {
		roleValues = append(roleValues, string(role))
	}
	return adminPermissions(adminaccess.PermissionsForRoles(roleValues))
}

func permissionsByRole(role AdminRole) []AdminPermission {
	return adminPermissions(adminaccess.PermissionsForRole(string(role)))
}

func adminPermissions(values []string) []AdminPermission {
	permissions := make([]AdminPermission, 0, len(values))
	for _, value := range values {
		permissions = append(permissions, AdminPermission(value))
	}
	return permissions
}

func validAdminRole(role AdminRole) bool {
	return adminaccess.IsRole(string(role))
}

func primaryRole(roles []AdminRole) AdminRole {
	for _, preferred := range []AdminRole{AdminRoleAdmin, AdminRoleSupport, AdminRoleAuditor, AdminRoleUser} {
		for _, role := range roles {
			if role == preferred {
				return preferred
			}
		}
	}

	return AdminRoleUser
}

func hasAdminPermission(actor AdminActor, permission AdminPermission) bool {
	for _, grant := range actor.Permissions {
		if grant == permission {
			return true
		}
	}

	return false
}
