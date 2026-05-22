import { requirePermission } from '$lib/application/admin-access/requirePermission';
import { requireAdmin } from '$lib/server/auth';

export function load({ locals }) {
	const admin = requireAdmin(locals);
	requirePermission(admin, 'admin.audit_entries.read');

	return {
		entries: [
			{
				id: 'audit_001',
				actor: 'CodeLinks Admin',
				action: 'user.status.lock',
				target: 'lena.schneider@example.com',
				createdAt: '2026-05-22T09:01:00.000Z'
			}
		]
	};
}
