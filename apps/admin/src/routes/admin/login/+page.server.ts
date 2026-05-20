import { loginToPlatform, PlatformRequestError } from '$lib/server/platformClient.js';
import { fail, redirect, type Actions } from '@sveltejs/kit';

export const actions: Actions = {
	default: async (event) => {
		const form = await event.request.formData();
		const email = String(form.get('email') ?? '').trim();
		const password = String(form.get('password') ?? '');
		const returnTo = String(form.get('returnTo') ?? '/admin');

		if (!email || !password) {
			return fail(400, { error: 'missing_credentials', email });
		}

		try {
			await loginToPlatform(event, { email, password });
		} catch (err) {
			if (err instanceof PlatformRequestError) {
				return fail(err.status === 503 ? 503 : 401, { error: err.code, email });
			}
			throw err;
		}

		throw redirect(303, safeReturnTo(returnTo));
	}
};

function safeReturnTo(returnTo: string): string {
	if (!returnTo.startsWith('/admin') || returnTo.startsWith('//')) {
		return '/admin';
	}
	return returnTo;
}
