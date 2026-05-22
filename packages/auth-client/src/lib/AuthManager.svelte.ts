export type ProductLicense = 'infra-link' | 'planer-link' | 'loka-link' | (string & {});

export interface UserSession {
	id: string;
	email: string;
	name: string;
	licenses: ProductLicense[];
}

export interface LoginCredentials {
	email: string;
	password: string;
}

export interface LoginResponse {
	accessToken: string;
	user: UserSession;
	expiresAt: string;
}

export interface AuthManagerOptions {
	authBaseUrl: string;
	fetcher?: typeof fetch;
}

interface JwtPayload {
	sub?: string;
	email?: string;
	name?: string;
	licenses?: string[];
	exp?: number;
}

export class AuthManager {
	private authBaseUrl: string;
	private fetcher: typeof fetch;

	accessToken = $state<string | null>(null);
	user = $state<UserSession | null>(null);
	expiresAt = $state<Date | null>(null);
	isLoading = $state(false);
	error = $state<string | null>(null);
	isAuthenticated = $derived(Boolean(this.accessToken && this.user));

	constructor(options: AuthManagerOptions) {
		this.authBaseUrl = options.authBaseUrl.replace(/\/$/, '');
		this.fetcher = options.fetcher ?? fetch;
	}

	login = async (credentials: LoginCredentials): Promise<UserSession> => {
		this.isLoading = true;
		this.error = null;

		try {
			const response = await this.fetcher(`${this.authBaseUrl}/api/auth/login`, {
				method: 'POST',
				headers: { 'content-type': 'application/json' },
				credentials: 'include',
				body: JSON.stringify(credentials)
			});

			const data = await parseJson<LoginResponse>(response);
			this.applySession(data);
			return data.user;
		} catch (error) {
			this.clearSession();
			this.error = error instanceof Error ? error.message : 'Login failed';
			throw error;
		} finally {
			this.isLoading = false;
		}
	};

	refresh = async (): Promise<UserSession | null> => {
		this.error = null;

		const response = await this.fetcher(`${this.authBaseUrl}/api/auth/refresh`, {
			method: 'POST',
			credentials: 'include'
		});

		if (response.status === 401) {
			this.clearSession();
			return null;
		}

		const data = await parseJson<LoginResponse>(response);
		this.applySession(data);
		return data.user;
	};

	logout = async (): Promise<void> => {
		try {
			await this.fetcher(`${this.authBaseUrl}/api/auth/logout`, {
				method: 'POST',
				credentials: 'include',
				headers: this.authorizationHeaders()
			});
		} finally {
			this.clearSession();
		}
	};

	mockPurchase = async (productId: ProductLicense): Promise<UserSession> => {
		const response = await this.fetcher(`${this.authBaseUrl}/api/licenses/mock-purchase`, {
			method: 'POST',
			headers: {
				...this.authorizationHeaders(),
				'content-type': 'application/json'
			},
			credentials: 'include',
			body: JSON.stringify({ productId })
		});

		const data = await parseJson<LoginResponse>(response);
		this.applySession(data);
		return data.user;
	};

	hasLicense(productId: string): boolean {
		return this.user?.licenses.includes(productId) ?? false;
	}

	authorizationHeaders(): Record<string, string> {
		return this.accessToken ? { authorization: `Bearer ${this.accessToken}` } : {};
	}

	private applySession(response: LoginResponse): void {
		this.accessToken = response.accessToken;
		this.expiresAt = new Date(response.expiresAt);
		this.user = response.user ?? userFromJwt(response.accessToken);
	}

	private clearSession(): void {
		this.accessToken = null;
		this.user = null;
		this.expiresAt = null;
	}
}

async function parseJson<T>(response: Response): Promise<T> {
	const body = (await response.json().catch(() => null)) as unknown;

	if (!response.ok) {
		const message = isErrorBody(body) ? body.error : `Request failed with ${response.status}`;
		throw new Error(message);
	}

	return body as T;
}

function isErrorBody(value: unknown): value is { error: string } {
	return typeof value === 'object' && value !== null && 'error' in value && typeof value.error === 'string';
}

function userFromJwt(token: string): UserSession {
	const payload = decodePayload(token);

	return {
		id: requireString(payload.sub, 'sub'),
		email: requireString(payload.email, 'email'),
		name: requireString(payload.name, 'name'),
		licenses: payload.licenses ?? []
	};
}

function decodePayload(token: string): JwtPayload {
	const [, payload] = token.split('.');
	if (!payload) {
		throw new Error('Access token is not a valid JWT');
	}

	const json = base64UrlDecode(payload);
	return JSON.parse(json) as JwtPayload;
}

function base64UrlDecode(input: string): string {
	const padded = input.replace(/-/g, '+').replace(/_/g, '/').padEnd(Math.ceil(input.length / 4) * 4, '=');

	if (typeof atob === 'function') {
		return atob(padded);
	}

	throw new Error('No base64 decoder is available in this environment');
}

function requireString(value: string | undefined, claim: string): string {
	if (!value) {
		throw new Error(`Access token is missing ${claim}`);
	}

	return value;
}
