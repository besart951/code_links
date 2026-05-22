import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class GetSmtpSettingsUseCase {
	constructor(private readonly repository: AdminReadRepository) {}

	execute(actor: AdminActor) {
		requirePermission(actor, 'admin.smtp_settings.read');
		return this.repository.getSmtpSettings();
	}
}
