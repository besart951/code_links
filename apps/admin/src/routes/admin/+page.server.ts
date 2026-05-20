import { getDashboard, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return {
			dashboard: await getDashboard(event)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
