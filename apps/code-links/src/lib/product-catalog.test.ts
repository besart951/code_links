import { describe, expect, it } from 'vitest';
import { products } from '@codelinks/config/products';

describe('product catalog', () => {
	it('exposes launch URLs for every portal product', () => {
		expect(products.map((product) => product.id)).toEqual(['infra-link', 'planer-link', 'loka-link']);
		expect(products.every((product) => product.appUrl.endsWith('.codelinks.localhost'))).toBe(true);
	});
});
