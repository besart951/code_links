import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import { adminLoad, parseLoginAttemptQuery } from '$lib/server/admin-route-helpers';
import { createAdminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

export async function load(event) {
	return adminLoad(async () => {
		const admin = requireAdmin(event.locals);
		const adminContainer = createAdminContainer(event);
		const attempts = await adminContainer.listLoginAttempts.execute(admin, parseLoginAttemptQuery(event.url));

		return {
			...attempts,
			maskIp: !canSeeRawIpAddress(admin)
		};
	});
}
