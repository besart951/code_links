<script lang="ts">
  import {
    localizedHref,
    type Locale,
    type LocalizedPricePlan,
    type SiteCopy
  } from '$lib/site';
  import SectionHeading from './SectionHeading.svelte';

  type PricingCopy = Pick<
    SiteCopy,
    | 'pricingEyebrow'
    | 'pricingTitle'
    | 'pricingDescription'
    | 'includedLabel'
    | 'priceFromLabel'
    | 'pricePerMonthLabel'
  >;

  interface Props {
    locale: Locale;
    plans: LocalizedPricePlan[];
    text: PricingCopy;
  }

  let { locale, plans, text }: Props = $props();

  const hrefForPlan = (plan: LocalizedPricePlan) =>
    plan.productSlug
      ? localizedHref(`/produkte/${plan.productSlug}`, locale)
      : 'mailto:hello@codelinks.ch?subject=CodeLinks%20Bundle';
</script>

<section id="pricing" class="mx-auto max-w-7xl px-5 pb-24 sm:px-8 lg:px-12" aria-labelledby="pricing-title">
  <SectionHeading
    id="pricing-title"
    align="center"
    eyebrow={text.pricingEyebrow}
    title={text.pricingTitle}
    description={text.pricingDescription}
  />

  <div class="mt-10 grid gap-5 lg:grid-cols-3">
    {#each plans as plan, index (plan.id)}
      <article
        class="price-plan rounded-lg border border-[var(--panel-border)] bg-[var(--panel)] p-6 shadow-[0_18px_42px_var(--panel-shadow)] motion-safe:animate-[panel-rise_580ms_cubic-bezier(0.16,1,0.3,1)_both]"
        class:featured={plan.highlighted}
        style={`animation-delay: ${index * 90}ms`}
      >
        <p class="text-sm font-black uppercase tracking-normal text-[var(--accent)]">{plan.audience}</p>
        <h3 class="mt-3 text-2xl font-black text-[var(--text)]">{plan.name}</h3>
        <p class="mt-3 min-h-14 text-sm leading-6 text-[var(--muted)]">{plan.summary}</p>

        <div class="mt-6 flex items-end gap-2">
          {#if plan.price.startsWith('CHF')}
            <span class="pb-1 text-sm font-black uppercase tracking-normal text-[var(--muted)]">
              {text.priceFromLabel}
            </span>
            <span class="text-4xl font-black leading-none text-[var(--text)]">{plan.price}</span>
            <span class="pb-1 text-sm font-bold text-[var(--muted)]">{text.pricePerMonthLabel}</span>
          {:else}
            <span class="text-4xl font-black leading-none text-[var(--text)]">{plan.price}</span>
          {/if}
        </div>

        <div class="mt-7">
          <p class="text-xs font-black uppercase tracking-normal text-[var(--muted)]">{text.includedLabel}</p>
          <ul class="mt-3 grid gap-3">
            {#each plan.features as feature (feature)}
              <li class="flex gap-3 text-sm font-bold leading-6 text-[var(--text)]">
                <span class="mt-2 h-2 w-2 shrink-0 rounded-full bg-[var(--accent)]"></span>
                <span>{feature}</span>
              </li>
            {/each}
          </ul>
        </div>

        <a
          class="mt-8 inline-flex min-h-11 w-full items-center justify-center rounded-lg bg-[var(--control-active-bg)] px-5 text-sm font-black text-[var(--control-active-text)] transition hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-[var(--accent)]"
          href={hrefForPlan(plan)}
        >
          {plan.cta}
        </a>
      </article>
    {/each}
  </div>
</section>

<style>
  .price-plan.featured {
    border-color: color-mix(in srgb, var(--accent) 62%, var(--panel-border));
    background:
      linear-gradient(180deg, color-mix(in srgb, var(--accent) 8%, transparent), transparent 36%),
      var(--panel);
    box-shadow:
      0 24px 58px var(--panel-shadow),
      inset 0 0 0 1px color-mix(in srgb, var(--accent) 34%, transparent);
  }
</style>
