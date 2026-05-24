import { describe, expect, it } from 'vitest';
import { parseLoginAttemptQuery, parseRole, parseUserListQuery, parseUserStatus } from './admin-route-helpers';

describe('admin route helpers', () => {
	it('parses user list filters with safe defaults', () => {
		const query = parseUserListQuery(
			new URL('http://admin-link.localhost/admin/users?query=demo&role=support&status=locked&page=2&pageSize=10')
		);

		expect(query).toMatchObject({
			query: 'demo',
			role: 'support',
			status: 'locked',
			page: 2,
			pageSize: 10,
			sort: { field: 'createdAt', direction: 'desc' }
		});
	});

	it('rejects invalid form enum values', () => {
		expect(parseUserStatus('deleted')).toBeUndefined();
		expect(parseRole('owner')).toBeUndefined();
	});

	it('parses login attempt filters', () => {
		expect(
			parseLoginAttemptQuery(new URL('http://admin-link.localhost/admin/login-history?query=demo&success=false'))
		).toMatchObject({ query: 'demo', success: false, page: 1, pageSize: 50 });
	});
});
