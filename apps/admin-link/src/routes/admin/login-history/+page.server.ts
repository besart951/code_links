import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import { createAdminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

export async function load(event) {
	const admin = requireAdmin(event.locals);
	const adminContainer = createAdminContainer(event);
	const attempts = await adminContainer.listLoginAttempts.execute(admin, {
		query: event.url.searchParams.get('query') || undefined,
		page: Number(event.url.searchParams.get('page') ?? 1),
		pageSize: 50
	});

	return {
		...attempts,
		maskIp: !canSeeRawIpAddress(admin)
	};
}
