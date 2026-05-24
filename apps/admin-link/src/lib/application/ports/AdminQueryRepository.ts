import type { LoginAttempt, LoginAttemptQuery } from '$lib/domain/auth-logs/types';
import type { AdminAuditEntry } from '$lib/domain/audit/types';
import type { SecurityEvent } from '$lib/domain/security/types';
import type { Paginated } from '$lib/domain/shared/pagination';
import type { DashboardSummary } from '$lib/domain/statistics/types';
import type { Notification } from '$lib/domain/notifications/types';
import type { RuntimeLogEntry } from '$lib/domain/runtime-logs/types';
import type { SmtpSettings } from '$lib/domain/smtp/types';
import type { ManagedUserDetail, UserListItem, UserListQuery } from '$lib/domain/users/types';

export interface AdminQueryRepository {
	getDashboardSummary(): Promise<DashboardSummary>;
	listUsers(query: UserListQuery): Promise<Paginated<UserListItem>>;
	getUserDetail(userId: string): Promise<ManagedUserDetail>;
	listLoginAttempts(query: LoginAttemptQuery): Promise<Paginated<LoginAttempt>>;
	listSecurityEvents(): Promise<SecurityEvent[]>;
	listNotifications(): Promise<Notification[]>;
	listAuditEntries(): Promise<AdminAuditEntry[]>;
	listRuntimeLogs(): Promise<RuntimeLogEntry[]>;
	getSmtpSettings(): Promise<SmtpSettings>;
}
