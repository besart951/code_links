import { describe, expect, it } from 'vitest';
import { safeRedirectTo } from './redirects';

describe('safeRedirectTo', () => {
	it('allows relative and CodeLinks-local redirects only', () => {
		expect(safeRedirectTo('/products')).toBe('/products');
		expect(safeRedirectTo('http://admin-link.codelinks.localhost/admin')).toBe(
			'http://admin-link.codelinks.localhost/admin'
		);
		expect(safeRedirectTo('https://evil.example/phish', '/login')).toBe('/login');
		expect(safeRedirectTo('//evil.example/phish', '/login')).toBe('/login');
	});
});
