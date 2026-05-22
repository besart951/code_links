import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';
import type { UserListQuery } from '$lib/domain/users/types';

export class ListUsersUseCase {
	constructor(private readonly repository: AdminReadRepository) {}

	execute(actor: AdminActor, query: UserListQuery) {
		requirePermission(actor, 'admin.users.read');
		return this.repository.listUsers(query);
	}
}
