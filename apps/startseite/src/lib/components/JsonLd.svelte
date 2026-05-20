<script lang="ts">
  interface Props {
    value: unknown;
  }

  let { value }: Props = $props();

  const json = $derived(JSON.stringify(value).replace(/</g, '\\u003c'));
</script>

<svelte:head>
  <script type="application/ld+json">{json}</script>
</svelte:head>
