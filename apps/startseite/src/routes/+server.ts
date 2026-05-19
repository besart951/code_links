import { defaultLocale } from '$lib/site';

export function GET() {
  return new Response(null, {
    status: 308,
    headers: {
      location: `/${defaultLocale}`
    }
  });
}
