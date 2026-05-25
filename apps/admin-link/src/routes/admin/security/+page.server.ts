import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import { createAdminContainer } from '$lib/server/admin-container';
import { adminLoad } from '$lib/server/admin-route-helpers';
import { requireAdmin } from '$lib/server/auth';

export async function load(event) {
	return adminLoad(async () => {
		const admin = requireAdmin(event.locals);
		const adminContainer = createAdminContainer(event);

		return {
			events: await adminContainer.listSecurityEvents.execute(admin),
			maskIp: !canSeeRawIpAddress(admin)
		};
	});
}
