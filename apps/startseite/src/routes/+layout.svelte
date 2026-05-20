<script lang="ts">
  import PageTransitions from '$lib/components/PageTransitions.svelte';
  import SiteFooter from '$lib/components/SiteFooter.svelte';
  import TopControls from '$lib/TopControls.svelte';
  import type { Snippet } from 'svelte';
  import type { LayoutData } from './$types';
  import '../app.css';

  let { data, children }: { data: LayoutData; children: Snippet } = $props();
</script>

<PageTransitions />
<TopControls locale={data.locale} path={data.pathname} />
{@render children()}
<SiteFooter locale={data.locale} products={data.products} />
