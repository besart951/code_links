import { requirePermission } from '$lib/application/admin-access/requirePermission';
import type { AdminReadRepository } from '$lib/application/ports/AdminReadRepository';
import type { AdminActor } from '$lib/domain/admin-access/types';
import type { LoginAttemptQuery } from '$lib/domain/auth-logs/types';

export class ListLoginAttemptsUseCase {
	constructor(private readonly repository: AdminReadRepository) {}

	execute(actor: AdminActor, query: LoginAttemptQuery) {
		requirePermission(actor, 'admin.auth_logs.read');
		return this.repository.listLoginAttempts(query);
	}
}
