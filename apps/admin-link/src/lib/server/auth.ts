import { error } from '@sveltejs/kit';

export function requireAdmin(locals: App.Locals) {
	if (!locals.admin) {
		error(403, 'Admin access required');
	}

	return locals.admin;
}
