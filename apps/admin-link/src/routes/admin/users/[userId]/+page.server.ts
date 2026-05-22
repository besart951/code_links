import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import { adminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

export async function load({ locals, params }) {
	const admin = requireAdmin(locals);
	const user = await adminContainer.getUserDetail.execute(admin, params.userId);
	const attempts = await adminContainer.listLoginAttempts.execute(admin, {
		userId: params.userId,
		page: 1,
		pageSize: 25
	});

	return {
		user,
		attempts: attempts.items,
		maskIp: !canSeeRawIpAddress(admin)
	};
}
