import { fail } from '@sveltejs/kit';
import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import {
	formString,
	adminActionFailure,
	adminLoad,
	parseRole,
	parseUserListQuery,
	parseUserStatus
} from '$lib/server/admin-route-helpers';
import { createAdminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

export async function load(event) {
	return adminLoad(async () => {
		const admin = requireAdmin(event.locals);
		const adminContainer = createAdminContainer(event);
		const query = parseUserListQuery(event.url);
		const users = await adminContainer.listUsers.execute(admin, query);

		return {
			...users,
			query,
			maskIp: !canSeeRawIpAddress(admin)
		};
	});
}

export const actions = {
	setStatus: async (event) => {
		const admin = requireAdmin(event.locals);
		const adminContainer = createAdminContainer(event);
		const formData = await event.request.formData();
		const userId = formString(formData, 'userId');
		const status = parseUserStatus(formData.get('status'));

		if (!userId || !status) {
			return fail(400, { message: 'Invalid status change request' });
		}

		try {
			await adminContainer.setUserStatus.execute(admin, userId, status);
		} catch (caught) {
			return adminActionFailure(caught);
		}

		return { ok: true };
	},
	setRole: async (event) => {
		const admin = requireAdmin(event.locals);
		const adminContainer = createAdminContainer(event);
		const formData = await event.request.formData();
		const userId = formString(formData, 'userId');
		const role = parseRole(formData.get('role'));

		if (!userId || !role) {
			return fail(400, { message: 'Invalid role change request' });
		}

		try {
			await adminContainer.setUserRole.execute(admin, userId, role);
		} catch (caught) {
			return adminActionFailure(caught);
		}

		return { ok: true };
	}
};
