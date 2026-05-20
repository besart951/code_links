<script lang="ts">
  import {
    localizedHref,
    type Locale,
    type LocalizedProductPage,
    type SiteCopy
  } from '$lib/site';
  import SectionHeading from './SectionHeading.svelte';

  type SelectionCopy = Pick<
    SiteCopy,
    | 'productSelectionEyebrow'
    | 'productSelectionTitle'
    | 'productSelectionDescription'
    | 'selectedProductLabel'
    | 'openProduct'
    | 'plannedAction'
  >;

  interface Props {
    locale: Locale;
    products: LocalizedProductPage[];
    text: SelectionCopy;
  }

  let { locale, products, text }: Props = $props();
  let activeSlug = $state<string | null>(null);

  const activeProductSlug = $derived(activeSlug ?? products[0]?.slug ?? '');
  const selectedProduct = $derived(
    products.find((product) => product.slug === activeProductSlug) ?? products[0]
  );
</script>

<section id="products" class="mx-auto max-w-7xl px-5 py-20 sm:px-8 lg:px-12" aria-labelledby="products-title">
  <SectionHeading
    id="products-title"
    eyebrow={text.productSelectionEyebrow}
    title={text.productSelectionTitle}
    description={text.productSelectionDescription}
  />

  <div class="mt-10 grid gap-5 lg:grid-cols-[0.95fr_1.05fr] lg:items-stretch">
    <div class="grid gap-3 sm:grid-cols-3 lg:grid-cols-1">
      {#each products as product, index (product.slug)}
        <button
          type="button"
          class="product-option rounded-lg border border-[var(--panel-border)] bg-[var(--panel)] p-5 text-left shadow-[0_16px_36px_var(--panel-shadow)] transition hover:-translate-y-0.5 hover:border-[var(--accent)] focus:outline-none focus:ring-2 focus:ring-[var(--accent)] motion-safe:animate-[panel-rise_520ms_cubic-bezier(0.16,1,0.3,1)_both]"
          class:selected={activeProductSlug === product.slug}
          style={`animation-delay: ${index * 70}ms`}
          aria-pressed={activeProductSlug === product.slug}
          onclick={() => (activeSlug = product.slug)}
        >
          <span class="text-xs font-black uppercase tracking-normal text-[var(--accent)]">{product.key}</span>
          <span class="mt-2 block text-xl font-black text-[var(--text)]">{product.name}</span>
          <span class="mt-3 block text-sm leading-6 text-[var(--muted)]">{product.summary}</span>
        </button>
      {/each}
    </div>

    {#if selectedProduct}
      <article
        class="rounded-lg border border-[var(--panel-border)] bg-[var(--panel)] p-6 shadow-[0_22px_50px_var(--panel-shadow)] sm:p-8 motion-safe:animate-[panel-rise_620ms_cubic-bezier(0.16,1,0.3,1)_both]"
        aria-label={text.selectedProductLabel}
      >
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div>
            <p class="text-sm font-black uppercase tracking-normal text-[var(--accent)]">
              {selectedProduct.status === 'available' ? text.openProduct : text.plannedAction}
            </p>
            <h3 class="mt-2 text-balance text-3xl font-black leading-tight text-(--text)">
              {selectedProduct.headline}
            </h3>
          </div>
          <span
            class="rounded-full border border-(--panel-border) px-3 py-1 text-xs font-black uppercase tracking-normal text-[var(--muted)]"
          >
            {selectedProduct.name}
          </span>
        </div>

        <p class="mt-5 max-w-3xl text-base leading-7 text-(--muted)">{selectedProduct.summary}</p>

        <div class="mt-7 grid gap-3 sm:grid-cols-2">
          {#each selectedProduct.features as feature (feature)}
            <div class="rounded-lg border border-(--panel-border) bg-[color-mix(in_srgb,var(--page-bg)_74%,var(--panel))] p-4 text-sm font-bold text-[var(--text)]">
              {feature}
            </div>
          {/each}
        </div>

        <a
          class="mt-8 inline-flex min-h-11 items-center rounded-lg bg-[var(--control-active-bg)] px-5 text-sm font-black text-[var(--control-active-text)] transition hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-[var(--accent)]"
          href={localizedHref(`/produkte/${selectedProduct.slug}`, locale)}
        >
          {selectedProduct.status === 'available' ? text.openProduct : text.plannedAction}
        </a>
      </article>
    {/if}
  </div>
</section>

<style>
  .product-option.selected {
    border-color: var(--accent);
    background: color-mix(in srgb, var(--panel) 78%, var(--accent) 8%);
    box-shadow:
      0 22px 48px var(--panel-shadow),
      inset 0 0 0 1px color-mix(in srgb, var(--accent) 38%, transparent);
  }
</style>
