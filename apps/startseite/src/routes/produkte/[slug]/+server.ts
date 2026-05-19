import { defaultLocale } from '$lib/site';

export function GET({ params }) {
  return new Response(null, {
    status: 308,
    headers: {
      location: `/${defaultLocale}/produkte/${params.slug}`
    }
  });
}
