import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import { createAdminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

export async function load(event) {
	const admin = requireAdmin(event.locals);
	const adminContainer = createAdminContainer(event);
	const user = await adminContainer.getUserDetail.execute(admin, event.params.userId);
	const attempts = await adminContainer.listLoginAttempts.execute(admin, {
		userId: event.params.userId,
		page: 1,
		pageSize: 25
	});

	return {
		user,
		attempts: attempts.items,
		maskIp: !canSeeRawIpAddress(admin)
	};
}
