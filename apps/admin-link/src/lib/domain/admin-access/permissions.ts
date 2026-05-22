import type { AdminActor, AdminPermission, AdminRole } from './types';

export const rolePermissions = {
	admin: [
		'admin.dashboard.read',
		'admin.users.read',
		'admin.users.update',
		'admin.users.lock',
		'admin.users.change_role',
		'admin.auth_logs.read',
		'admin.security_events.read',
		'admin.smtp_settings.read',
		'admin.smtp_settings.update',
		'admin.notifications.read',
		'admin.audit_entries.read'
	],
	support: [
		'admin.dashboard.read',
		'admin.users.read',
		'admin.users.update',
		'admin.users.lock',
		'admin.auth_logs.read',
		'admin.security_events.read',
		'admin.notifications.read'
	],
	auditor: [
		'admin.dashboard.read',
		'admin.auth_logs.read',
		'admin.security_events.read',
		'admin.notifications.read',
		'admin.audit_entries.read'
	]
} satisfies Record<AdminRole, AdminPermission[]>;

export function permissionsForRoles(roles: AdminRole[]) {
	return [...new Set(roles.flatMap((role) => rolePermissions[role]))];
}

export function hasPermission(actor: AdminActor, permission: AdminPermission) {
	return actor.permissions.includes(permission);
}

export function canSeeRawIpAddress(actor: AdminActor) {
	return actor.roles.includes('admin') || actor.roles.includes('support');
}
