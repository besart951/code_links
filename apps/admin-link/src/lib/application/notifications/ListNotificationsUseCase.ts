import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminQueryRepository } from '$lib/application/ports/AdminQueryRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class ListNotificationsUseCase {
	constructor(private readonly repository: AdminQueryRepository) {}

	execute(actor: AdminActor) {
		requirePermission(actor, 'admin.notifications.read');
		return this.repository.listNotifications();
	}
}
