import { getPricingPlans, getProducts, getSiteCopy, locales, parseLocale } from '$lib/site';
import type { PageLoad } from './$types';

export const prerender = true;

export function entries() {
  return locales.map((locale) => ({ locale }));
}

export const load: PageLoad = ({ params }) => {
  const locale = parseLocale(params.locale);

  return {
    locale,
    siteCopy: getSiteCopy(locale),
    pricingPlans: getPricingPlans(locale),
    products: getProducts(locale)
  };
};
