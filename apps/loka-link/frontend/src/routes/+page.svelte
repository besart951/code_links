<script lang="ts">
	import { env } from '$env/dynamic/public';
	import { onMount } from 'svelte';
	import { AuthManager } from '@codelinks/auth-client';
	import { ProductAppShell } from '@codelinks/ui-library';
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

<ProductAppShell
	brand={m.app_title()}
	headline={m.headline()}
	body={m.body()}
	status={auth.user ? m.status_signed_in({ name: auth.user.name }) : m.status_signed_out()}
	refreshLabel={m.refresh_session()}
	loadApiLabel={m.load_api()}
	{apiMessage}
	isAuthenticated={auth.isAuthenticated}
	userName={auth.user?.name}
	onRefresh={() => auth.refresh()}
	onLogout={() => auth.logout()}
	onLoadAPI={loadProtectedAPI}
/>
