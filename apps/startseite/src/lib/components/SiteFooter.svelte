<script lang="ts">
  import {
    getSiteCopy,
    localizedHref,
    site,
    type Locale,
    type LocalizedProductPage
  } from '$lib/site';

  interface Props {
    locale: Locale;
    products: LocalizedProductPage[];
  }

  let { locale, products }: Props = $props();

  const text = $derived(getSiteCopy(locale));
  const year = new Date().getFullYear();
</script>

<footer class="border-t border-[var(--panel-border)] bg-[color-mix(in_srgb,var(--panel)_80%,var(--page-bg))]">
  <div class="mx-auto grid max-w-7xl gap-10 px-5 py-12 sm:px-8 lg:grid-cols-[1.25fr_2fr] lg:px-12">
    <div>
      <a class="text-2xl font-black text-[var(--text)] no-underline" href={localizedHref('/', locale)}>
        {site.title}
      </a>
      <p class="mt-4 max-w-md text-sm leading-6 text-[var(--muted)]">{text.footerTagline}</p>
    </div>

    <div class="grid gap-8 sm:grid-cols-3">
      <nav aria-label={text.footerProductsTitle}>
        <h2 class="text-sm font-black uppercase tracking-normal text-[var(--text)]">
          {text.footerProductsTitle}
        </h2>
        <div class="mt-4 grid gap-3">
          {#each products as product (product.slug)}
            <a
              class="text-sm font-bold text-[var(--muted)] no-underline transition hover:text-[var(--accent)]"
              href={localizedHref(`/produkte/${product.slug}`, locale)}
            >
              {product.name}
            </a>
          {/each}
        </div>
      </nav>

      <nav aria-label={text.footerPlatformTitle}>
        <h2 class="text-sm font-black uppercase tracking-normal text-[var(--text)]">
          {text.footerPlatformTitle}
        </h2>
        <div class="mt-4 grid gap-3">
          <a class="text-sm font-bold text-[var(--muted)] no-underline transition hover:text-[var(--accent)]" href="https://auth.codelinks.ch">
            {text.signIn}
          </a>
          <a class="text-sm font-bold text-[var(--muted)] no-underline transition hover:text-[var(--accent)]" href={localizedHref('/', locale) + '#products'}>
            ProductAccess
          </a>
          <a class="text-sm font-bold text-[var(--muted)] no-underline transition hover:text-[var(--accent)]" href={localizedHref('/', locale) + '#pricing'}>
            {text.pricingEyebrow}
          </a>
        </div>
      </nav>

      <nav aria-label={text.footerLegalTitle}>
        <h2 class="text-sm font-black uppercase tracking-normal text-[var(--text)]">
          {text.footerLegalTitle}
        </h2>
        <div class="mt-4 grid gap-3">
          <a class="text-sm font-bold text-[var(--muted)] no-underline transition hover:text-[var(--accent)]" href="mailto:hello@codelinks.ch">
            {text.footerContact}
          </a>
        </div>
      </nav>
    </div>
  </div>

  <div class="mx-auto flex max-w-7xl flex-wrap items-center justify-between gap-3 border-t border-[var(--panel-border)] px-5 py-5 text-xs font-bold text-[var(--muted)] sm:px-8 lg:px-12">
    <span>© {year} {site.title}</span>
    <span>codelinks.ch</span>
  </div>
</footer>
