import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminCommandRepository } from '$lib/application/ports/AdminCommandRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';
import type { UpdateSmtpSettingsInput } from '$lib/domain/smtp/types';

export class UpdateSmtpSettingsUseCase {
	constructor(private readonly repository: AdminCommandRepository) {}

	execute(actor: AdminActor, input: UpdateSmtpSettingsInput) {
		requirePermission(actor, 'admin.smtp_settings.update');
		return this.repository.updateSmtpSettings(input);
	}
}
