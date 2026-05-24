import { createAdminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

export async function load(event) {
	const admin = requireAdmin(event.locals);
	const adminContainer = createAdminContainer(event);

	return {
		entries: await adminContainer.listAuditEntries.execute(admin)
	};
}
