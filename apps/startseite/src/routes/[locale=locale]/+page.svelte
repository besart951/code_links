<script lang="ts">
  import JsonLd from '$lib/components/JsonLd.svelte';
  import LandingHero from '$lib/components/LandingHero.svelte';
  import PricingOverview from '$lib/components/PricingOverview.svelte';
  import ProductSelection from '$lib/components/ProductSelection.svelte';
  import { alternateUrls, localizedAbsoluteUrl, site } from '$lib/site';
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();

  const text = $derived(data.siteCopy);
  const canonical = $derived(localizedAbsoluteUrl('/', data.locale));
  const organizationJsonLd = $derived({
    '@context': 'https://schema.org',
    '@type': 'Organization',
    name: 'CodeLinks',
    url: site.url,
    description: text.description,
    sameAs: []
  });
</script>

<svelte:head>
  <title>{text.title}</title>
  <meta
    name="description"
    content={text.description}
  />
  <link rel="canonical" href={canonical} />
  {#each alternateUrls('/') as alternate (alternate.locale)}
    <link rel="alternate" hreflang={alternate.hrefLang} href={alternate.href} />
  {/each}
  <link rel="alternate" hreflang="x-default" href={localizedAbsoluteUrl('/', 'de')} />
  <meta property="og:type" content="website" />
  <meta property="og:site_name" content="CodeLinks" />
  <meta property="og:title" content={text.title} />
  <meta property="og:description" content={text.description} />
  <meta property="og:url" content={canonical} />
  <meta name="twitter:card" content="summary_large_image" />
</svelte:head>

<JsonLd value={organizationJsonLd} />

<main class="min-h-screen">
  <LandingHero
    eyebrow={text.eyebrow}
    title="CodeLinks"
    description={text.description}
    primary={{ label: text.signIn, href: 'https://auth.codelinks.ch', external: true }}
    secondary={{ label: text.viewProducts, href: '#products' }}
  />

  <ProductSelection locale={data.locale} products={data.products} {text} />
  <PricingOverview locale={data.locale} plans={data.pricingPlans} {text} />
</main>
