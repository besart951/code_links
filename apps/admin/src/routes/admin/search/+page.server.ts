import { mapPlatformError, searchAdmin } from '$lib/server/platformClient.js';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	const query = event.url.searchParams.get('q') ?? '';
	try {
		return {
			query,
			search: await searchAdmin(event, query)
		};
	} catch (err) {
		mapPlatformError(err);
	}
};
