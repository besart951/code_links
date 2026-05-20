import { listSettings, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return {
			settings: await listSettings(event)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
