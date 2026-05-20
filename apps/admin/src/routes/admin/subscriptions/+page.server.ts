import { listSubscriptions, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return {
			subscriptions: await listSubscriptions(event)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
