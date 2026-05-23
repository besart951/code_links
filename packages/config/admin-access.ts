import contract from './admin-access.json';

export const roles = contract.roles as unknown as readonly ['user', 'admin', 'support', 'auditor'];
export const adminRoles = contract.adminRoles as unknown as readonly ['admin', 'support', 'auditor'];
export const adminPermissions = contract.permissions as unknown as readonly [
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
];

export type Role = (typeof roles)[number];
export type AdminRole = (typeof adminRoles)[number];
export type AdminPermission = (typeof adminPermissions)[number];

export const rolePermissions = contract.rolePermissions as Record<AdminRole, AdminPermission[]>;

export function isRole(value: string): value is Role {
	return roles.includes(value as Role);
}

export function isAdminRole(value: string): value is AdminRole {
	return adminRoles.includes(value as AdminRole);
}
