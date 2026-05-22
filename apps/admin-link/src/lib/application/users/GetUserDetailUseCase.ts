import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class GetUserDetailUseCase {
	constructor(private readonly repository: AdminReadRepository) {}

	execute(actor: AdminActor, userId: string) {
		requirePermission(actor, 'admin.users.read');
		return this.repository.getUserDetail(userId);
	}
}
