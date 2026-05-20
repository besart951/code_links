import { listUsers, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return {
			users: await listUsers(event)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
