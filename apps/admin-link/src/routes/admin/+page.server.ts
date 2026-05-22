import { adminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';
import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';

export async function load({ locals }) {
	const admin = requireAdmin(locals);

	return {
		summary: await adminContainer.getDashboardSummary.execute(admin),
		maskIp: !canSeeRawIpAddress(admin)
	};
}
