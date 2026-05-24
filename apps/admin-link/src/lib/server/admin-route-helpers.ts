import { isRole } from '@codelinks/config/admin-access';
import type { LoginAttemptQuery } from '$lib/domain/auth-logs/types';
import type { UserListQuery, UserStatus } from '$lib/domain/users/types';

export function parseUserStatus(value: FormDataEntryValue | string | null): UserStatus | undefined {
	return value === 'active' || value === 'disabled' || value === 'locked' ? value : undefined;
}

export function parseUserListQuery(url: URL): UserListQuery {
	const role = url.searchParams.get('role');

	return {
		query: url.searchParams.get('query') || undefined,
		role: role && role !== 'all' && isRole(role) ? role : undefined,
		status: parseUserStatus(url.searchParams.get('status')),
		page: parsePositiveInteger(url.searchParams.get('page'), 1),
		pageSize: parsePositiveInteger(url.searchParams.get('pageSize'), 25),
		sort: {
			field: 'createdAt',
			direction: 'desc'
		}
	};
}

export function parseLoginAttemptQuery(url: URL): LoginAttemptQuery {
	return {
		query: url.searchParams.get('query') || undefined,
		userId: url.searchParams.get('userId') || undefined,
		success: parseOptionalBoolean(url.searchParams.get('success')),
		page: parsePositiveInteger(url.searchParams.get('page'), 1),
		pageSize: 50
	};
}

export function formString(formData: FormData, key: string) {
	return String(formData.get(key) ?? '').trim();
}

export function parseRole(value: FormDataEntryValue | string | null) {
	const role = String(value ?? '');
	return isRole(role) ? role : undefined;
}

function parsePositiveInteger(value: string | null, fallback: number) {
	const parsed = Number(value ?? fallback);
	return Number.isInteger(parsed) && parsed > 0 ? parsed : fallback;
}

function parseOptionalBoolean(value: string | null) {
	if (value === 'true') return true;
	if (value === 'false') return false;
	return undefined;
}
