<script lang="ts">
  import TopControls from '$lib/TopControls.svelte';
  import type { Snippet } from 'svelte';
  import type { LayoutData } from './$types';

  let { data, children }: { data: LayoutData; children: Snippet } = $props();
</script>

<TopControls locale={data.locale} path={data.pathname} />
{@render children()}

<style>
  :global(:root) {
    color-scheme: light;
    --page-bg: #f4f7f5;
    --text: #14211f;
    --muted: #52615d;
    --panel: #ffffff;
    --panel-border: #d7dfda;
    --panel-shadow: rgba(21, 45, 40, 0.12);
    --hero-bg: #dfe9e4;
    --hero-media:
      linear-gradient(120deg, rgba(20, 33, 31, 0.9), rgba(20, 33, 31, 0.18)),
      radial-gradient(circle at 75% 25%, #78b7a2, transparent 28%),
      linear-gradient(135deg, #173b36, #e2ede7);
    --hero-text: #ffffff;
    --control-bg: rgba(255, 255, 255, 0.88);
    --control-border: rgba(215, 223, 218, 0.9);
    --control-active-bg: #14211f;
    --control-active-text: #ffffff;
  }

  :global(html[data-theme='dark']) {
    color-scheme: dark;
    --page-bg: #101817;
    --text: #eef7f3;
    --muted: #b5c7bf;
    --panel: #182421;
    --panel-border: #2d403a;
    --panel-shadow: rgba(0, 0, 0, 0.32);
    --hero-bg: #101817;
    --hero-media:
      linear-gradient(120deg, rgba(8, 18, 16, 0.96), rgba(8, 18, 16, 0.28)),
      radial-gradient(circle at 75% 25%, #4f8f7b, transparent 28%),
      linear-gradient(135deg, #0d211f, #22352f);
    --hero-text: #f7fffb;
    --control-bg: rgba(16, 24, 23, 0.88);
    --control-border: rgba(82, 97, 93, 0.8);
    --control-active-bg: #eef7f3;
    --control-active-text: #101817;
  }

  :global(body) {
    margin: 0;
    min-width: 320px;
    font-family:
      Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI",
      sans-serif;
    color: var(--text);
    background: var(--page-bg);
  }
</style>
