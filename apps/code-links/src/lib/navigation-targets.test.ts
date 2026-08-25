import { describe, expect, it } from 'vitest';
import { resolveAuthAppUrl, resolveProductAppUrl } from './navigation-targets';

describe('navigation targets', () => {
	it('uses canonical codelinks URLs away from local dev hosts', () => {
		expect(resolveAuthAppUrl({}, 'code-links.codelinks.localhost')).toBe('http://auth.codelinks.localhost');
		expect(resolveProductAppUrl('infra-link', 'http://infra-link.codelinks.localhost', {}, 'code-links.codelinks.localhost')).toBe(
			'http://infra-link.codelinks.localhost'
		);
	});

	it('uses stable localhost ports in local dev mode', () => {
		expect(resolveAuthAppUrl({}, 'localhost')).toBe('http://localhost:5174');
		expect(resolveProductAppUrl('infra-link', 'http://infra-link.codelinks.localhost', {}, 'localhost')).toBe(
			'http://localhost:5175'
		);
		expect(resolveProductAppUrl('planer-link', 'http://planer-link.codelinks.localhost', {}, 'localhost')).toBe(
			'http://localhost:5176'
		);
		expect(resolveProductAppUrl('loka-link', 'http://loka-link.codelinks.localhost', {}, 'localhost')).toBe(
			'http://localhost:5177'
		);
	});

	it('lets public env overrides win over local defaults', () => {
		expect(resolveAuthAppUrl({ PUBLIC_AUTH_APP_URL: 'http://auth.example.test/' }, 'localhost')).toBe(
			'http://auth.example.test'
		);
		expect(
			resolveProductAppUrl(
				'infra-link',
				'http://infra-link.codelinks.localhost',
				{ PUBLIC_INFRA_LINK_APP_URL: 'http://infra.example.test/' },
				'localhost'
			)
		).toBe('http://infra.example.test');
	});
});
