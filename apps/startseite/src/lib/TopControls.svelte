<script lang="ts">
  import { browser } from '$app/environment';
  import {
    getSiteCopy,
    localeLabels,
    localeNames,
    locales,
    localizedHref,
    type Locale
  } from '$lib/site';

  type Theme = 'light' | 'dark';

  let { locale, path }: { locale: Locale; path: string } = $props();
  let theme = $state<Theme>('light');
  let initialized = false;

  const text = $derived(getSiteCopy(locale));

  $effect(() => {
    if (!browser || initialized) {
      return;
    }

    const stored = window.localStorage.getItem('codelinks-theme');
    if (stored === 'light' || stored === 'dark') {
      theme = stored;
    } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      theme = 'dark';
    }

    initialized = true;
  });

  $effect(() => {
    if (!browser) {
      return;
    }

    document.documentElement.dataset.theme = theme;
    window.localStorage.setItem('codelinks-theme', theme);
  });
</script>

<nav class="top-controls" aria-label={text.controlsLabel}>
  <div class="language" aria-label={text.languageLabel}>
    {#each locales as target}
      <a
        href={localizedHref(path, target)}
        aria-label={localeNames[target]}
        aria-current={target === locale ? 'page' : undefined}
        class:active={target === locale}
      >
        {localeLabels[target]}
      </a>
    {/each}
  </div>

  <div class="theme" role="group" aria-label={text.themeLabel}>
    <button
      type="button"
      class:active={theme === 'light'}
      aria-pressed={theme === 'light'}
      onclick={() => (theme = 'light')}
    >
      {text.lightTheme}
    </button>
    <button
      type="button"
      class:active={theme === 'dark'}
      aria-pressed={theme === 'dark'}
      onclick={() => (theme = 'dark')}
    >
      {text.darkTheme}
    </button>
  </div>
</nav>

<style>
  .top-controls {
    position: fixed;
    z-index: 20;
    top: 16px;
    right: 16px;
    display: flex;
    max-width: calc(100vw - 32px);
    flex-wrap: wrap;
    justify-content: flex-end;
    gap: 8px;
    color: var(--text);
  }

  .language,
  .theme {
    display: inline-flex;
    overflow: hidden;
    border: 1px solid var(--control-border);
    border-radius: 8px;
    background: var(--control-bg);
    box-shadow: 0 12px 32px var(--panel-shadow);
    backdrop-filter: blur(12px);
  }

  a,
  button {
    min-width: 42px;
    min-height: 38px;
    border: 0;
    border-radius: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0 12px;
    background: transparent;
    color: var(--text);
    font: inherit;
    font-size: 13px;
    font-weight: 800;
    letter-spacing: 0;
    text-decoration: none;
    white-space: nowrap;
    cursor: pointer;
  }

  a + a,
  button + button {
    border-left: 1px solid var(--control-border);
  }

  .active {
    background: var(--control-active-bg);
    color: var(--control-active-text);
  }

  @media (max-width: 560px) {
    .top-controls {
      top: 10px;
      right: 10px;
      max-width: calc(100vw - 20px);
    }

    a,
    button {
      min-width: 36px;
      min-height: 34px;
      padding: 0 9px;
      font-size: 12px;
    }
  }
</style>
