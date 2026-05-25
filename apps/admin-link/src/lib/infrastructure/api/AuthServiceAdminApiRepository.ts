import type { AdminRepository } from '$lib/application/ports/AdminRepository';
import type { AdminAuditEntry } from '$lib/domain/audit/types';
import type { LoginAttempt, LoginAttemptQuery } from '$lib/domain/auth-logs/types';
import type { SecurityEvent } from '$lib/domain/security/types';
import type { Paginated } from '$lib/domain/shared/pagination';
import type { DashboardSummary } from '$lib/domain/statistics/types';
import type { Notification } from '$lib/domain/notifications/types';
import type { RuntimeLogEntry } from '$lib/domain/runtime-logs/types';
import type { SmtpSettings, UpdateSmtpSettingsInput } from '$lib/domain/smtp/types';
import type { ManagedUserDetail, UserListItem, UserListQuery, UserStatus } from '$lib/domain/users/types';

type HeaderSource = () => Record<string, string> | Promise<Record<string, string>>;

type AuthServiceErrorBody = {
	error?: string;
	code?: string;
	message?: string;
};

export class AdminApiError extends Error {
	constructor(
		readonly status: number,
		readonly code: string,
		message: string
	) {
		super(message);
		this.name = 'AdminApiError';
	}
}

export class AuthServiceAdminApiRepository implements AdminRepository {
	constructor(
		private readonly baseUrl: string,
		private readonly fetchImpl: typeof fetch,
		private readonly requestHeaders: Record<string, string>,
		private readonly commandHeaderSource?: HeaderSource
	) {}

	private async getJson<T>(path: string): Promise<T> {
		const response = await this.fetchImpl(`${this.baseUrl}${path}`, { headers: this.requestHeaders });

		if (!response.ok) {
			throw await adminApiError(response);
		}

		return response.json() as Promise<T>;
	}

	private async commandHeaders() {
		return this.commandHeaderSource ? await this.commandHeaderSource() : this.requestHeaders;
	}

	async getDashboardSummary(): Promise<DashboardSummary> {
		const stats = await this.getJson<{
			users: DashboardSummary['users'];
			loginAttempts: DashboardSummary['loginAttempts'];
			passwordResetRequests: number;
			notifications: number;
			openSecurityEvents: number;
			topCountries: { key: string; count: number }[];
			topIpAddresses: { key: string; count: number }[];
		}>('/api/admin/dashboard/stats');

		return {
			users: stats.users,
			loginAttempts: stats.loginAttempts,
			passwordResetRequests: stats.passwordResetRequests,
			notifications: stats.notifications,
			security: {
				openEvents: stats.openSecurityEvents,
				suspiciousAttempts: stats.openSecurityEvents
			},
			topCountries: stats.topCountries.map((entry) => ({ countryCode: entry.key, count: entry.count })),
			topIpAddresses: stats.topIpAddresses.map((entry) => ({ ipAddress: entry.key, count: entry.count })),
			trend: [],
			recentActivity: [],
			highlightedEvents: []
		};
	}

	listUsers(query: UserListQuery) {
		const params = new URLSearchParams({
			page: String(query.page),
			pageSize: String(query.pageSize),
			sort: query.sort.field,
			direction: query.sort.direction
		});

		if (query.query) params.set('query', query.query);
		if (query.role) params.set('role', query.role);
		if (query.status) params.set('status', query.status);

		return this.getJson<Paginated<UserListItem>>(`/api/admin/users?${params}`);
	}

	getUserDetail(userId: string) {
		return this.getJson<ManagedUserDetail>(`/api/admin/users/${userId}`);
	}

	listLoginAttempts(query: LoginAttemptQuery) {
		const params = new URLSearchParams({
			page: String(query.page),
			pageSize: String(query.pageSize)
		});

		if (query.userId) params.set('userId', query.userId);
		if (query.query) params.set('query', query.query);
		if (query.success !== undefined) params.set('success', String(query.success));

		return this.getJson<Paginated<LoginAttempt>>(`/api/admin/login-attempts?${params}`);
	}

	listSecurityEvents() {
		return this.getJson<SecurityEvent[]>('/api/admin/security-events');
	}

	listNotifications() {
		return this.getJson<Notification[]>('/api/admin/notifications');
	}

	listAuditEntries() {
		return this.getJson<AdminAuditEntry[]>('/api/admin/audit-entries');
	}

	listRuntimeLogs() {
		return this.getJson<RuntimeLogEntry[]>('/api/admin/runtime-logs');
	}

	getSmtpSettings() {
		return this.getJson<SmtpSettings>('/api/admin/settings/smtp');
	}

	async updateSmtpSettings(input: UpdateSmtpSettingsInput) {
		const response = await this.fetchImpl(`${this.baseUrl}/api/admin/settings/smtp`, {
			method: 'PUT',
			headers: {
				...(await this.commandHeaders()),
				'content-type': 'application/json'
			},
			body: JSON.stringify(input)
		});

		if (!response.ok) {
			throw await adminApiError(response);
		}

		return response.json() as Promise<SmtpSettings>;
	}

	async sendTestEmail(recipient: string) {
		const response = await this.fetchImpl(`${this.baseUrl}/api/admin/settings/smtp/test-email`, {
			method: 'POST',
			headers: {
				...(await this.commandHeaders()),
				'content-type': 'application/json'
			},
			body: JSON.stringify({ recipient })
		});

		if (!response.ok) {
			throw await adminApiError(response);
		}
	}

	async setUserStatus(userId: string, status: UserStatus) {
		const response = await this.fetchImpl(`${this.baseUrl}/api/admin/users/${userId}/status`, {
			method: 'PATCH',
			headers: {
				...(await this.commandHeaders()),
				'content-type': 'application/json'
			},
			body: JSON.stringify({ status })
		});

		if (!response.ok) {
			throw await adminApiError(response);
		}
	}

	async setUserRole(userId: string, role: string) {
		const response = await this.fetchImpl(`${this.baseUrl}/api/admin/users/${userId}/role`, {
			method: 'PATCH',
			headers: {
				...(await this.commandHeaders()),
				'content-type': 'application/json'
			},
			body: JSON.stringify({ role })
		});

		if (!response.ok) {
			throw await adminApiError(response);
		}
	}
}

async function adminApiError(response: Response) {
	const body = (await response.json().catch(() => ({}))) as AuthServiceErrorBody;
	const code = body.code ?? `http_${response.status}`;
	const message = body.error ?? body.message ?? `Admin API request failed: ${response.status}`;
	return new AdminApiError(response.status, code, message);
}
