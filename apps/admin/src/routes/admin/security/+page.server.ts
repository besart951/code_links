import { listSecurity, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return {
			security: await listSecurity(event)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
