import { fail } from '@sveltejs/kit';
import { postAuthJson } from '$lib/server/auth-api';

export const actions = {
	default: async (event) => {
		const formData = await event.request.formData();
		const email = String(formData.get('email') ?? '');

		const result = await postAuthJson<{
			message?: string;
			debugResetToken?: string;
			debugResetUrl?: string;
		}>(event, '/api/auth/forgot-password', { email });

		if (!result.ok) {
			return fail(result.status, {
				email,
				error: result.body.error ?? 'Reset-Link konnte nicht angefordert werden.'
			});
		}

		return {
			success: true,
			email,
			message:
				result.body.message ??
				'Falls ein Konto mit dieser E-Mail existiert, wurde ein Reset-Link versendet.',
			debugResetUrl: result.body.debugResetUrl,
			debugResetToken: result.body.debugResetToken
		};
	}
};
