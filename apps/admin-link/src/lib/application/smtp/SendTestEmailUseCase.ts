import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminCommandRepository } from '$lib/application/ports/AdminCommandRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class SendTestEmailUseCase {
	constructor(private readonly repository: AdminCommandRepository) {}

	execute(actor: AdminActor, recipient: string) {
		requirePermission(actor, 'admin.smtp_settings.update');
		return this.repository.sendTestEmail(recipient);
	}
}
