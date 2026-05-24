import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminCommandRepository } from '$lib/application/ports/AdminCommandRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class SetUserRoleUseCase {
	constructor(private readonly repository: AdminCommandRepository) {}

	execute(actor: AdminActor, userId: string, role: string) {
		requirePermission(actor, 'admin.users.change_role');
		return this.repository.setUserRole(userId, role);
	}
}
