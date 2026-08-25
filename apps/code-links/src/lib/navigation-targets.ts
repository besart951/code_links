import type { ProductId } from '@codelinks/config/products';

type PublicEnv = Record<string, string | undefined>;

const canonicalAuthAppUrl = 'http://auth.codelinks.localhost';

const productPorts: Record<ProductId, number> = {
	'infra-link': 5175,
	'planer-link': 5176,
	'loka-link': 5177
};

const productEnvKeys: Record<ProductId, string> = {
	'infra-link': 'PUBLIC_INFRA_LINK_APP_URL',
	'planer-link': 'PUBLIC_PLANER_LINK_APP_URL',
	'loka-link': 'PUBLIC_LOKA_LINK_APP_URL'
};

export function resolveAuthAppUrl(env: PublicEnv, hostname: string): string {
	return normalizeUrl(env.PUBLIC_AUTH_APP_URL ?? localDevUrl(hostname, 5174) ?? canonicalAuthAppUrl);
}

export function resolveProductAppUrl(
	productId: ProductId,
	canonicalUrl: string,
	env: PublicEnv,
	hostname: string
): string {
	return normalizeUrl(env[productEnvKeys[productId]] ?? localDevUrl(hostname, productPorts[productId]) ?? canonicalUrl);
}

export function isLocalDevHost(hostname: string): boolean {
	return hostname === 'localhost' || hostname === '127.0.0.1' || hostname === '::1' || hostname === '[::1]';
}

function localDevUrl(hostname: string, port: number): string | null {
	if (!isLocalDevHost(hostname)) return null;

	return `http://${formatHostForUrl(hostname)}:${port}`;
}

function formatHostForUrl(hostname: string): string {
	return hostname.includes(':') && !hostname.startsWith('[') ? `[${hostname}]` : hostname;
}

function normalizeUrl(url: string): string {
	return url.replace(/\/$/, '');
}
