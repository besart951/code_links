import { describe, expect, it } from 'vitest';
import { platformJson, PlatformRequestError } from './platformClient.js';
import type { RequestEvent } from '@sveltejs/kit';

function eventWithFetch(fetchImpl: RequestEvent['fetch']): RequestEvent {
	return {
		fetch: fetchImpl,
		request: new Request('http://localhost/admin')
	} as RequestEvent;
}

describe('platformJson', () => {
	it('maps Platform network failures to a typed unavailable error', async () => {
		const event = eventWithFetch(async () => {
			throw new TypeError('fetch failed');
		});

		await expect(platformJson(event, '/api/v1/admin/me')).rejects.toMatchObject({
			status: 503,
			code: 'platform_unavailable'
		});
		await expect(platformJson(event, '/api/v1/admin/me')).rejects.toBeInstanceOf(
			PlatformRequestError
		);
	});
});
