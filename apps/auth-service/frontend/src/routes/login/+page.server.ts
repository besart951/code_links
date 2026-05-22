import { fail, redirect } from '@sveltejs/kit';
import { postAuthJson } from '$lib/server/auth-api';
import { safeRedirectTo } from '$lib/server/redirects';

export const actions = {
	default: async (event) => {
		const formData = await event.request.formData();
		const email = String(formData.get('email') ?? '');
		const password = String(formData.get('password') ?? '');
		const redirectTo = safeRedirectTo(formData.get('redirectTo') ?? event.url.searchParams.get('redirectTo'), '/');

		const result = await postAuthJson(event, '/api/auth/login', { email, password });
		if (!result.ok) {
			return fail(result.status, {
				email,
				error: result.body.error ?? 'Login fehlgeschlagen.'
			});
		}

		throw redirect(303, redirectTo);
	}
};
