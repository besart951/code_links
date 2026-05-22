<script lang="ts">
	import type { Snippet } from 'svelte';

	type ButtonVariant = 'primary' | 'secondary' | 'ghost' | 'danger';

	interface Props {
		children?: Snippet;
		href?: string;
		type?: 'button' | 'submit' | 'reset';
		variant?: ButtonVariant;
		disabled?: boolean;
		class?: string;
		onclick?: (event: MouseEvent) => void;
	}

	let {
		children,
		href,
		type = 'button',
		variant = 'primary',
		disabled = false,
		class: className = '',
		onclick
	}: Props = $props();

	const base =
		'inline-flex min-h-10 items-center justify-center rounded-md px-4 py-2 text-sm font-medium transition focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50';

	const variants: Record<ButtonVariant, string> = {
		primary: 'bg-zinc-950 text-white hover:bg-zinc-800 focus-visible:ring-zinc-950',
		secondary: 'border border-zinc-300 bg-white text-zinc-950 hover:bg-zinc-100 focus-visible:ring-zinc-400',
		ghost: 'text-zinc-700 hover:bg-zinc-100 hover:text-zinc-950 focus-visible:ring-zinc-400',
		danger: 'bg-red-600 text-white hover:bg-red-700 focus-visible:ring-red-600'
	};

	let classes = $derived(`${base} ${variants[variant]} ${className}`.trim());
</script>

{#if href}
	<a class={classes} {href} aria-disabled={disabled}>
		{@render children?.()}
	</a>
{:else}
	<button class={classes} {type} {disabled} {onclick}>
		{@render children?.()}
	</button>
{/if}
