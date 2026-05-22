import type { Paginated } from '$lib/domain/shared/pagination';
import type { LoginAttempt, LoginAttemptQuery } from '$lib/domain/auth-logs/types';
import type { SecurityEvent } from '$lib/domain/security/types';
import type { DashboardSummary } from '$lib/domain/statistics/types';
import type { Notification } from '$lib/domain/notifications/types';
import type { SmtpSettings, UpdateSmtpSettingsInput } from '$lib/domain/smtp/types';
import type { ManagedUserDetail, UserListItem, UserListQuery, UserStatus } from '$lib/domain/users/types';

export interface AdminReadRepository {
	getDashboardSummary(): Promise<DashboardSummary>;
	listUsers(query: UserListQuery): Promise<Paginated<UserListItem>>;
	getUserDetail(userId: string): Promise<ManagedUserDetail>;
	listLoginAttempts(query: LoginAttemptQuery): Promise<Paginated<LoginAttempt>>;
	listSecurityEvents(): Promise<SecurityEvent[]>;
	listNotifications(): Promise<Notification[]>;
	getSmtpSettings(): Promise<SmtpSettings>;
	updateSmtpSettings(input: UpdateSmtpSettingsInput): Promise<SmtpSettings>;
	sendTestEmail(recipient: string): Promise<void>;
	setUserStatus(userId: string, status: UserStatus): Promise<void>;
	setUserRole(userId: string, role: string): Promise<void>;
}
