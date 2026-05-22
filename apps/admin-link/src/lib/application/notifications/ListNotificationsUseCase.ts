import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class ListNotificationsUseCase {
	constructor(private readonly repository: AdminReadRepository) {}

	execute(actor: AdminActor) {
		requirePermission(actor, 'admin.notifications.read');
		return this.repository.listNotifications();
	}
}
