import { createAdminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';
import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';

export async function load(event) {
	const admin = requireAdmin(event.locals);
	const adminContainer = createAdminContainer(event);

	return {
		summary: await adminContainer.getDashboardSummary.execute(admin),
		maskIp: !canSeeRawIpAddress(admin)
	};
}
