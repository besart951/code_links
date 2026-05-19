import { canonicalPathFromLocalizedPath, localeFromPath, parseLocale } from '$lib/site';
import type { LayoutLoad } from './$types';

export const ssr = true;
export const csr = true;

export const load: LayoutLoad = ({ params, url }) => {
  const locale = parseLocale(params.locale ?? localeFromPath(url.pathname));

  return {
    locale,
    pathname: canonicalPathFromLocalizedPath(url.pathname)
  };
};
