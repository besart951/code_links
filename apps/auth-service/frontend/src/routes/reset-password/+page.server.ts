import { fail } from '@sveltejs/kit';
import { postAuthJson } from '$lib/server/auth-api';

export const actions = {
	default: async (event) => {
		const formData = await event.request.formData();
		const token = String(formData.get('token') ?? event.url.searchParams.get('token') ?? '');
		const password = String(formData.get('password') ?? '');
		const confirmPassword = String(formData.get('confirmPassword') ?? '');

		if (password !== confirmPassword) {
			return fail(400, { token, error: 'Die Passwörter stimmen nicht überein.' });
		}

		const result = await postAuthJson(event, '/api/auth/reset-password', { token, password });
		if (!result.ok) {
			return fail(result.status, {
				token,
				error: result.body.error ?? 'Passwort konnte nicht geändert werden.'
			});
		}

		return { success: true, message: 'Passwort geändert. Du kannst dich jetzt anmelden.' };
	}
};
