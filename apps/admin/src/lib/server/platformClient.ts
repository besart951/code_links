import { env } from '$env/dynamic/private';
import { error, type RequestEvent } from '@sveltejs/kit';
import type {
	AdminDashboardSummary,
	AdminMeResponse,
	AdminProductSummary,
	AdminSearchResponse,
	AdminSetting,
	AdminTenantSummary,
	AdminUserSummary,
	AuditLogEntry,
	NotificationDeliverySummary,
	NotificationTemplateSummary,
	PageResponse,
	SecurityEventSummary,
	SubscriptionStatus
} from '@codelinks/contracts';

export class PlatformRequestError extends Error {
	constructor(
		readonly status: number,
		readonly code: string
	) {
		super(code);
	}
}

type PlatformFetchOptions = {
	search?: URLSearchParams;
};

const baseUrls = () =>
	env.PLATFORM_API_URL ? [env.PLATFORM_API_URL] : ['http://localhost:8080', 'http://localhost:8081'];

export async function platformJson<T>(
	event: RequestEvent,
	path: string,
	options: PlatformFetchOptions = {}
): Promise<T> {
	let response: Response;
	try {
		response = await fetchPlatform(event, path, options, {
			headers: {
				accept: 'application/json',
				cookie: event.request.headers.get('cookie') ?? ''
			}
		});
	} catch {
		throw new PlatformRequestError(503, 'platform_unavailable');
	}

	if (!response.ok) {
		let code = 'platform_request_failed';
		try {
			const body = (await response.json()) as { error?: string };
			code = body.error ?? code;
		} catch {
			// Keep the stable fallback code.
		}
		throw new PlatformRequestError(response.status, code);
	}

	return (await response.json()) as T;
}

export async function requireAdmin(event: RequestEvent): Promise<AdminMeResponse> {
	return platformJson<AdminMeResponse>(event, '/api/v1/admin/me');
}

export async function loginToPlatform(
	event: RequestEvent,
	input: { email: string; password: string }
): Promise<void> {
	let response: Response;
	try {
		response = await fetchPlatform(event, '/api/v1/auth/login', {}, {
			method: 'POST',
			headers: {
				accept: 'application/json',
				'content-type': 'application/json'
			},
			body: JSON.stringify(input)
		});
	} catch {
		throw new PlatformRequestError(503, 'platform_unavailable');
	}

	if (!response.ok) {
		let code = 'login_failed';
		try {
			const body = (await response.json()) as { error?: string };
			code = body.error ?? code;
		} catch {
			// Keep the stable fallback code.
		}
		throw new PlatformRequestError(response.status, code);
	}

	for (const rawCookie of getSetCookieHeaders(response.headers)) {
		const cookie = parseSetCookie(rawCookie);
		if (!cookie || !['platform_session', 'refresh_token', 'csrf_token'].includes(cookie.name)) {
			continue;
		}
		event.cookies.set(cookie.name, cookie.value, {
			path: cookie.path ?? '/',
			httpOnly: cookie.httpOnly,
			secure: cookie.secure,
			sameSite: cookie.sameSite,
			expires: cookie.expires
		});
	}
}

export async function getDashboard(event: RequestEvent): Promise<AdminDashboardSummary> {
	return platformJson<AdminDashboardSummary>(event, '/api/v1/admin/dashboard');
}

export async function searchAdmin(event: RequestEvent, query: string): Promise<AdminSearchResponse> {
	const params = new URLSearchParams();
	if (query) params.set('q', query);
	return platformJson<AdminSearchResponse>(event, '/api/v1/admin/search', { search: params });
}

export async function listTenants(event: RequestEvent): Promise<PageResponse<AdminTenantSummary>> {
	return platformJson<PageResponse<AdminTenantSummary>>(event, '/api/v1/admin/tenants');
}

export async function getTenant(event: RequestEvent, tenantId: string): Promise<AdminTenantSummary> {
	return platformJson<AdminTenantSummary>(event, `/api/v1/admin/tenants/${encodeURIComponent(tenantId)}`);
}

export async function listUsers(event: RequestEvent): Promise<PageResponse<AdminUserSummary>> {
	return platformJson<PageResponse<AdminUserSummary>>(event, '/api/v1/admin/users');
}

export async function getUser(event: RequestEvent, userId: string): Promise<AdminUserSummary> {
	return platformJson<AdminUserSummary>(event, `/api/v1/admin/users/${encodeURIComponent(userId)}`);
}

export async function listAudit(event: RequestEvent): Promise<PageResponse<AuditLogEntry>> {
	return platformJson<PageResponse<AuditLogEntry>>(event, '/api/v1/admin/audit');
}

export async function listNotifications(event: RequestEvent): Promise<{
	templates: NotificationTemplateSummary[];
	deliveries: NotificationDeliverySummary[];
}> {
	return platformJson(event, '/api/v1/admin/notifications');
}

export async function listSecurity(event: RequestEvent): Promise<PageResponse<SecurityEventSummary>> {
	return platformJson<PageResponse<SecurityEventSummary>>(event, '/api/v1/admin/security');
}

export async function listSettings(event: RequestEvent): Promise<PageResponse<AdminSetting>> {
	return platformJson<PageResponse<AdminSetting>>(event, '/api/v1/admin/settings');
}

export async function listProducts(event: RequestEvent): Promise<{ products: AdminProductSummary[] }> {
	return platformJson(event, '/api/v1/admin/products');
}

export async function listSubscriptions(event: RequestEvent): Promise<
	PageResponse<{
		id: string;
		tenant_id: string;
		tenant_name: string;
		product_key: string;
		plan_name: string;
		status: SubscriptionStatus;
		current_period_end: string | null;
	}>
> {
	return platformJson(event, '/api/v1/admin/subscriptions');
}

export function mapPlatformError(err: unknown): never {
	if (err instanceof PlatformRequestError) {
		throw error(err.status, err.code);
	}
	throw err;
}

function getSetCookieHeaders(headers: Headers): string[] {
	const withGetter = headers as Headers & { getSetCookie?: () => string[] };
	const cookies = withGetter.getSetCookie?.();
	if (cookies && cookies.length > 0) {
		return cookies;
	}
	const combined = headers.get('set-cookie');
	if (!combined) {
		return [];
	}
	return combined.split(/,(?=\s*[^;,]+=)/g).map((value) => value.trim());
}

async function fetchPlatform(
	event: RequestEvent,
	path: string,
	options: PlatformFetchOptions,
	init: RequestInit
): Promise<Response> {
	let lastError: unknown;
	for (const baseUrl of baseUrls()) {
		const url = new URL(path, baseUrl);
		if (options.search) {
			options.search.forEach((value, key) => url.searchParams.set(key, value));
		}
		try {
			return await event.fetch(url, init);
		} catch (err) {
			lastError = err;
		}
	}
	throw lastError;
}

function parseSetCookie(rawCookie: string):
	| {
			name: string;
			value: string;
			path?: string;
			httpOnly: boolean;
			secure: boolean;
			sameSite: 'lax' | 'strict' | 'none';
			expires?: Date;
	  }
	| undefined {
	const [nameValue, ...attributes] = rawCookie.split(';');
	const [name, ...valueParts] = nameValue.split('=');
	const value = valueParts.join('=');
	if (!name || !value) {
		return undefined;
	}

	let path: string | undefined;
	let expires: Date | undefined;
	let httpOnly = false;
	let secure = false;
	let sameSite: 'lax' | 'strict' | 'none' = 'lax';

	for (const attribute of attributes) {
		const [rawKey, ...rawValueParts] = attribute.trim().split('=');
		const key = rawKey.toLowerCase();
		const rawValue = rawValueParts.join('=');
		if (key === 'path') {
			path = rawValue || '/';
		} else if (key === 'expires') {
			const parsed = new Date(rawValue);
			if (!Number.isNaN(parsed.getTime())) {
				expires = parsed;
			}
		} else if (key === 'httponly') {
			httpOnly = true;
		} else if (key === 'secure') {
			secure = true;
		} else if (key === 'samesite') {
			const normalized = rawValue.toLowerCase();
			if (normalized === 'strict' || normalized === 'none') {
				sameSite = normalized;
			}
		}
	}

	return { name, value, path, httpOnly, secure, sameSite, expires };
}
