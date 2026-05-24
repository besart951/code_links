import { dev } from '$app/environment';
import { env } from '$env/dynamic/private';
import type { RequestEvent } from '@sveltejs/kit';
import { ListActivityEventsUseCase } from '$lib/application/activity/ListActivityEventsUseCase';
import { ListAuditEntriesUseCase } from '$lib/application/audit/ListAuditEntriesUseCase';
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
import type { AdminRepository } from '$lib/application/ports/AdminRepository';
import { AuthServiceAdminApiRepository } from '$lib/infrastructure/api/AuthServiceAdminApiRepository';
import { MockAdminRepository } from '$lib/server/mock-admin-repository';

function createAdminRepository(event: RequestEvent): AdminRepository {
	if (
		env.ADMIN_LINK_DATA_SOURCE === 'mock' ||
		env.ADMIN_LINK_MOCK_AUTH === 'true' ||
		(dev && env.ADMIN_LINK_MOCK_REPOSITORY === 'true')
	) {
		return new MockAdminRepository();
	}

	const baseUrl = (env.AUTH_API_BASE_URL ?? 'http://localhost:8080').replace(/\/$/, '');
	return new AuthServiceAdminApiRepository(baseUrl, event.fetch, {
		accept: 'application/json',
		cookie: event.request.headers.get('cookie') ?? ''
	});
}

export function createAdminContainer(event: RequestEvent) {
	const repository = createAdminRepository(event);
	const repositories = {
		query: repository,
		command: repository
	};

	return {
		getDashboardSummary: new GetDashboardSummaryUseCase(repositories.query),
		listUsers: new ListUsersUseCase(repositories.query),
		getUserDetail: new GetUserDetailUseCase(repositories.query),
		listLoginAttempts: new ListLoginAttemptsUseCase(repositories.query),
		listSecurityEvents: new ListSecurityEventsUseCase(repositories.query),
		listNotifications: new ListNotificationsUseCase(repositories.query),
		listAuditEntries: new ListAuditEntriesUseCase(repositories.query),
		listActivityEvents: new ListActivityEventsUseCase(repositories.query),
		getSmtpSettings: new GetSmtpSettingsUseCase(repositories.query),
		updateSmtpSettings: new UpdateSmtpSettingsUseCase(repositories.command),
		sendTestEmail: new SendTestEmailUseCase(repositories.command),
		setUserStatus: new SetUserStatusUseCase(repositories.command),
		setUserRole: new SetUserRoleUseCase(repositories.command)
	};
}
