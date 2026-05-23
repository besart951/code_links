import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import type { LoginAttempt, LoginAttemptQuery } from '$lib/domain/auth-logs/types';
import type { SecurityEvent } from '$lib/domain/security/types';
import type { Paginated } from '$lib/domain/shared/pagination';
import type { DashboardSummary } from '$lib/domain/statistics/types';
import type { Notification } from '$lib/domain/notifications/types';
import type { SmtpSettings, UpdateSmtpSettingsInput } from '$lib/domain/smtp/types';
import type { ManagedUserDetail, UserListItem, UserListQuery, UserStatus } from '$lib/domain/users/types';

export class AuthServiceAdminApiRepository implements AdminReadRepository {
	constructor(
		private readonly baseUrl: string,
		private readonly fetchImpl: typeof fetch,
		private readonly requestHeaders: Record<string, string>
	) {}

	private async getJson<T>(path: string): Promise<T> {
		const response = await this.fetchImpl(`${this.baseUrl}${path}`, { headers: this.requestHeaders });

		if (!response.ok) {
			throw new Error(`Admin API request failed: ${response.status}`);
		}

		return response.json() as Promise<T>;
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

	getSmtpSettings() {
		return this.getJson<SmtpSettings>('/api/admin/settings/smtp');
	}

	async updateSmtpSettings(input: UpdateSmtpSettingsInput) {
		const response = await this.fetchImpl(`${this.baseUrl}/api/admin/settings/smtp`, {
			method: 'PUT',
			headers: {
				...this.requestHeaders,
				'content-type': 'application/json'
			},
			body: JSON.stringify(input)
		});

		if (!response.ok) {
			throw new Error(`Admin API request failed: ${response.status}`);
		}

		return response.json() as Promise<SmtpSettings>;
	}

	async sendTestEmail(recipient: string) {
		const response = await this.fetchImpl(`${this.baseUrl}/api/admin/settings/smtp/test-email`, {
			method: 'POST',
			headers: {
				...this.requestHeaders,
				'content-type': 'application/json'
			},
			body: JSON.stringify({ recipient })
		});

		if (!response.ok) {
			throw new Error(`Admin API request failed: ${response.status}`);
		}
	}

	async setUserStatus(userId: string, status: UserStatus) {
		const response = await this.fetchImpl(`${this.baseUrl}/api/admin/users/${userId}/status`, {
			method: 'PATCH',
			headers: {
				...this.requestHeaders,
				'content-type': 'application/json'
			},
			body: JSON.stringify({ status })
		});

		if (!response.ok) {
			throw new Error(`Admin API request failed: ${response.status}`);
		}
	}

	async setUserRole(userId: string, role: string) {
		const response = await this.fetchImpl(`${this.baseUrl}/api/admin/users/${userId}/role`, {
			method: 'PATCH',
			headers: {
				...this.requestHeaders,
				'content-type': 'application/json'
			},
			body: JSON.stringify({ role })
		});

		if (!response.ok) {
			throw new Error(`Admin API request failed: ${response.status}`);
		}
	}
}
