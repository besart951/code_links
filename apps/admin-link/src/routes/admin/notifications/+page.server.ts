import { adminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

export async function load({ locals }) {
	const admin = requireAdmin(locals);

	return {
		notifications: await adminContainer.listNotifications.execute(admin)
	};
}
