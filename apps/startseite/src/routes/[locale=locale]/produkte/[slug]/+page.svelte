<script lang="ts">
  import JsonLd from '$lib/components/JsonLd.svelte';
  import { alternateUrls, localizedAbsoluteUrl, localizedHref, site } from '$lib/site';
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();

  const product = $derived(data.product);
  const text = $derived(data.siteCopy);
  const productPath = $derived(`/produkte/${product.slug}`);
  const canonical = $derived(localizedAbsoluteUrl(productPath, data.locale));
  const productJsonLd = $derived({
    '@context': 'https://schema.org',
    '@type': 'SoftwareApplication',
    name: product.name,
    applicationCategory: 'BusinessApplication',
    operatingSystem: 'Web',
    url: product.appUrl,
    description: product.summary,
    offers: {
      '@type': 'Offer',
      availability:
        product.status === 'available'
          ? 'https://schema.org/InStock'
          : 'https://schema.org/PreOrder'
    }
  });
</script>

<svelte:head>
  <title>{product.name} | CodeLinks</title>
  <meta name="description" content={product.summary} />
  <link rel="canonical" href={canonical} />
  {#each alternateUrls(productPath) as alternate (alternate.locale)}
    <link rel="alternate" hreflang={alternate.hrefLang} href={alternate.href} />
  {/each}
  <link rel="alternate" hreflang="x-default" href={localizedAbsoluteUrl(productPath, 'de')} />
  <meta property="og:type" content="website" />
  <meta property="og:site_name" content="CodeLinks" />
  <meta property="og:title" content={`${product.name} | CodeLinks`} />
  <meta property="og:description" content={product.summary} />
  <meta property="og:url" content={canonical} />
  <meta name="twitter:card" content="summary_large_image" />
</svelte:head>

<JsonLd value={productJsonLd} />

<main class="min-h-screen px-5 pb-20 pt-24 sm:px-8 lg:px-16">
  <a
    class="font-black text-[var(--muted)] no-underline transition hover:text-[var(--accent)]"
    href={localizedHref('/', data.locale)}
  >
    {text.homeLabel}
  </a>

  <section class="max-w-5xl py-14 sm:py-20">
    <p class="mb-4 text-sm font-black uppercase tracking-normal text-[var(--accent)]">{product.key}</p>
    <h1 class="text-balance text-5xl font-black leading-none tracking-normal text-[var(--text)] sm:text-7xl lg:text-8xl">
      {product.name}
    </h1>
    <p class="mt-7 max-w-3xl text-balance text-2xl font-black leading-tight text-[var(--text)] sm:text-4xl">
      {product.headline}
    </p>
    <p class="mt-5 max-w-3xl text-lg leading-8 text-[var(--muted)]">{product.summary}</p>
    <a
      class="mt-8 inline-flex min-h-11 items-center rounded-lg bg-[var(--control-active-bg)] px-5 text-sm font-black text-[var(--control-active-text)] no-underline transition hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-[var(--accent)]"
      href={product.appUrl}
    >
      {product.status === 'available' ? text.openProduct : text.plannedAction}
    </a>
  </section>

  <section
    class="grid max-w-6xl gap-3 sm:grid-cols-2 lg:grid-cols-4"
    aria-label={`${product.name} ${text.featuresLabel}`}
  >
    {#each product.features as feature (feature)}
      <div class="rounded-lg border border-[var(--panel-border)] bg-[var(--panel)] p-5 font-black text-[var(--text)] shadow-[0_16px_36px_var(--panel-shadow)]">
        {feature}
      </div>
    {/each}
  </section>
</main>
