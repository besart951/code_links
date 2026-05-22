import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class SetUserRoleUseCase {
	constructor(private readonly repository: AdminReadRepository) {}

	execute(actor: AdminActor, userId: string, role: string) {
		requirePermission(actor, 'admin.users.change_role');
		return this.repository.setUserRole(userId, role);
	}
}
