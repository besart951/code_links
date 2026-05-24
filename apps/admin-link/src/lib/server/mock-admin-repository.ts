import { error } from '@sveltejs/kit';
import type { AdminRepository } from '$lib/application/ports/AdminRepository';
import type { AdminAuditEntry } from '$lib/domain/audit/types';
import type { LoginAttempt, LoginAttemptQuery } from '$lib/domain/auth-logs/types';
import type { Paginated } from '$lib/domain/shared/pagination';
import type { DashboardSummary } from '$lib/domain/statistics/types';
import type { ManagedUserDetail, UserListItem, UserListQuery, UserStatus } from '$lib/domain/users/types';
import type { SecurityEvent } from '$lib/domain/security/types';
import type { Notification } from '$lib/domain/notifications/types';
import type { RuntimeLogEntry } from '$lib/domain/runtime-logs/types';
import type { SmtpSettings, UpdateSmtpSettingsInput } from '$lib/domain/smtp/types';

const users: UserListItem[] = [
	{
		id: 'usr_anna',
		name: 'Anna Keller',
		email: 'anna.keller@example.com',
		primaryRole: 'admin',
		status: 'active',
		emailVerified: true,
		createdAt: '2026-05-02T09:15:00.000Z',
		lastLoginAt: '2026-05-22T09:21:00.000Z',
		successfulLoginCount: 42,
		failedLoginCount: 1,
		lastKnownIpAddress: '203.0.113.18',
		lastLoginCountryCode: 'CH'
	},
	{
		id: 'usr_mario',
		name: 'Mario Rossi',
		email: 'mario.rossi@example.com',
		primaryRole: 'support',
		status: 'active',
		emailVerified: true,
		createdAt: '2026-04-18T12:44:00.000Z',
		lastLoginAt: '2026-05-21T18:32:00.000Z',
		successfulLoginCount: 27,
		failedLoginCount: 4,
		lastKnownIpAddress: '198.51.100.44',
		lastLoginCountryCode: 'IT'
	},
	{
		id: 'usr_lena',
		name: 'Lena Schneider',
		email: 'lena.schneider@example.com',
		primaryRole: 'user',
		status: 'locked',
		emailVerified: true,
		createdAt: '2026-05-16T07:05:00.000Z',
		lastLoginAt: '2026-05-20T11:08:00.000Z',
		successfulLoginCount: 8,
		failedLoginCount: 18,
		lastKnownIpAddress: '192.0.2.77',
		lastLoginCountryCode: 'DE'
	},
	{
		id: 'usr_sofia',
		name: 'Sofia Martin',
		email: 'sofia.martin@example.com',
		primaryRole: 'auditor',
		status: 'disabled',
		emailVerified: false,
		createdAt: '2026-03-29T13:27:00.000Z',
		lastLoginAt: null,
		successfulLoginCount: 0,
		failedLoginCount: 2,
		lastKnownIpAddress: null,
		lastLoginCountryCode: null
	},
	{
		id: 'usr_noah',
		name: 'Noah Dubois',
		email: 'noah.dubois@example.com',
		primaryRole: 'user',
		status: 'active',
		emailVerified: true,
		createdAt: '2026-05-19T16:15:00.000Z',
		lastLoginAt: '2026-05-22T06:41:00.000Z',
		successfulLoginCount: 3,
		failedLoginCount: 0,
		lastKnownIpAddress: '203.0.113.81',
		lastLoginCountryCode: 'FR'
	}
];

const loginAttempts: LoginAttempt[] = [
	{
		id: 'login_001',
		userId: 'usr_anna',
		emailAttempted: 'anna.keller@example.com',
		occurredAt: '2026-05-22T09:21:00.000Z',
		ipAddress: '203.0.113.18',
		countryCode: 'CH',
		city: 'Zurich',
		success: true,
		failureReason: null,
		authMethod: 'password',
		device: { browser: 'Chrome', os: 'Windows', userAgent: 'Chrome 125 / Windows' },
		riskScore: 8
	},
	{
		id: 'login_002',
		userId: 'usr_lena',
		emailAttempted: 'lena.schneider@example.com',
		occurredAt: '2026-05-22T08:58:00.000Z',
		ipAddress: '192.0.2.77',
		countryCode: 'DE',
		city: 'Berlin',
		success: false,
		failureReason: 'account_locked',
		authMethod: 'password',
		device: { browser: 'Firefox', os: 'Linux', userAgent: 'Firefox 126 / Linux' },
		riskScore: 89
	},
	{
		id: 'login_003',
		userId: null,
		emailAttempted: 'unknown@example.com',
		occurredAt: '2026-05-22T08:53:00.000Z',
		ipAddress: '192.0.2.77',
		countryCode: 'DE',
		city: 'Berlin',
		success: false,
		failureReason: 'unknown_email',
		authMethod: 'password',
		device: { browser: 'Firefox', os: 'Linux', userAgent: 'Firefox 126 / Linux' },
		riskScore: 76
	},
	{
		id: 'login_004',
		userId: 'usr_mario',
		emailAttempted: 'mario.rossi@example.com',
		occurredAt: '2026-05-21T18:32:00.000Z',
		ipAddress: '198.51.100.44',
		countryCode: 'IT',
		city: 'Milan',
		success: true,
		failureReason: null,
		authMethod: 'refresh_token',
		device: { browser: 'Safari', os: 'macOS', userAgent: 'Safari 18 / macOS' },
		riskScore: 13
	},
	{
		id: 'login_005',
		userId: 'usr_noah',
		emailAttempted: 'noah.dubois@example.com',
		occurredAt: '2026-05-22T06:41:00.000Z',
		ipAddress: '203.0.113.81',
		countryCode: 'FR',
		city: 'Lyon',
		success: true,
		failureReason: null,
		authMethod: 'password',
		device: { browser: 'Edge', os: 'Windows', userAgent: 'Edge 125 / Windows' },
		riskScore: 18
	}
];

const securityEvents: SecurityEvent[] = [
	{
		id: 'sec_001',
		userId: 'usr_lena',
		type: 'many_failed_logins',
		severity: 'high',
		status: 'open',
		summary: '18 fehlgeschlagene Logins in den letzten 24 Stunden',
		detectedAt: '2026-05-22T09:00:00.000Z',
		resolvedAt: null,
		sourceIpAddress: '192.0.2.77',
		countryCode: 'DE'
	},
	{
		id: 'sec_002',
		userId: null,
		type: 'many_failures_from_ip',
		severity: 'medium',
		status: 'open',
		summary: 'Mehrere unbekannte E-Mail-Adressen von derselben IP',
		detectedAt: '2026-05-22T08:55:00.000Z',
		resolvedAt: null,
		sourceIpAddress: '192.0.2.77',
		countryCode: 'DE'
	}
];

const notifications: Notification[] = [
	{
		id: 'not_001',
		type: 'password_reset',
		channel: 'email',
		recipient: 'lena.schneider@example.com',
		subject: 'CodeLinks Passwort zurücksetzen',
		status: 'sent',
		createdAt: '2026-05-22T08:45:00.000Z',
		sentAt: '2026-05-22T08:45:07.000Z'
	},
	{
		id: 'not_002',
		type: 'smtp_test',
		channel: 'email',
		recipient: 'admin@example.com',
		subject: 'CodeLinks SMTP Test',
		status: 'sent',
		createdAt: '2026-05-21T14:12:00.000Z',
		sentAt: '2026-05-21T14:12:03.000Z'
	}
];

const auditEntries: AdminAuditEntry[] = [
	{
		id: 'audit_001',
		actorUserId: 'usr_anna',
		action: 'admin.users.lock',
		targetType: 'user',
		targetId: 'usr_lena',
		reason: 'suspicious_failed_logins',
		ipAddress: '203.0.113.18',
		createdAt: '2026-05-22T09:01:00.000Z'
	},
	{
		id: 'audit_002',
		actorUserId: 'usr_anna',
		action: 'admin.smtp_settings.test_email',
		targetType: 'smtp_settings',
		targetId: 'default',
		reason: 'admin@example.com',
		ipAddress: '203.0.113.18',
		createdAt: '2026-05-21T14:12:04.000Z'
	}
];

const runtimeLogs: RuntimeLogEntry[] = [
	{
		id: 'runtime_001',
		occurredAt: '2026-05-22T09:04:00.000Z',
		level: 'info',
		source: 'auth-service',
		message: 'auth-service listening on :8080',
		raw: '2026/05/22 09:04:00.000000 auth-service listening on :8080'
	},
	{
		id: 'runtime_002',
		occurredAt: '2026-05-22T09:05:00.000Z',
		level: 'fatal',
		source: 'auth-service',
		message: 'fatal: database unavailable',
		raw: '2026/05/22 09:05:00.000000 fatal: database unavailable'
	}
];

let smtpSettings: SmtpSettings = {
	host: 'smtp.example.com',
	port: 587,
	username: 'mailer',
	hasPassword: true,
	encryption: 'starttls',
	fromEmail: 'no-reply@codelinks.dev',
	fromName: 'CodeLinks',
	replyToEmail: 'support@codelinks.dev',
	active: true,
	updatedAt: '2026-05-22T07:30:00.000Z'
};

function paginate<T>(items: T[], page: number, pageSize: number): Paginated<T> {
	const start = (page - 1) * pageSize;

	return {
		items: items.slice(start, start + pageSize),
		total: items.length,
		page,
		pageSize
	};
}

function compareNullableString(left: string | null, right: string | null) {
	return (left ?? '').localeCompare(right ?? '');
}

export class MockAdminRepository implements AdminRepository {
	async getDashboardSummary(): Promise<DashboardSummary> {
		const successful = loginAttempts.filter((attempt) => attempt.success).length;
		const failed = loginAttempts.length - successful;

		return {
			users: {
				total: users.length,
				active: users.filter((user) => user.status === 'active').length,
				locked: users.filter((user) => user.status === 'locked').length,
				newLast7Days: 3,
				newLast30Days: 4
			},
			loginAttempts: {
				total: loginAttempts.length,
				successful,
				failed
			},
			passwordResetRequests: 6,
			notifications: notifications.length,
			security: {
				openEvents: securityEvents.filter((event) => event.status === 'open').length,
				suspiciousAttempts: loginAttempts.filter((attempt) => attempt.riskScore >= 70).length
			},
			topCountries: [
				{ countryCode: 'DE', count: 2 },
				{ countryCode: 'CH', count: 1 },
				{ countryCode: 'IT', count: 1 }
			],
			topIpAddresses: [
				{ ipAddress: '192.0.2.77', count: 2 },
				{ ipAddress: '203.0.113.18', count: 1 },
				{ ipAddress: '198.51.100.44', count: 1 }
			],
			trend: [
				{ date: '2026-05-16', successful: 9, failed: 2 },
				{ date: '2026-05-17', successful: 12, failed: 1 },
				{ date: '2026-05-18', successful: 11, failed: 3 },
				{ date: '2026-05-19', successful: 16, failed: 2 },
				{ date: '2026-05-20', successful: 15, failed: 5 },
				{ date: '2026-05-21', successful: 18, failed: 4 },
				{ date: '2026-05-22', successful: 7, failed: 9 }
			],
			recentActivity: loginAttempts.slice(0, 5),
			highlightedEvents: securityEvents.filter((event) => event.status === 'open')
		};
	}

	async listUsers(query: UserListQuery) {
		const search = query.query?.toLowerCase();
		let result = [...users];

		if (search) {
			result = result.filter((user) => user.name.toLowerCase().includes(search) || user.email.toLowerCase().includes(search));
		}

		if (query.role) {
			result = result.filter((user) => user.primaryRole === query.role);
		}

		if (query.status) {
			result = result.filter((user) => user.status === query.status);
		}

		result.sort((left, right) => {
			const direction = query.sort.direction === 'asc' ? 1 : -1;

			switch (query.sort.field) {
				case 'createdAt':
					return (Date.parse(left.createdAt) - Date.parse(right.createdAt)) * direction;
				case 'lastLoginAt':
					return compareNullableString(left.lastLoginAt, right.lastLoginAt) * direction;
				default:
					return String(left[query.sort.field]).localeCompare(String(right[query.sort.field])) * direction;
			}
		});

		return paginate(result, query.page, query.pageSize);
	}

	async getUserDetail(userId: string): Promise<ManagedUserDetail> {
		const user = users.find((candidate) => candidate.id === userId);

		if (!user) {
			error(404, 'User not found');
		}

		const attempts = loginAttempts.filter((attempt) => attempt.userId === userId);

		return {
			...user,
			roles: [
				{
					role: user.primaryRole,
					grantedAt: user.createdAt,
					grantedBy: 'system'
				}
			],
			productLicenses: user.id === 'usr_anna' ? ['infra-link', 'planer-link'] : user.id === 'usr_noah' ? ['infra-link'] : [],
			knownIpAddresses: [...new Set(attempts.map((attempt) => attempt.ipAddress))],
			loginCountries: [...new Set(attempts.map((attempt) => attempt.countryCode))],
			usedDevices: [...new Set(attempts.map((attempt) => `${attempt.device.browser} / ${attempt.device.os}`))]
		};
	}

	async listLoginAttempts(query: LoginAttemptQuery) {
		let result = [...loginAttempts];
		const search = query.query?.toLowerCase();

		if (query.userId) {
			result = result.filter((attempt) => attempt.userId === query.userId);
		}

		if (query.success !== undefined) {
			result = result.filter((attempt) => attempt.success === query.success);
		}

		if (search) {
			result = result.filter((attempt) => attempt.emailAttempted.toLowerCase().includes(search) || attempt.ipAddress.includes(search));
		}

		return paginate(result, query.page, query.pageSize);
	}

	async listSecurityEvents() {
		return securityEvents;
	}

	async listNotifications() {
		return notifications;
	}

	async listAuditEntries() {
		return auditEntries;
	}

	async listRuntimeLogs() {
		return runtimeLogs;
	}

	async getSmtpSettings() {
		return smtpSettings;
	}

	async updateSmtpSettings(input: UpdateSmtpSettingsInput) {
		smtpSettings = {
			host: input.host,
			port: input.port,
			username: input.username,
			hasPassword: input.password ? true : smtpSettings.hasPassword,
			encryption: input.encryption,
			fromEmail: input.fromEmail,
			fromName: input.fromName,
			replyToEmail: input.replyToEmail,
			active: input.active,
			updatedAt: new Date().toISOString()
		};

		return smtpSettings;
	}

	async sendTestEmail(recipient: string) {
		notifications.unshift({
			id: `not_${crypto.randomUUID()}`,
			type: 'smtp_test',
			channel: 'email',
			recipient,
			subject: 'CodeLinks SMTP Test',
			status: 'sent',
			createdAt: new Date().toISOString(),
			sentAt: new Date().toISOString()
		});
	}

	async setUserStatus(userId: string, status: UserStatus) {
		const user = users.find((candidate) => candidate.id === userId);

		if (!user) {
			error(404, 'User not found');
		}

		user.status = status;
	}

	async setUserRole(userId: string, role: string) {
		const user = users.find((candidate) => candidate.id === userId);

		if (!user) {
			error(404, 'User not found');
		}

		if (role === 'admin' || role === 'support' || role === 'auditor' || role === 'user') {
			user.primaryRole = role;
		}
	}
}
