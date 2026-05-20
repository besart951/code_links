import { listProducts, mapPlatformError } from '$lib/server/platformClient.js';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async (event) => {
	try {
		const { products } = await listProducts(event);
		const product = products.find((item) => item.product_key === event.params.productId);
		if (!product) throw error(404, 'product_not_found');
		return { product };
	} catch (err) {
		mapPlatformError(err);
	}
};
