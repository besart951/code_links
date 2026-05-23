import { fail } from '@sveltejs/kit';
import { isRole } from '@codelinks/config/admin-access';
import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import type { UserListQuery, UserStatus } from '$lib/domain/users/types';
import { createAdminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

function parseStatus(value: string | null): UserStatus | undefined {
	return value === 'active' || value === 'disabled' || value === 'locked' ? value : undefined;
}

function parseUserQuery(url: URL): UserListQuery {
	const role = url.searchParams.get('role');

	return {
		query: url.searchParams.get('query') || undefined,
		role: role && role !== 'all' ? (role as UserListQuery['role']) : undefined,
		status: parseStatus(url.searchParams.get('status')),
		page: Number(url.searchParams.get('page') ?? 1),
		pageSize: Number(url.searchParams.get('pageSize') ?? 25),
		sort: {
			field: 'createdAt',
			direction: 'desc'
		}
	};
}

export async function load(event) {
	const admin = requireAdmin(event.locals);
	const adminContainer = createAdminContainer(event);
	const query = parseUserQuery(event.url);
	const users = await adminContainer.listUsers.execute(admin, query);

	return {
		...users,
		query,
		maskIp: !canSeeRawIpAddress(admin)
	};
}

export const actions = {
	setStatus: async (event) => {
		const admin = requireAdmin(event.locals);
		const adminContainer = createAdminContainer(event);
		const formData = await event.request.formData();
		const userId = String(formData.get('userId') ?? '');
		const status = parseStatus(String(formData.get('status') ?? ''));

		if (!userId || !status) {
			return fail(400, { message: 'Invalid status change request' });
		}

		await adminContainer.setUserStatus.execute(admin, userId, status);

		return { ok: true };
	},
	setRole: async (event) => {
		const admin = requireAdmin(event.locals);
		const adminContainer = createAdminContainer(event);
		const formData = await event.request.formData();
		const userId = String(formData.get('userId') ?? '');
		const role = String(formData.get('role') ?? '');

		if (!userId || !isRole(role)) {
			return fail(400, { message: 'Invalid role change request' });
		}

		await adminContainer.setUserRole.execute(admin, userId, role);

		return { ok: true };
	}
};
