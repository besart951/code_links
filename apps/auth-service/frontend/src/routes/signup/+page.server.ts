import { fail } from '@sveltejs/kit';
import { postAuthJson } from '$lib/server/auth-api';

export const actions = {
	default: async (event) => {
		const formData = await event.request.formData();
		const name = String(formData.get('name') ?? '');
		const email = String(formData.get('email') ?? '');
		const password = String(formData.get('password') ?? '');
		const confirmPassword = String(formData.get('confirmPassword') ?? '');
		const acceptedTerms = formData.get('acceptedTerms') === 'on';

		if (password !== confirmPassword) {
			return fail(400, { name, email, error: 'Die Passwörter stimmen nicht überein.' });
		}

		const result = await postAuthJson<{
			message?: string;
			debugVerificationToken?: string;
			verificationUrl?: string;
		}>(event, '/api/auth/signup', {
			name,
			email,
			password,
			acceptedTerms
		});

		if (!result.ok) {
			return fail(result.status, {
				name,
				email,
				error: result.body.error ?? 'Registrierung fehlgeschlagen.'
			});
		}

		return {
			success: true,
			message:
				result.body.message ?? 'Account erstellt. Bitte bestätige deine E-Mail-Adresse vor dem Login.',
			debugVerificationToken: result.body.debugVerificationToken,
			verificationUrl: result.body.verificationUrl
		};
	}
};
