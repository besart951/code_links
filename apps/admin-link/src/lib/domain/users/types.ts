import type { AdminRole } from '$lib/domain/admin-access/types';

export type UserStatus = 'active' | 'disabled' | 'locked';
export type UserSortField = 'name' | 'email' | 'primaryRole' | 'status' | 'createdAt' | 'lastLoginAt';
export type SortDirection = 'asc' | 'desc';

export interface UserListQuery {
	query?: string;
	role?: AdminRole | 'user';
	status?: UserStatus;
	page: number;
	pageSize: number;
	sort: {
		field: UserSortField;
		direction: SortDirection;
	};
}

export interface UserListItem {
	id: string;
	name: string;
	email: string;
	primaryRole: AdminRole | 'user';
	status: UserStatus;
	emailVerified: boolean;
	createdAt: string;
	lastLoginAt: string | null;
	successfulLoginCount: number;
	failedLoginCount: number;
	lastKnownIpAddress: string | null;
	lastLoginCountryCode: string | null;
}

export interface UserPermissionGrant {
	role: AdminRole | 'user';
	grantedAt: string;
	grantedBy: string;
}

export interface ManagedUserDetail extends UserListItem {
	roles: UserPermissionGrant[];
	productLicenses: string[];
	knownIpAddresses: string[];
	loginCountries: string[];
	usedDevices: string[];
}
