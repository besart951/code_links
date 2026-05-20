import { PlatformRequestError, requireAdmin } from '$lib/server/platformClient.js';
import { error, redirect, type Handle } from '@sveltejs/kit';

const publicAdminPaths = new Set(['/admin/login']);

export const handle: Handle = async ({ event, resolve }) => {
	if (!event.url.pathname.startsWith('/admin') || publicAdminPaths.has(event.url.pathname)) {
		return resolve(event);
	}

	try {
		const admin = await requireAdmin(event);
		if (!admin.superadmin) {
			throw error(403, 'superadmin_required');
		}
		event.locals.admin = admin;
	} catch (err) {
		if (err instanceof PlatformRequestError) {
			if (err.status === 401) {
				const returnTo = `${event.url.pathname}${event.url.search}`;
				throw redirect(303, `/admin/login?returnTo=${encodeURIComponent(returnTo)}`);
			}
			if (err.status === 403) {
				throw error(403, 'superadmin_required');
			}
			if (err.status === 503) {
				const returnTo = `${event.url.pathname}${event.url.search}`;
				throw redirect(
					303,
					`/admin/login?returnTo=${encodeURIComponent(returnTo)}&error=${encodeURIComponent(err.code)}`
				);
			}
		}
		throw err;
	}

	return resolve(event);
};
