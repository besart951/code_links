import { getAuthJson } from '$lib/server/auth-api';

export const load = async (event) => {
	const token = event.url.searchParams.get('token');
	if (!token) {
		return { verified: false, error: 'Der Bestätigungslink ist unvollständig.' };
	}

	const result = await getAuthJson<{ status?: string; email?: string }>(
		event,
		`/api/auth/verify-email?token=${encodeURIComponent(token)}`
	);

	return {
		verified: result.ok,
		email: result.body.email,
		error: result.ok ? null : (result.body.error ?? 'Der Bestätigungslink ist ungültig oder abgelaufen.')
	};
};
