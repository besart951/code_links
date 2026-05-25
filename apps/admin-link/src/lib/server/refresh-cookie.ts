import type { RequestEvent } from '@sveltejs/kit';

export function forwardRefreshCookie(event: RequestEvent, response: Response) {
	const headers = response.headers as Headers & { getSetCookie?: () => string[] };
	const cookies = headers.getSetCookie?.() ?? splitSetCookieHeader(headers.get('set-cookie'));

	for (const rawCookie of cookies) {
		const parsed = parseSetCookie(rawCookie);
		if (parsed.name !== 'refresh_token') continue;

		if (!parsed.value || parsed.maxAge === -1) {
			event.cookies.delete(parsed.name, {
				path: parsed.path ?? '/',
				domain: parsed.domain,
				secure: parsed.secure,
				sameSite: 'lax'
			});
			continue;
		}

		event.cookies.set(parsed.name, parsed.value, {
			path: parsed.path ?? '/',
			domain: parsed.domain,
			httpOnly: true,
			secure: parsed.secure,
			sameSite: 'lax',
			maxAge: parsed.maxAge
		});
	}
}

function splitSetCookieHeader(value: string | null): string[] {
	if (!value) return [];

	return value.split(/,(?=\s*refresh_token=)/).map((cookie) => cookie.trim());
}

function parseSetCookie(rawCookie: string) {
	const [nameValue, ...attributes] = rawCookie.split(';').map((part) => part.trim());
	const [name, value = ''] = nameValue.split('=');
	const parsed: {
		name: string;
		value: string;
		path?: string;
		domain?: string;
		secure: boolean;
		maxAge?: number;
	} = {
		name,
		value,
		secure: false
	};

	for (const attribute of attributes) {
		const [key, rawValue = ''] = attribute.split('=');
		switch (key.toLowerCase()) {
			case 'path':
				parsed.path = rawValue;
				break;
			case 'domain':
				parsed.domain = rawValue;
				break;
			case 'secure':
				parsed.secure = true;
				break;
			case 'max-age':
				parsed.maxAge = Number.parseInt(rawValue, 10);
				break;
		}
	}

	return parsed;
}
