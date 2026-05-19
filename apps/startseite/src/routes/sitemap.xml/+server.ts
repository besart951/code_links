import { alternateUrls, canonicalContentPaths, localizedAbsoluteUrl, locales } from '$lib/site';

export const prerender = true;

export function GET() {
  const urls = canonicalContentPaths()
    .flatMap((path) =>
      locales.map((locale) => {
        const alternates = [
          ...alternateUrls(path),
          { hrefLang: 'x-default', href: localizedAbsoluteUrl(path, 'de') }
        ]
          .map(
            (alternate) =>
              `    <xhtml:link rel="alternate" hreflang="${alternate.hrefLang}" href="${alternate.href}" />`
          )
          .join('\n');

        return `
  <url>
    <loc>${localizedAbsoluteUrl(path, locale)}</loc>
${alternates}
  </url>`;
      })
    )
    .join('');

  return new Response(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"
  xmlns:xhtml="http://www.w3.org/1999/xhtml">${urls}
</urlset>
`, {
    headers: {
      'content-type': 'application/xml; charset=utf-8'
    }
  });
}
