import { describe, expect, it, vi } from 'vitest';
import { AuthServiceAdminApiRepository } from './AuthServiceAdminApiRepository';

describe('AuthServiceAdminApiRepository', () => {
	it('forwards request cookies on query calls', async () => {
		const fetchImpl = vi.fn(async () => Response.json({ items: [], total: 0, page: 1, pageSize: 25 }));
		const repository = new AuthServiceAdminApiRepository('http://auth-service:8080', fetchImpl, {
			accept: 'application/json',
			cookie: 'refresh_token=abc'
		});

		await repository.listUsers({
			page: 1,
			pageSize: 25,
			sort: { field: 'createdAt', direction: 'desc' }
		});

		expect(fetchImpl).toHaveBeenCalledWith(
			'http://auth-service:8080/api/admin/users?page=1&pageSize=25&sort=createdAt&direction=desc',
			{ headers: { accept: 'application/json', cookie: 'refresh_token=abc' } }
		);
	});

	it('sends command bodies to the Auth Service Admin API', async () => {
		const fetchImpl = vi.fn(async () => new Response(null, { status: 204 }));
		const repository = new AuthServiceAdminApiRepository(
			'http://auth-service:8080',
			fetchImpl,
			{
				accept: 'application/json',
				cookie: 'refresh_token=abc'
			},
			async () => {
				return {
					accept: 'application/json',
					authorization: 'Bearer access-1'
				};
			}
		);

		await repository.setUserStatus('user-1', 'locked');

		expect(fetchImpl).toHaveBeenCalledWith('http://auth-service:8080/api/admin/users/user-1/status', {
			method: 'PATCH',
			headers: {
				accept: 'application/json',
				authorization: 'Bearer access-1',
				'content-type': 'application/json'
			},
			body: JSON.stringify({ status: 'locked' })
		});
	});
});
