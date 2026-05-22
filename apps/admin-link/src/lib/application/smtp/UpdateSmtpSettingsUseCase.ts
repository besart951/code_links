import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';
import type { UpdateSmtpSettingsInput } from '$lib/domain/smtp/types';

export class UpdateSmtpSettingsUseCase {
	constructor(private readonly repository: AdminReadRepository) {}

	execute(actor: AdminActor, input: UpdateSmtpSettingsInput) {
		requirePermission(actor, 'admin.smtp_settings.update');
		return this.repository.updateSmtpSettings(input);
	}
}
