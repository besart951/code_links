<script lang="ts">
  interface Props {
    eyebrow: string;
    title: string;
    description: string;
    id?: string;
    align?: 'left' | 'center';
  }

  let { eyebrow, title, description, id, align = 'left' }: Props = $props();
</script>

<div class={`max-w-3xl ${align === 'center' ? 'mx-auto text-center' : ''}`}>
  <p class="mb-3 text-sm font-black uppercase tracking-normal text-[var(--accent)]">{eyebrow}</p>
  <h2 {id} class="text-balance text-3xl font-black leading-tight text-[var(--text)] sm:text-4xl">
    {title}
  </h2>
  <p class="mt-4 text-base leading-7 text-[var(--muted)] sm:text-lg">{description}</p>
</div>
