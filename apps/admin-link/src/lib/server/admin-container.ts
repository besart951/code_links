import { dev } from '$app/environment';
import { env } from '$env/dynamic/private';
import type { RequestEvent } from '@sveltejs/kit';
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
import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import { AuthServiceAdminApiRepository } from '$lib/infrastructure/api/AuthServiceAdminApiRepository';
import { MockAdminRepository } from '$lib/server/mock-admin-repository';

function createAdminRepository(event: RequestEvent): AdminReadRepository {
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

	return {
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
}
