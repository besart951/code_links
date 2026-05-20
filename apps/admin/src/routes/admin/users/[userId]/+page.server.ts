import { getUser, listAudit, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return {
			user: await getUser(event, event.params.userId),
			audit: await listAudit(event)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
