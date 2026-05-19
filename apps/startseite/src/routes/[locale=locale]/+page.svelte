<script lang="ts">
  import { alternateUrls, jsonLd, localizedAbsoluteUrl, localizedHref, site } from '$lib/site';
  import type { PageData } from './$types';

  let { data }: { data: PageData } = $props();

  const text = $derived(data.siteCopy);
  const canonical = $derived(localizedAbsoluteUrl('/', data.locale));
  const organizationJsonLd = $derived(jsonLd({
    '@context': 'https://schema.org',
    '@type': 'Organization',
    name: 'CodeLinks',
    url: site.url,
    description: text.description,
    sameAs: []
  }));
</script>

<svelte:head>
  <title>{text.title}</title>
  <meta
    name="description"
    content={text.description}
  />
  <link rel="canonical" href={canonical} />
  {#each alternateUrls('/') as alternate}
    <link rel="alternate" hreflang={alternate.hrefLang} href={alternate.href} />
  {/each}
  <link rel="alternate" hreflang="x-default" href={localizedAbsoluteUrl('/', 'de')} />
  <meta property="og:type" content="website" />
  <meta property="og:site_name" content="CodeLinks" />
  <meta property="og:title" content={text.title} />
  <meta property="og:description" content={text.description} />
  <meta property="og:url" content={canonical} />
  <meta name="twitter:card" content="summary_large_image" />
  {@html organizationJsonLd}
</svelte:head>

<main class="page">
  <section class="hero">
    <div class="hero__media" aria-hidden="true">
      <div class="signal signal--one"></div>
      <div class="signal signal--two"></div>
      <div class="grid"></div>
    </div>
    <div class="hero__content">
      <p class="eyebrow">{text.eyebrow}</p>
      <h1>CodeLinks</h1>
      <p class="lead">{text.description}</p>
      <div class="actions">
        <a href="https://auth.codelinks.ch">{text.signIn}</a>
        <a href={localizedHref('/produkte/infra-link', data.locale)} class="secondary">{text.viewProducts}</a>
      </div>
    </div>
  </section>

  <section class="products" aria-label={text.productsLabel}>
    {#each data.products as product}
      <a class="product" href={localizedHref(`/produkte/${product.slug}`, data.locale)}>
        <span>{product.name}</span>
        <p>{product.summary}</p>
      </a>
    {/each}
  </section>
</main>

<style>
  .page {
    min-height: 100vh;
  }

  .hero {
    position: relative;
    min-height: 78vh;
    overflow: hidden;
    display: grid;
    place-items: center start;
    padding: clamp(32px, 8vw, 96px);
    background: var(--hero-bg);
  }

  .hero__media {
    position: absolute;
    inset: 0;
    background: var(--hero-media);
  }

  .grid {
    position: absolute;
    inset: 0;
    background-image:
      linear-gradient(rgba(255, 255, 255, 0.14) 1px, transparent 1px),
      linear-gradient(90deg, rgba(255, 255, 255, 0.14) 1px, transparent 1px);
    background-size: 56px 56px;
    mask-image: linear-gradient(90deg, black, transparent 75%);
  }

  .signal {
    position: absolute;
    width: 34vw;
    aspect-ratio: 1;
    border: 1px solid rgba(255, 255, 255, 0.36);
    transform: rotate(18deg);
  }

  .signal--one {
    right: 10vw;
    top: 12vh;
  }

  .signal--two {
    right: 24vw;
    bottom: -8vh;
  }

  .hero__content {
    position: relative;
    max-width: 760px;
    color: var(--hero-text);
  }

  .eyebrow {
    margin: 0 0 18px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0;
  }

  h1 {
    margin: 0;
    font-size: clamp(56px, 12vw, 132px);
    line-height: 0.88;
    letter-spacing: 0;
  }

  .lead {
    max-width: 620px;
    margin: 28px 0 0;
    font-size: clamp(18px, 2vw, 24px);
    line-height: 1.45;
  }

  .actions {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    margin-top: 34px;
  }

  .actions a {
    display: inline-flex;
    min-height: 44px;
    align-items: center;
    border-radius: 8px;
    padding: 0 18px;
    background: #ffffff;
    color: #14211f;
    font-weight: 700;
    text-decoration: none;
  }

  .actions .secondary {
    background: rgba(255, 255, 255, 0.16);
    color: #ffffff;
    border: 1px solid rgba(255, 255, 255, 0.4);
  }

  .products {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: 16px;
    max-width: 1120px;
    margin: -56px auto 0;
    padding: 0 24px 64px;
    position: relative;
  }

  .product {
    min-height: 132px;
    border-radius: 8px;
    border: 1px solid var(--panel-border);
    background: var(--panel);
    padding: 22px;
    color: var(--text);
    text-decoration: none;
    box-shadow: 0 18px 40px var(--panel-shadow);
  }

  .product span {
    display: block;
    font-size: 20px;
    font-weight: 800;
  }

  .product p {
    margin: 12px 0 0;
    color: var(--muted);
    line-height: 1.5;
  }

  @media (max-width: 680px) {
    .hero {
      align-items: end;
      min-height: 72vh;
      padding: 28px 22px 86px;
    }

    .products {
      margin-top: -40px;
    }
  }
</style>
