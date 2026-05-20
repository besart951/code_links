import { listAudit, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return {
			audit: await listAudit(event)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
