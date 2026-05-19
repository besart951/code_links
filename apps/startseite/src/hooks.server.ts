import type { Handle } from '@sveltejs/kit';
import { htmlLang, localeFromPath } from '$lib/site';

export const handle: Handle = async ({ event, resolve }) => {
  const locale = localeFromPath(event.url.pathname);

  return resolve(event, {
    transformPageChunk: ({ html }) => html.replace('<html lang="de">', `<html lang="${htmlLang[locale]}">`)
  });
};
