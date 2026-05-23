import type { AdminPermission, AdminRole } from '@codelinks/config/admin-access';

export type { AdminPermission, AdminRole };

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
