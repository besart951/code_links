import { getDashboard, listSecurity, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return {
			dashboard: await getDashboard(event),
			security: await listSecurity(event)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
