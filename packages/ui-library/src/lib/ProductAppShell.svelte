<script lang="ts">
	import Button from './Button.svelte';
	import GlobalNavbar from './GlobalNavbar.svelte';

	interface Props {
		brand: string;
		headline: string;
		body: string;
		status: string;
		refreshLabel: string;
		loadApiLabel: string;
		apiMessage?: string;
		isAuthenticated?: boolean;
		userName?: string | null;
		onRefresh: () => void;
		onLogout: () => void;
		onLoadAPI: () => void;
	}

	let {
		brand,
		headline,
		body,
		status,
		refreshLabel,
		loadApiLabel,
		apiMessage = '',
		isAuthenticated = false,
		userName = null,
		onRefresh,
		onLogout,
		onLoadAPI
	}: Props = $props();
</script>

<GlobalNavbar
	{brand}
	{isAuthenticated}
	{userName}
	onLogin={onRefresh}
	{onLogout}
	loginLabel={refreshLabel}
/>

<main class="mx-auto grid min-h-screen max-w-4xl content-center px-4 py-10">
	<section class="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
		<h1 class="text-4xl font-semibold text-zinc-950">{headline}</h1>
		<p class="mt-3 max-w-2xl text-base leading-7 text-zinc-600">{body}</p>
		<p class="mt-5 text-sm text-zinc-600">{status}</p>
		<div class="mt-6 flex flex-wrap gap-3">
			<Button onclick={onRefresh}>{refreshLabel}</Button>
			<Button variant="secondary" onclick={onLoadAPI}>{loadApiLabel}</Button>
		</div>
		{#if apiMessage}
			<p class="mt-5 rounded-md bg-zinc-50 px-3 py-2 text-sm text-zinc-700">{apiMessage}</p>
		{/if}
	</section>
</main>
