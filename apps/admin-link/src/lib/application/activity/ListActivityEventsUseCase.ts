import { hasPermission } from '$lib/domain/admin-access/permissions';
import type { AdminQueryRepository } from '$lib/application/ports/AdminQueryRepository';
import type { ActivityEvent, ActivityEventTone } from '$lib/domain/activity/types';
import type { AdminActor } from '$lib/domain/admin-access/types';
import type { LoginAttempt } from '$lib/domain/auth-logs/types';
import type { AdminAuditEntry } from '$lib/domain/audit/types';
import type { Notification } from '$lib/domain/notifications/types';
import type { RuntimeLogEntry } from '$lib/domain/runtime-logs/types';
import type { SecurityEvent } from '$lib/domain/security/types';

export class ListActivityEventsUseCase {
	constructor(private readonly repository: AdminQueryRepository) {}

	async execute(actor: AdminActor): Promise<ActivityEvent[]> {
		const [attempts, securityEvents, notifications, auditEntries, runtimeLogs] = await Promise.all([
			hasPermission(actor, 'admin.auth_logs.read')
				? this.repository.listLoginAttempts({ page: 1, pageSize: 100 })
				: Promise.resolve({ items: [], total: 0, page: 1, pageSize: 100 }),
			hasPermission(actor, 'admin.security_events.read') ? this.repository.listSecurityEvents() : Promise.resolve([]),
			hasPermission(actor, 'admin.notifications.read') ? this.repository.listNotifications() : Promise.resolve([]),
			hasPermission(actor, 'admin.audit_entries.read') ? this.repository.listAuditEntries() : Promise.resolve([]),
			hasPermission(actor, 'admin.audit_entries.read') ? this.repository.listRuntimeLogs() : Promise.resolve([])
		]);

		return [
			...attempts.items.map(loginAttemptEvent),
			...securityEvents.map(securityEventEvent),
			...notifications.map(notificationEvent),
			...auditEntries.map(auditEntryEvent),
			...runtimeLogs.map(runtimeLogEvent)
		]
			.sort((left, right) => Date.parse(left.occurredAt) - Date.parse(right.occurredAt))
			.slice(-200);
	}
}

function loginAttemptEvent(attempt: LoginAttempt): ActivityEvent {
	const tone: ActivityEventTone = attempt.success ? 'success' : attempt.riskScore >= 70 ? 'danger' : 'warning';

	return {
		id: `auth:${attempt.id}`,
		source: 'auth',
		tone,
		title: attempt.success ? 'Login erfolgreich' : 'Login fehlgeschlagen',
		message: attempt.success
			? `${attempt.emailAttempted} wurde authentifiziert.`
			: `${attempt.emailAttempted} wurde abgelehnt${attempt.failureReason ? ` (${attempt.failureReason})` : ''}.`,
		occurredAt: attempt.occurredAt,
		details: compactDetails([
			['IP', attempt.ipAddress],
			['Ort', attempt.city ? `${attempt.city}, ${attempt.countryCode}` : attempt.countryCode],
			['Gerät', `${attempt.device.browser} / ${attempt.device.os}`],
			['Risiko', String(attempt.riskScore)]
		])
	};
}

function securityEventEvent(event: SecurityEvent): ActivityEvent {
	return {
		id: `security:${event.id}`,
		source: 'security',
		tone: event.severity === 'critical' || event.severity === 'high' ? 'danger' : 'warning',
		title: 'Security Event',
		message: event.summary,
		occurredAt: event.detectedAt,
		details: compactDetails([
			['Schweregrad', event.severity],
			['Status', event.status],
			['Quelle', event.sourceIpAddress ?? ''],
			['Land', event.countryCode ?? '']
		])
	};
}

function notificationEvent(notification: Notification): ActivityEvent {
	return {
		id: `notification:${notification.id}`,
		source: 'notification',
		tone: notification.status === 'failed' ? 'danger' : notification.status === 'sent' ? 'info' : 'neutral',
		title: `Notification ${notification.status}`,
		message: `${notification.subject} an ${notification.recipient}`,
		occurredAt: notification.sentAt ?? notification.createdAt,
		details: compactDetails([
			['Typ', notification.type],
			['Kanal', notification.channel],
			['Status', notification.status]
		])
	};
}

function auditEntryEvent(entry: AdminAuditEntry): ActivityEvent {
	return {
		id: `audit:${entry.id}`,
		source: 'audit',
		tone: 'info',
		title: 'Admin Audit',
		message: `${entry.action} auf ${entry.targetType}:${entry.targetId}`,
		occurredAt: entry.createdAt,
		details: compactDetails([
			['Actor', entry.actorUserId],
			['Grund', entry.reason],
			['IP', entry.ipAddress]
		])
	};
}

function runtimeLogEvent(entry: RuntimeLogEntry): ActivityEvent {
	return {
		id: `runtime:${entry.id}`,
		source: 'runtime',
		tone: entry.level === 'fatal' ? 'danger' : 'neutral',
		title: entry.level === 'fatal' ? 'Runtime Fatal' : 'Runtime Log',
		message: entry.message,
		occurredAt: entry.occurredAt,
		details: compactDetails([
			['Quelle', entry.source],
			['Level', entry.level],
			['Raw', entry.raw]
		])
	};
}

function compactDetails(entries: [string, string][]) {
	return entries.filter(([, value]) => value !== '').map(([label, value]) => ({ label, value }));
}
