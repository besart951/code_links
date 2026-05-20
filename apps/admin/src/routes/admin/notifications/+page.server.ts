import { listNotifications, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return await listNotifications(event);
	} catch (err) {
		mapPlatformError(err);
	}
};
