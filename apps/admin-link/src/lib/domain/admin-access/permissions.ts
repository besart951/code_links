import type { AdminActor, AdminPermission } from './types';

export function hasPermission(actor: AdminActor, permission: AdminPermission) {
	return actor.permissions.includes(permission);
}

export function canSeeRawIpAddress(actor: AdminActor) {
	return hasPermission(actor, 'admin.users.update');
}
