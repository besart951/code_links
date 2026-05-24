import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminQueryRepository } from '$lib/application/ports/AdminQueryRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class ListSecurityEventsUseCase {
	constructor(private readonly repository: AdminQueryRepository) {}

	execute(actor: AdminActor) {
		requirePermission(actor, 'admin.security_events.read');
		return this.repository.listSecurityEvents();
	}
}
