<script lang="ts">
  import { alternateUrls, jsonLd, localizedAbsoluteUrl, localizedHref, site } from '$lib/site';
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();

  const product = $derived(data.product);
  const text = $derived(data.siteCopy);
  const productPath = $derived(`/produkte/${product.slug}`);
  const canonical = $derived(localizedAbsoluteUrl(productPath, data.locale));
  const productJsonLd = $derived(jsonLd({
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
  }));
</script>

<svelte:head>
  <title>{product.name} | CodeLinks</title>
  <meta name="description" content={product.summary} />
  <link rel="canonical" href={canonical} />
  {#each alternateUrls(productPath) as alternate}
    <link rel="alternate" hreflang={alternate.hrefLang} href={alternate.href} />
  {/each}
  <link rel="alternate" hreflang="x-default" href={localizedAbsoluteUrl(productPath, 'de')} />
  <meta property="og:type" content="website" />
  <meta property="og:site_name" content="CodeLinks" />
  <meta property="og:title" content={`${product.name} | CodeLinks`} />
  <meta property="og:description" content={product.summary} />
  <meta property="og:url" content={canonical} />
  <meta name="twitter:card" content="summary_large_image" />
  {@html productJsonLd}
</svelte:head>

<main class="product-page">
  <a class="back" href={localizedHref('/', data.locale)}>{text.homeLabel}</a>

  <section class="hero">
    <p class="eyebrow">{product.key}</p>
    <h1>{product.name}</h1>
    <p class="lead">{product.headline}</p>
    <p class="summary">{product.summary}</p>
    <a class="cta" href={product.appUrl}>
      {product.status === 'available' ? text.openProduct : text.plannedAction}
    </a>
  </section>

  <section class="features" aria-label={`${product.name} ${text.featuresLabel}`}>
    {#each product.features as feature}
      <div>{feature}</div>
    {/each}
  </section>
</main>

<style>
  .product-page {
    min-height: 100vh;
    padding: 76px clamp(20px, 6vw, 72px) 72px;
  }

  .back {
    color: var(--muted);
    font-weight: 700;
    text-decoration: none;
  }

  .hero {
    max-width: 860px;
    padding: clamp(56px, 10vw, 120px) 0 48px;
  }

  .eyebrow {
    margin: 0 0 14px;
    color: var(--muted);
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0;
  }

  h1 {
    margin: 0;
    font-size: clamp(48px, 8vw, 104px);
    line-height: 0.95;
    letter-spacing: 0;
  }

  .lead {
    margin: 28px 0 0;
    max-width: 720px;
    font-size: clamp(22px, 3vw, 34px);
    line-height: 1.18;
  }

  .summary {
    max-width: 680px;
    margin: 20px 0 0;
    color: var(--muted);
    font-size: 18px;
    line-height: 1.6;
  }

  .cta {
    display: inline-flex;
    min-height: 44px;
    align-items: center;
    margin-top: 30px;
    border-radius: 8px;
    background: var(--control-active-bg);
    color: var(--control-active-text);
    padding: 0 18px;
    font-weight: 800;
    text-decoration: none;
  }

  .features {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 12px;
    max-width: 980px;
  }

  .features div {
    border: 1px solid var(--panel-border);
    border-radius: 8px;
    background: var(--panel);
    padding: 18px;
    font-weight: 700;
  }
</style>
