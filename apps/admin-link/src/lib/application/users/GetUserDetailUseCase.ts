import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminQueryRepository } from '$lib/application/ports/AdminQueryRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';

export class GetUserDetailUseCase {
	constructor(private readonly repository: AdminQueryRepository) {}

	execute(actor: AdminActor, userId: string) {
		requirePermission(actor, 'admin.users.read');
		return this.repository.getUserDetail(userId);
	}
}
