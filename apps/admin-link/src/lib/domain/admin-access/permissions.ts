import type { AdminActor, AdminPermission, AdminRole } from './types';
import { rolePermissions as sharedRolePermissions } from '@codelinks/config/admin-access';

export const rolePermissions = sharedRolePermissions satisfies Record<AdminRole, AdminPermission[]>;

export function permissionsForRoles(roles: AdminRole[]) {
	return [...new Set(roles.flatMap((role) => rolePermissions[role]))];
}

export function hasPermission(actor: AdminActor, permission: AdminPermission) {
	return actor.permissions.includes(permission);
}

export function canSeeRawIpAddress(actor: AdminActor) {
	return hasPermission(actor, 'admin.users.update');
}
