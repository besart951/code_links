import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';
import type { UserStatus } from '$lib/domain/users/types';

export class SetUserStatusUseCase {
	constructor(private readonly repository: AdminReadRepository) {}

	execute(actor: AdminActor, userId: string, status: UserStatus) {
		requirePermission(actor, 'admin.users.update');
		return this.repository.setUserStatus(userId, status);
	}
}
