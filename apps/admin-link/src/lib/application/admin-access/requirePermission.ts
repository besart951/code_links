import { error } from '@sveltejs/kit';
import { hasPermission } from '$lib/domain/admin-access/permissions';
import type { AdminActor, AdminPermission } from '$lib/domain/admin-access/types';

export function requirePermission(actor: AdminActor, permission: AdminPermission) {
	if (!hasPermission(actor, permission)) {
		error(403, 'Missing admin permission');
	}
}
