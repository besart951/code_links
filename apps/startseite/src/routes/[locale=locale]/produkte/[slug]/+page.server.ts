import { error } from '@sveltejs/kit';
import { getProductBySlug, getSiteCopy, locales, parseLocale, products } from '$lib/site';

export const prerender = true;

export function entries() {
  return locales.flatMap((locale) => products.map((product) => ({ locale, slug: product.slug })));
}

export function load({ params }) {
  const locale = parseLocale(params.locale);
  const product = getProductBySlug(params.slug, locale);
  if (!product) {
    throw error(404, 'product_not_found');
  }

  return {
    locale,
    product,
    siteCopy: getSiteCopy(locale)
  };
}
