import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminQueryRepository } from '$lib/application/ports/AdminQueryRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class GetSmtpSettingsUseCase {
	constructor(private readonly repository: AdminQueryRepository) {}

	execute(actor: AdminActor) {
		requirePermission(actor, 'admin.smtp_settings.read');
		return this.repository.getSmtpSettings();
	}
}
