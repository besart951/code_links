<script lang="ts">
	import { env } from '$env/dynamic/public';
	import { onMount } from 'svelte';
	import { ExternalLink, LoaderCircle, LogIn, ShoppingCart } from '@lucide/svelte';
	import { AuthManager, type ProductLicense } from '@codelinks/auth-client';
	import { products, type ProductId } from '@codelinks/config/products';
	import { GlobalNavbar } from '@codelinks/ui-library';
	import { Button } from '@codelinks/ui-library/components/ui/button';
	import * as m from '$lib/paraglide/messages';

	const auth = new AuthManager({
		authBaseUrl: env.PUBLIC_AUTH_BASE_URL ?? 'http://localhost:8080'
	});

	let pendingProduct = $state<ProductId | null>(null);
	let statusMessage = $state('');

	onMount(() => {
		void auth.refresh();
	});

	const navLinks = $derived([{ label: m.nav_products(), href: '#products', active: true }]);
	const heroStatus = $derived(
		auth.user ? m.hero_status_signed_in({ name: auth.user.name }) : m.hero_status_signed_out()
	);

	async function handleLogin() {
		try {
			await auth.login({ email: 'demo@codelinks.dev', password: 'password' });
			statusMessage = '';
		} catch {
			statusMessage = m.login_failed();
		}
	}

	async function handleLogout() {
		await auth.logout();
		statusMessage = '';
	}

	async function handlePurchase(productId: ProductId) {
		pendingProduct = productId;
		statusMessage = '';

		try {
			if (!auth.isAuthenticated) {
				await auth.login({ email: 'demo@codelinks.dev', password: 'password' });
			}

			await auth.mockPurchase(productId as ProductLicense);
			const product = products.find((item) => item.id === productId);
			statusMessage = m.purchase_success({ productName: product?.name ?? productId });
		} catch {
			statusMessage = m.purchase_failed();
		} finally {
			pendingProduct = null;
		}
	}

	function productDescription(productId: ProductId): string {
		switch (productId) {
			case 'infra-link':
				return m.products_infra_description();
			case 'planer-link':
				return m.products_planer_description();
			case 'loka-link':
				return m.products_loka_description();
		}
	}
</script>

<svelte:head>
	<title>{m.app_title()}</title>
	<meta name="description" content={m.hero_body()} />
</svelte:head>

<GlobalNavbar
	brand={m.app_title()}
	links={navLinks}
	isAuthenticated={auth.isAuthenticated}
	userName={auth.user?.name}
	loginLabel={m.nav_login()}
	logoutLabel={m.nav_logout()}
	onLogin={handleLogin}
	onLogout={handleLogout}
/>

<main class="min-h-screen bg-zinc-50 text-zinc-950">
	<section class="border-b border-zinc-200 bg-white">
		<div class="mx-auto grid max-w-7xl gap-8 px-4 py-14 sm:px-6 lg:grid-cols-[1.05fr_0.95fr] lg:px-8 lg:py-20">
			<div class="flex flex-col justify-center">
				<p class="text-sm font-semibold uppercase tracking-normal text-blue-700">{m.hero_eyebrow()}</p>
				<h1 class="mt-4 text-5xl font-semibold tracking-normal text-zinc-950 sm:text-6xl">{m.hero_title()}</h1>
				<p class="mt-5 max-w-2xl text-lg leading-8 text-zinc-600">{m.hero_body()}</p>
				<div class="mt-8 flex flex-wrap items-center gap-3">
					<Button onclick={handleLogin} disabled={auth.isAuthenticated || auth.isLoading} size="lg">
						{#if auth.isLoading}
							<LoaderCircle class="animate-spin" />
						{:else}
							<LogIn />
						{/if}
						{m.nav_login()}
					</Button>
					<p class="text-sm text-zinc-600">{heroStatus}</p>
				</div>
				{#if statusMessage}
					<p class="mt-4 rounded-md border border-zinc-200 bg-zinc-50 px-3 py-2 text-sm text-zinc-700">{statusMessage}</p>
				{/if}
			</div>

			<div class="grid content-end gap-3">
				{#each products as product (product.id)}
					<div class="rounded-lg border border-zinc-200 bg-zinc-50 p-4 shadow-sm">
						<div class="flex items-start justify-between gap-4">
							<div>
								<h2 class="text-base font-semibold text-zinc-950">{product.name}</h2>
								<p class="mt-1 text-sm leading-6 text-zinc-600">{productDescription(product.id)}</p>
							</div>
							<span class="rounded-full bg-white px-2.5 py-1 text-xs font-medium text-zinc-600 ring-1 ring-zinc-200">
								{auth.hasLicense(product.id) ? m.license_active() : m.license_missing()}
							</span>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<section id="products" class="mx-auto max-w-7xl px-4 py-12 sm:px-6 lg:px-8">
		<div class="max-w-3xl">
			<h2 class="text-2xl font-semibold tracking-normal text-zinc-950">{m.products_title()}</h2>
			<p class="mt-3 text-base leading-7 text-zinc-600">{m.products_body()}</p>
		</div>

		<div class="mt-8 grid gap-4 md:grid-cols-3">
			{#each products as product (product.id)}
				<article class="flex min-h-64 flex-col justify-between rounded-lg border border-zinc-200 bg-white p-5 shadow-sm">
					<div>
						<div class="flex items-center justify-between gap-3">
							<h3 class="text-lg font-semibold text-zinc-950">{product.name}</h3>
							<span class="rounded-full px-2.5 py-1 text-xs font-medium {auth.hasLicense(product.id) ? 'bg-emerald-50 text-emerald-700 ring-1 ring-emerald-200' : 'bg-zinc-100 text-zinc-600 ring-1 ring-zinc-200'}">
								{auth.hasLicense(product.id) ? m.license_active() : m.license_missing()}
							</span>
						</div>
						<p class="mt-4 text-sm leading-6 text-zinc-600">{productDescription(product.id)}</p>
					</div>

					<div class="mt-6">
						{#if auth.hasLicense(product.id)}
							<Button href={product.appUrl} class="w-full">
								<ExternalLink />
								{m.open_app()}
							</Button>
						{:else}
							<Button
								class="w-full"
								variant="secondary"
								disabled={pendingProduct === product.id}
								onclick={() => handlePurchase(product.id)}
							>
								{#if pendingProduct === product.id}
									<LoaderCircle class="animate-spin" />
									{m.buying()}
								{:else}
									<ShoppingCart />
									{m.buy()}
								{/if}
							</Button>
						{/if}
					</div>
				</article>
			{/each}
		</div>
	</section>
</main>
