import { describe, expect, it } from 'vitest';
import { assertMockAdminAllowed } from './hooks.server';

describe('admin mock auth guard', () => {
	it('blocks mock auth in production', () => {
		try {
			assertMockAdminAllowed(true, false);
			expect.unreachable('expected production mock auth guard to throw');
		} catch (caught) {
			expect(caught).toMatchObject({
				status: 500,
				body: { message: 'ADMIN_LINK_MOCK_AUTH is development-only' }
			});
		}
	});

	it('allows mock auth in development', () => {
		expect(() => assertMockAdminAllowed(true, true)).not.toThrow();
	});
});
