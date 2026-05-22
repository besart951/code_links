import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import { adminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

export async function load({ locals }) {
	const admin = requireAdmin(locals);
	const attempts = await adminContainer.listLoginAttempts.execute(admin, {
		success: false,
		page: 1,
		pageSize: 50
	});

	return {
		...attempts,
		maskIp: !canSeeRawIpAddress(admin)
	};
}
