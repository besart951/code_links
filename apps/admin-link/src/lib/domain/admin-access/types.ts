export type AdminRole = 'admin' | 'support' | 'auditor';

export type AdminPermission =
	| 'admin.dashboard.read'
	| 'admin.users.read'
	| 'admin.users.update'
	| 'admin.users.change_role'
	| 'admin.auth_logs.read'
	| 'admin.security_events.read'
	| 'admin.users.lock'
	| 'admin.smtp_settings.read'
	| 'admin.smtp_settings.update'
	| 'admin.notifications.read'
	| 'admin.audit_entries.read';

export interface AdminActor {
	id: string;
	email: string;
	name: string;
	roles: AdminRole[];
	permissions: AdminPermission[];
}

export interface AdminNavItem {
	title: string;
	url: string;
	permission: AdminPermission;
}
