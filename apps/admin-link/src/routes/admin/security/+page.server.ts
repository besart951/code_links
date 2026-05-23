import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import { createAdminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

export async function load(event) {
	const admin = requireAdmin(event.locals);
	const adminContainer = createAdminContainer(event);

	return {
		events: await adminContainer.listSecurityEvents.execute(admin),
		maskIp: !canSeeRawIpAddress(admin)
	};
}
