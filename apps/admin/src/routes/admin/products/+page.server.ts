import { listProducts, mapPlatformError } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		return await listProducts(event);
	} catch (err) {
		mapPlatformError(err);
	}
};
