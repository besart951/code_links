<script lang="ts">
	import { env } from '$env/dynamic/public';
	import { onMount } from 'svelte';
	import { AuthManager } from '@codelinks/auth-client';
	import { Button, GlobalNavbar } from '@codelinks/ui-library';
	import * as m from '$lib/paraglide/messages';

	const productId = 'loka-link';
	const auth = new AuthManager({
		authBaseUrl: env.PUBLIC_AUTH_BASE_URL ?? 'http://auth.codelinks.localhost'
	});

	let apiMessage = $state('');

	onMount(() => {
		void auth.refresh();
	});

	async function loadProtectedAPI() {
		if (!auth.isAuthenticated) {
			await auth.refresh();
		}

		const response = await fetch('/api/me', { headers: auth.authorizationHeaders() });
		apiMessage = response.ok ? m.api_success({ productId }) : m.api_failed();
	}
</script>

<svelte:head>
	<title>{m.app_title()}</title>
</svelte:head>

<GlobalNavbar brand={m.app_title()} isAuthenticated={auth.isAuthenticated} userName={auth.user?.name} onLogin={() => auth.refresh()} onLogout={() => auth.logout()} loginLabel={m.refresh_session()} />

<main class="mx-auto grid min-h-screen max-w-4xl content-center px-4 py-10">
	<section class="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm">
		<h1 class="text-4xl font-semibold text-zinc-950">{m.headline()}</h1>
		<p class="mt-3 max-w-2xl text-base leading-7 text-zinc-600">{m.body()}</p>
		<p class="mt-5 text-sm text-zinc-600">{auth.user ? m.status_signed_in({ name: auth.user.name }) : m.status_signed_out()}</p>
		<div class="mt-6 flex flex-wrap gap-3">
			<Button onclick={() => auth.refresh()}>{m.refresh_session()}</Button>
			<Button variant="secondary" onclick={loadProtectedAPI}>{m.load_api()}</Button>
		</div>
		{#if apiMessage}
			<p class="mt-5 rounded-md bg-zinc-50 px-3 py-2 text-sm text-zinc-700">{apiMessage}</p>
		{/if}
	</section>
</main>
