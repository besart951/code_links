import { GetDashboardSummaryUseCase } from '$lib/application/dashboard/GetDashboardSummaryUseCase';
import { ListLoginAttemptsUseCase } from '$lib/application/login-statistics/ListLoginAttemptsUseCase';
import { ListNotificationsUseCase } from '$lib/application/notifications/ListNotificationsUseCase';
import { ListSecurityEventsUseCase } from '$lib/application/security-events/ListSecurityEventsUseCase';
import { GetSmtpSettingsUseCase } from '$lib/application/smtp/GetSmtpSettingsUseCase';
import { SendTestEmailUseCase } from '$lib/application/smtp/SendTestEmailUseCase';
import { UpdateSmtpSettingsUseCase } from '$lib/application/smtp/UpdateSmtpSettingsUseCase';
import { GetUserDetailUseCase } from '$lib/application/users/GetUserDetailUseCase';
import { ListUsersUseCase } from '$lib/application/users/ListUsersUseCase';
import { SetUserRoleUseCase } from '$lib/application/users/SetUserRoleUseCase';
import { SetUserStatusUseCase } from '$lib/application/users/SetUserStatusUseCase';
import { MockAdminRepository } from '$lib/server/mock-admin-repository';

const repository = new MockAdminRepository();

export const adminContainer = {
	getDashboardSummary: new GetDashboardSummaryUseCase(repository),
	listUsers: new ListUsersUseCase(repository),
	getUserDetail: new GetUserDetailUseCase(repository),
	listLoginAttempts: new ListLoginAttemptsUseCase(repository),
	listSecurityEvents: new ListSecurityEventsUseCase(repository),
	listNotifications: new ListNotificationsUseCase(repository),
	getSmtpSettings: new GetSmtpSettingsUseCase(repository),
	updateSmtpSettings: new UpdateSmtpSettingsUseCase(repository),
	sendTestEmail: new SendTestEmailUseCase(repository),
	setUserStatus: new SetUserStatusUseCase(repository),
	setUserRole: new SetUserRoleUseCase(repository)
};
