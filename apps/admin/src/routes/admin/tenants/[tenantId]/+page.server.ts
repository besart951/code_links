import { getTenant, listAudit, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return {
			tenant: await getTenant(event, event.params.tenantId),
			audit: await listAudit(event)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
