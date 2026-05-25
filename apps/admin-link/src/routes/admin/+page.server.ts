import { createAdminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';
import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import { adminLoad } from '$lib/server/admin-route-helpers';

export async function load(event) {
	return adminLoad(async () => {
		const admin = requireAdmin(event.locals);
		const adminContainer = createAdminContainer(event);

		return {
			summary: await adminContainer.getDashboardSummary.execute(admin),
			maskIp: !canSeeRawIpAddress(admin)
		};
	});
}
