import { permissionsForRoles } from '$lib/domain/admin-access/permissions';
import type { AdminActor, AdminRole } from '$lib/domain/admin-access/types';

function normalizeRole(value: string | null): AdminRole {
	if (value === 'support' || value === 'auditor' || value === 'admin') {
		return value;
	}

	return 'admin';
}

export function createMockAdminActor(roleValue: string | null): AdminActor {
	const role = normalizeRole(roleValue);

	return {
		id: `dev-${role}`,
		email: `${role}@codelinks.localhost`,
		name: role === 'admin' ? 'CodeLinks Admin' : role === 'support' ? 'Support Team' : 'Audit Reader',
		roles: [role],
		permissions: permissionsForRoles([role])
	};
}
