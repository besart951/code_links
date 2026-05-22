import { describe, expect, it } from 'vitest';
import { AuthManager } from './AuthManager.svelte.js';

describe('AuthManager', () => {
	it('stores the logged-in user and checks licenses', async () => {
		expect.assertions(3);

		const manager = new AuthManager({
			authBaseUrl: 'http://auth.test',
			fetcher: async () =>
				new Response(
					JSON.stringify({
						accessToken: 'header.payload.signature',
						expiresAt: '2026-05-22T12:00:00.000Z',
						user: {
							id: 'user-1',
							email: 'demo@codelinks.dev',
							name: 'Demo User',
							licenses: ['infra-link']
						}
					}),
					{ status: 200, headers: { 'content-type': 'application/json' } }
				)
		});

		await manager.login({ email: 'demo@codelinks.dev', password: 'password' });

		expect(manager.isAuthenticated).toBe(true);
		expect(manager.hasLicense('infra-link')).toBe(true);
		expect(manager.hasLicense('planer-link')).toBe(false);
	});
});
