import { describe, expect, it } from 'vitest';
import { assertAdminMockRepositoryAllowed, shouldUseMockRepository } from './admin-container';

describe('admin container guards', () => {
	it('blocks mock repository in production', () => {
		try {
			assertAdminMockRepositoryAllowed(true, false);
			expect.unreachable('expected production mock repository guard to throw');
		} catch (caught) {
			expect(caught).toMatchObject({
				status: 500,
				body: { message: 'Admin mock repository is development-only' }
			});
		}
	});

	it('only allows ADMIN_LINK_MOCK_REPOSITORY in development', () => {
		expect(shouldUseMockRepository({ ADMIN_LINK_MOCK_REPOSITORY: 'true' }, true)).toBe(true);
		expect(shouldUseMockRepository({ ADMIN_LINK_MOCK_REPOSITORY: 'true' }, false)).toBe(false);
	});
});
