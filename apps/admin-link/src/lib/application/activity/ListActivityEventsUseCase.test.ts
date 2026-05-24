import { describe, expect, it } from 'vitest';
import { ListActivityEventsUseCase } from './ListActivityEventsUseCase';
import type { AdminQueryRepository } from '$lib/application/ports/AdminQueryRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

const admin: AdminActor = {
	id: 'admin',
	name: 'Admin',
	email: 'admin@example.com',
	roles: ['admin'],
	permissions: ['admin.auth_logs.read', 'admin.security_events.read', 'admin.notifications.read', 'admin.audit_entries.read']
};

describe('ListActivityEventsUseCase', () => {
	it('merges admin activity sources oldest first', async () => {
		const repository = {
			getDashboardSummary: async () => {
				throw new Error('not used');
			},
			listUsers: async () => {
				throw new Error('not used');
			},
			getUserDetail: async () => {
				throw new Error('not used');
			},
			listLoginAttempts: async () => ({
				items: [
					{
						id: 'login-1',
						userId: null,
						emailAttempted: 'admin@example.com',
						occurredAt: '2026-05-22T09:00:00.000Z',
						ipAddress: '192.0.2.1',
						countryCode: 'CH',
						city: 'Zurich',
						success: true,
						failureReason: null,
						authMethod: 'password',
						device: { browser: 'Chrome', os: 'Windows', userAgent: 'Chrome' },
						riskScore: 8
					}
				],
				total: 1,
				page: 1,
				pageSize: 100
			}),
			listSecurityEvents: async () => [
				{
					id: 'sec-1',
					userId: null,
					type: 'many_failed_logins',
					severity: 'high',
					status: 'open',
					summary: 'Viele fehlgeschlagene Logins',
					detectedAt: '2026-05-22T08:00:00.000Z',
					resolvedAt: null,
					sourceIpAddress: '192.0.2.2',
					countryCode: 'CH'
				}
			],
			listNotifications: async () => [],
			listAuditEntries: async () => [
				{
					id: 'audit-1',
					actorUserId: 'admin',
					action: 'admin.users.lock',
					targetType: 'user',
					targetId: 'usr-1',
					reason: '',
					ipAddress: '',
					createdAt: '2026-05-22T10:00:00.000Z'
				}
			],
			listRuntimeLogs: async () => [
				{
					id: 'runtime-1',
					occurredAt: '2026-05-22T11:00:00.000Z',
					level: 'fatal',
					source: 'auth-service',
					message: 'fatal: database unavailable',
					raw: '2026/05/22 11:00:00.000000 fatal: database unavailable'
				}
			],
			getSmtpSettings: async () => {
				throw new Error('not used');
			}
		} satisfies AdminQueryRepository;

		const events = await new ListActivityEventsUseCase(repository).execute(admin);

		expect(events.map((event) => event.id)).toEqual(['security:sec-1', 'auth:login-1', 'audit:audit-1', 'runtime:runtime-1']);
		expect(events.at(-1)?.tone).toBe('danger');
	});
});
