import { requireAdmin } from '$lib/server/auth';

export function load({ locals }) {
	return {
		admin: requireAdmin(locals),
		adminMode: locals.adminMode
	};
}
