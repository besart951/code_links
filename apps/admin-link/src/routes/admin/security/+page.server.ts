import { canSeeRawIpAddress } from '$lib/domain/admin-access/permissions';
import { adminContainer } from '$lib/server/admin-container';
import { requireAdmin } from '$lib/server/auth';

export async function load({ locals }) {
	const admin = requireAdmin(locals);

	return {
		events: await adminContainer.listSecurityEvents.execute(admin),
		maskIp: !canSeeRawIpAddress(admin)
	};
}
