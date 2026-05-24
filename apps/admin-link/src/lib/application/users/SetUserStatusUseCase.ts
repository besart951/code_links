import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminCommandRepository } from '$lib/application/ports/AdminCommandRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';
import type { UserStatus } from '$lib/domain/users/types';

export class SetUserStatusUseCase {
	constructor(private readonly repository: AdminCommandRepository) {}

	execute(actor: AdminActor, userId: string, status: UserStatus) {
		requirePermission(actor, 'admin.users.update');
		return this.repository.setUserStatus(userId, status);
	}
}
