import { dev } from '$app/environment';
import { env } from '$env/dynamic/private';
import { error, redirect, type Handle } from '@sveltejs/kit';
import { createMockAdminActor } from '$lib/server/mock-admin-auth';
import type { AdminActor } from '$lib/domain/admin-access/types';

export const handle: Handle = async ({ event, resolve }) => {
	const allowMockAdmin = env.ADMIN_LINK_MOCK_AUTH === 'true';
	if (!dev && allowMockAdmin) {
		error(500, 'ADMIN_LINK_MOCK_AUTH is development-only');
	}
	const requestedRole = allowMockAdmin ? event.url.searchParams.get('as') ?? event.cookies.get('admin_role') ?? null : null;

	event.locals.adminMode = allowMockAdmin ? 'mock' : 'real';
	event.locals.admin = allowMockAdmin ? createMockAdminActor(requestedRole) : await fetchAdminActor(event);

	if (event.url.pathname.startsWith('/admin') && !event.locals.admin) {
		const loginUrl = new URL(env.AUTH_FRONTEND_LOGIN_URL ?? 'http://auth.codelinks.localhost/login');
		loginUrl.searchParams.set('redirectTo', event.url.toString());
		throw redirect(303, loginUrl.toString());
	}

	return resolve(event);
};

async function fetchAdminActor(event: Parameters<Handle>[0]['event']): Promise<AdminActor | null> {
	const baseUrl = env.AUTH_API_BASE_URL ?? 'http://localhost:8080';
	try {
		const response = await event.fetch(`${baseUrl}/api/admin/me`, {
			headers: {
				accept: 'application/json',
				cookie: event.request.headers.get('cookie') ?? ''
			}
		});
		if (!response.ok) {
			return null;
		}

		return (await response.json()) as AdminActor;
	} catch {
		return null;
	}
}
