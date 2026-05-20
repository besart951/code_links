import { listTenants, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return {
			tenants: await listTenants(event)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
