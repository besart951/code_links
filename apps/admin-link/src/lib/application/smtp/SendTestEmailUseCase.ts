import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class SendTestEmailUseCase {
	constructor(private readonly repository: AdminReadRepository) {}

	execute(actor: AdminActor, recipient: string) {
		requirePermission(actor, 'admin.smtp_settings.update');
		return this.repository.sendTestEmail(recipient);
	}
}
