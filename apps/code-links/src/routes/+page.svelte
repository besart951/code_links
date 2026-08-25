<script lang="ts">
	import { env } from '$env/dynamic/public';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import {
		ArrowRight,
		CheckCircle2,
		ExternalLink,
		Globe2,
		KeyRound,
		Languages,
		LayoutDashboard,
		LogIn,
		Menu,
		Moon,
		Network,
		ServerCog,
		ShieldCheck,
		Sparkles,
		Sun,
		UsersRound
	} from '@lucide/svelte';
	import { AuthManager } from '@codelinks/auth-client';
	import { products, type ProductId } from '@codelinks/config/products';
	import { Button } from '@codelinks/ui-library/components/ui/button';
	import * as DropdownMenu from '@codelinks/ui-library/components/ui/dropdown-menu';
	import * as Sheet from '@codelinks/ui-library/components/ui/sheet';
	import { resolveAuthAppUrl, resolveProductAppUrl } from '$lib/navigation-targets';
	import { getLocale, locales, localizeHref, type Locale } from '$lib/paraglide/runtime';
	import * as m from '$lib/paraglide/messages';

	type ThemePreference = 'system' | 'light' | 'dark';

	const auth = new AuthManager({
		authBaseUrl: env.PUBLIC_AUTH_BASE_URL ?? 'http://localhost:8080'
	});

	let mobileMenuOpen = $state(false);
	let themePreference = $state<ThemePreference>('system');
	let returnTo = $state('/');
	let browserHostname = $state('');

	const currentLocale = $derived(getLocale());
	const currentPath = $derived(`${page.url.pathname}${page.url.search}${page.url.hash}`);
	const authAppBaseUrl = $derived(resolveAuthAppUrl(env, browserHostname));
	const loginHref = $derived(`${authAppBaseUrl}/login?redirectTo=${encodeURIComponent(returnTo)}`);
	const heroStatus = $derived(
		auth.user ? m.hero_status_signed_in({ name: auth.user.name }) : m.hero_status_signed_out()
	);

	const localeLabels: Record<Locale, { label: string; short: string }> = {
		de: { label: 'Deutsch', short: 'DE' },
		en: { label: 'English', short: 'EN' },
		fr: { label: 'Francais', short: 'FR' },
		it: { label: 'Italiano', short: 'IT' },
		es: { label: 'Espanol', short: 'ES' }
	};

	const navItems = $derived([
		{ label: m.nav_who(), href: '#who' },
		{ label: m.nav_motivation(), href: '#motivation' },
		{ label: m.nav_story(), href: '#story' },
		{ label: m.nav_products(), href: '#products' }
	]);

	const languageOptions = $derived(
		locales.map((locale) => ({
			locale,
			label: localeLabels[locale].label,
			short: localeLabels[locale].short,
			href: localizeHref(currentPath || '/', { locale })
		}))
	);

	const platformSignals = $derived([
		{ label: m.visual_identity_label(), value: m.visual_identity_value() },
		{ label: m.visual_license_label(), value: m.visual_license_value() },
		{ label: m.visual_products_label(), value: m.visual_products_value() }
	]);

	const highlights = $derived([
		{
			title: m.feature_identity_title(),
			body: m.feature_identity_body()
		},
		{
			title: m.feature_license_title(),
			body: m.feature_license_body()
		},
		{
			title: m.feature_backends_title(),
			body: m.feature_backends_body()
		}
	]);

	const motivationItems = $derived([
		{ title: m.motivation_one_title(), body: m.motivation_one_body() },
		{ title: m.motivation_two_title(), body: m.motivation_two_body() },
		{ title: m.motivation_three_title(), body: m.motivation_three_body() }
	]);

	const storySteps = $derived([
		{ label: m.story_one_label(), title: m.story_one_title(), body: m.story_one_body() },
		{ label: m.story_two_label(), title: m.story_two_title(), body: m.story_two_body() },
		{ label: m.story_three_label(), title: m.story_three_title(), body: m.story_three_body() }
	]);

	const proofPoints = $derived([
		{ value: '01', label: m.proof_identity() },
		{ value: '03', label: m.proof_products() },
		{ value: '05', label: m.proof_locales() }
	]);

	onMount(() => {
		returnTo = window.location.href;
		browserHostname = window.location.hostname;

		const savedTheme = localStorage.getItem('theme');
		if (savedTheme === 'light' || savedTheme === 'dark') {
			themePreference = savedTheme;
		}

		applyTheme(themePreference);

		const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
		const syncSystemTheme = () => applyTheme(themePreference);
		mediaQuery.addEventListener('change', syncSystemTheme);

		void auth.refresh().catch(() => undefined);

		return () => mediaQuery.removeEventListener('change', syncSystemTheme);
	});

	function setTheme(preference: ThemePreference) {
		themePreference = preference;

		if (typeof localStorage !== 'undefined') {
			if (preference === 'system') {
				localStorage.removeItem('theme');
			} else {
				localStorage.setItem('theme', preference);
			}
		}

		applyTheme(preference);
	}

	function applyTheme(preference: ThemePreference) {
		if (typeof window === 'undefined') return;

		const dark =
			preference === 'dark' ||
			(preference === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches);

		document.documentElement.classList.toggle('dark', dark);
		document.documentElement.style.colorScheme = dark ? 'dark' : 'light';
	}

	async function handleLogout() {
		await auth.logout().catch(() => undefined);
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

	function productTone(productId: ProductId): string {
		switch (productId) {
			case 'infra-link':
				return 'border-cyan-200 bg-cyan-50 text-cyan-800 dark:border-cyan-400/30 dark:bg-cyan-400/10 dark:text-cyan-100';
			case 'planer-link':
				return 'border-amber-200 bg-amber-50 text-amber-800 dark:border-amber-400/30 dark:bg-amber-400/10 dark:text-amber-100';
			case 'loka-link':
				return 'border-emerald-200 bg-emerald-50 text-emerald-800 dark:border-emerald-400/30 dark:bg-emerald-400/10 dark:text-emerald-100';
		}
	}

	function productAppUrl(productId: ProductId, canonicalUrl: string): string {
		return resolveProductAppUrl(productId, canonicalUrl, env, browserHostname);
	}
</script>

<svelte:head>
	<title>{m.app_title()} | {m.meta_title_suffix()}</title>
	<meta name="description" content={m.meta_description()} />
</svelte:head>

<header class="sticky top-0 z-50 border-b border-zinc-200/80 bg-white/90 backdrop-blur-xl dark:border-white/10 dark:bg-zinc-950/90">
	<nav class="mx-auto flex min-h-16 w-full max-w-[90rem] items-center justify-between gap-3 px-4 sm:px-6 lg:px-8">
		<a href="#top" class="flex items-center gap-2 text-base font-semibold text-zinc-950 dark:text-white" aria-label={m.app_title()}>
			<span class="grid size-9 place-items-center rounded-md border border-zinc-200 bg-zinc-950 text-white shadow-sm dark:border-white/10 dark:bg-white dark:text-zinc-950">
				<Network class="size-4" />
			</span>
			<span>{m.app_title()}</span>
		</a>

		<div class="hidden items-center gap-1 md:flex">
			{#each navItems as item (item.href)}
				<a
					href={item.href}
					class="rounded-md px-3 py-2 text-sm font-medium text-zinc-600 transition hover:bg-zinc-100 hover:text-zinc-950 dark:text-zinc-300 dark:hover:bg-white/10 dark:hover:text-white"
				>
					{item.label}
				</a>
			{/each}
		</div>

		<div class="flex items-center gap-2">
			<DropdownMenu.Root>
				<DropdownMenu.Trigger
					aria-label={m.nav_language()}
					class="inline-flex h-9 min-w-9 items-center justify-center gap-1 rounded-md border border-zinc-200 bg-white px-2 text-sm font-medium text-zinc-700 transition hover:bg-zinc-100 dark:border-white/10 dark:bg-white/5 dark:text-zinc-100 dark:hover:bg-white/10"
				>
					<Languages class="size-4" />
					<span class="hidden sm:inline">{localeLabels[currentLocale].short}</span>
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end" class="w-56">
					{#each languageOptions as language (language.locale)}
						<DropdownMenu.Item class="p-0">
							<a
								href={language.href}
								data-sveltekit-reload
								aria-current={currentLocale === language.locale ? 'true' : undefined}
								aria-label={language.label}
								class="flex w-full items-center justify-between gap-3 rounded-2xl px-3 py-2 text-sm text-zinc-800 hover:bg-zinc-100 dark:text-zinc-100 dark:hover:bg-white/10"
							>
								<span>{language.label}</span>
								<span class="text-xs text-zinc-500 dark:text-zinc-400">{language.short}</span>
							</a>
						</DropdownMenu.Item>
					{/each}
				</DropdownMenu.Content>
			</DropdownMenu.Root>

			<DropdownMenu.Root>
				<DropdownMenu.Trigger
					aria-label={m.nav_theme()}
					class="inline-flex size-9 items-center justify-center rounded-md border border-zinc-200 bg-white text-zinc-700 transition hover:bg-zinc-100 dark:border-white/10 dark:bg-white/5 dark:text-zinc-100 dark:hover:bg-white/10"
				>
					{#if themePreference === 'dark'}
						<Moon class="size-4" />
					{:else if themePreference === 'light'}
						<Sun class="size-4" />
					{:else}
						<Sparkles class="size-4" />
					{/if}
				</DropdownMenu.Trigger>
				<DropdownMenu.Content align="end" class="w-48">
					<DropdownMenu.Item onclick={() => setTheme('system')}>
						<Sparkles class="size-4" />
						{m.nav_theme_system()}
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => setTheme('light')}>
						<Sun class="size-4" />
						{m.nav_theme_light()}
					</DropdownMenu.Item>
					<DropdownMenu.Item onclick={() => setTheme('dark')}>
						<Moon class="size-4" />
						{m.nav_theme_dark()}
					</DropdownMenu.Item>
				</DropdownMenu.Content>
			</DropdownMenu.Root>

			{#if auth.isAuthenticated}
				<span class="hidden max-w-36 truncate text-sm text-zinc-600 dark:text-zinc-300 lg:inline">{auth.user?.name}</span>
				<Button variant="secondary" onclick={handleLogout} class="hidden sm:inline-flex">{m.nav_logout()}</Button>
			{:else}
				<Button href={loginHref} class="hidden sm:inline-flex">
					<LogIn />
					{m.nav_login()}
				</Button>
			{/if}

			<Sheet.Root bind:open={mobileMenuOpen}>
				<Sheet.Trigger
					aria-label={m.nav_menu_open()}
					class="inline-flex size-9 items-center justify-center rounded-md border border-zinc-200 bg-white text-zinc-700 transition hover:bg-zinc-100 dark:border-white/10 dark:bg-white/5 dark:text-zinc-100 dark:hover:bg-white/10 md:hidden"
				>
					<Menu class="size-4" />
				</Sheet.Trigger>
				<Sheet.Content class="w-[min(22rem,92vw)] border-zinc-200 bg-white dark:border-white/10 dark:bg-zinc-950">
					<Sheet.Header class="border-b border-zinc-200 px-5 py-5 text-left dark:border-white/10">
						<Sheet.Title class="text-base font-semibold text-zinc-950 dark:text-white">{m.nav_menu_title()}</Sheet.Title>
						<Sheet.Description class="text-sm text-zinc-600 dark:text-zinc-400">{m.nav_menu_description()}</Sheet.Description>
					</Sheet.Header>
					<div class="grid gap-2 px-5 py-5">
						{#each navItems as item (item.href)}
							<a
								href={item.href}
								onclick={() => (mobileMenuOpen = false)}
								class="rounded-md px-3 py-3 text-sm font-medium text-zinc-700 hover:bg-zinc-100 dark:text-zinc-200 dark:hover:bg-white/10"
							>
								{item.label}
							</a>
						{/each}
						<div class="mt-4 border-t border-zinc-200 pt-4 dark:border-white/10">
							{#if auth.isAuthenticated}
								<p class="mb-3 truncate px-3 text-sm text-zinc-600 dark:text-zinc-300">{auth.user?.name}</p>
								<Button variant="secondary" onclick={handleLogout} class="w-full">{m.nav_logout()}</Button>
							{:else}
								<Button href={loginHref} class="w-full" onclick={() => (mobileMenuOpen = false)}>
									<LogIn />
									{m.nav_login()}
								</Button>
							{/if}
						</div>
					</div>
				</Sheet.Content>
			</Sheet.Root>
		</div>
	</nav>
</header>

<main id="top" class="bg-zinc-50 text-zinc-950 dark:bg-zinc-950 dark:text-white">
	<section class="relative isolate overflow-hidden border-b border-zinc-200 bg-zinc-50 dark:border-white/10 dark:bg-zinc-950">
		<div class="hero-scene pointer-events-none absolute inset-0" aria-hidden="true">
			<div class="scene-grid"></div>
			<div class="scene-panel scene-panel-a">
				<p>{m.visual_identity_label()}</p>
				<strong>{m.visual_identity_value()}</strong>
				<span></span>
			</div>
			<div class="scene-panel scene-panel-b">
				<p>{m.visual_license_label()}</p>
				<strong>{m.visual_license_value()}</strong>
				<span></span>
			</div>
			<div class="scene-panel scene-panel-c">
				<p>{m.visual_products_label()}</p>
				<strong>{m.visual_products_value()}</strong>
				<span></span>
			</div>
			<div class="scene-path scene-path-one"></div>
			<div class="scene-path scene-path-two"></div>
		</div>

		<div class="relative mx-auto flex min-h-[76svh] max-w-[90rem] flex-col justify-center px-4 py-20 sm:px-6 lg:px-8">
			<div class="max-w-3xl">
				<p class="inline-flex items-center gap-2 rounded-md border border-cyan-200 bg-cyan-50 px-3 py-1 text-sm font-medium text-cyan-900 dark:border-cyan-400/30 dark:bg-cyan-400/10 dark:text-cyan-100">
					<Sparkles class="size-4" />
					{m.hero_eyebrow()}
				</p>
				<h1 class="mt-6 max-w-3xl text-5xl font-semibold text-zinc-950 sm:text-6xl lg:text-7xl dark:text-white">{m.hero_title()}</h1>
				<p class="mt-6 max-w-2xl text-lg leading-8 text-zinc-700 dark:text-zinc-300">{m.hero_body()}</p>

				<div class="mt-8 flex flex-col gap-3 sm:flex-row sm:items-center">
					<Button href={loginHref} size="lg" class="w-full sm:w-auto">
						<LogIn />
						{m.hero_primary_cta()}
					</Button>
					<Button href="#products" size="lg" variant="secondary" class="w-full sm:w-auto">
						<LayoutDashboard />
						{m.hero_secondary_cta()}
					</Button>
				</div>

				<p class="mt-5 max-w-2xl text-sm leading-6 text-zinc-600 dark:text-zinc-400">{heroStatus}</p>
			</div>

			<div class="mt-10 grid max-w-3xl gap-3 sm:grid-cols-3">
				{#each proofPoints as point (point.label)}
					<div class="rounded-lg border border-zinc-200 bg-white/80 p-4 shadow-sm backdrop-blur dark:border-white/10 dark:bg-white/5">
						<p class="text-2xl font-semibold text-zinc-950 dark:text-white">{point.value}</p>
						<p class="mt-1 text-sm leading-5 text-zinc-600 dark:text-zinc-400">{point.label}</p>
					</div>
				{/each}
			</div>

			<div class="product-os-mobile mt-8 rounded-lg border border-zinc-200 bg-white p-4 shadow-sm dark:border-white/10 dark:bg-white/5 md:hidden">
				{#each platformSignals as signal (signal.label)}
					<div class="flex items-center justify-between border-b border-zinc-100 py-3 last:border-b-0 dark:border-white/10">
						<span class="text-sm text-zinc-600 dark:text-zinc-400">{signal.label}</span>
						<strong class="text-sm text-zinc-950 dark:text-white">{signal.value}</strong>
					</div>
				{/each}
			</div>
		</div>
	</section>

	<section id="who" class="border-b border-zinc-200 bg-white dark:border-white/10 dark:bg-zinc-900/60">
		<div class="mx-auto grid max-w-[90rem] gap-8 px-4 py-16 sm:px-6 lg:grid-cols-[0.9fr_1.1fr] lg:px-8 lg:py-24">
			<div>
				<h2 class="text-sm font-semibold text-cyan-700 dark:text-cyan-300">{m.section_who_kicker()}</h2>
				<p class="mt-3 text-3xl font-semibold text-zinc-950 sm:text-4xl dark:text-white">{m.section_who_title()}</p>
				<p class="mt-5 text-base leading-7 text-zinc-600 dark:text-zinc-300">{m.section_who_body()}</p>
			</div>

			<div id="work" class="grid gap-4 md:grid-cols-3">
				{#each highlights as highlight, index (highlight.title)}
					<article class="rounded-lg border border-zinc-200 bg-zinc-50 p-5 shadow-sm transition hover:-translate-y-1 hover:shadow-md dark:border-white/10 dark:bg-white/5">
						<div class="mb-5 inline-grid size-10 place-items-center rounded-md border border-zinc-200 bg-white text-zinc-950 dark:border-white/10 dark:bg-zinc-950 dark:text-white">
							{#if index === 0}
								<KeyRound class="size-5" />
							{:else if index === 1}
								<ShieldCheck class="size-5" />
							{:else}
								<ServerCog class="size-5" />
							{/if}
						</div>
						<h3 class="text-base font-semibold text-zinc-950 dark:text-white">{highlight.title}</h3>
						<p class="mt-3 text-sm leading-6 text-zinc-600 dark:text-zinc-400">{highlight.body}</p>
					</article>
				{/each}
			</div>
		</div>
	</section>

	<section id="motivation" class="border-b border-zinc-200 bg-zinc-50 dark:border-white/10 dark:bg-zinc-950">
		<div class="mx-auto max-w-[90rem] px-4 py-16 sm:px-6 lg:px-8 lg:py-24">
			<div class="max-w-3xl">
				<h2 class="text-sm font-semibold text-amber-700 dark:text-amber-300">{m.section_motivation_kicker()}</h2>
				<p class="mt-3 text-3xl font-semibold text-zinc-950 sm:text-4xl dark:text-white">{m.section_motivation_title()}</p>
				<p class="mt-5 text-base leading-7 text-zinc-600 dark:text-zinc-300">{m.section_motivation_body()}</p>
			</div>

			<div class="mt-10 grid gap-4 lg:grid-cols-3">
				{#each motivationItems as item (item.title)}
					<article class="rounded-lg border border-zinc-200 bg-white p-6 shadow-sm dark:border-white/10 dark:bg-white/5">
						<CheckCircle2 class="size-5 text-emerald-600 dark:text-emerald-300" />
						<h3 class="mt-5 text-lg font-semibold text-zinc-950 dark:text-white">{item.title}</h3>
						<p class="mt-3 text-sm leading-6 text-zinc-600 dark:text-zinc-400">{item.body}</p>
					</article>
				{/each}
			</div>
		</div>
	</section>

	<section id="story" class="border-b border-zinc-200 bg-white dark:border-white/10 dark:bg-zinc-900/60">
		<div class="mx-auto grid max-w-[90rem] gap-10 px-4 py-16 sm:px-6 lg:grid-cols-[0.8fr_1.2fr] lg:px-8 lg:py-24">
			<div>
				<h2 class="text-sm font-semibold text-rose-700 dark:text-rose-300">{m.section_story_kicker()}</h2>
				<p class="mt-3 text-3xl font-semibold text-zinc-950 sm:text-4xl dark:text-white">{m.section_story_title()}</p>
				<p class="mt-5 text-base leading-7 text-zinc-600 dark:text-zinc-300">{m.section_story_body()}</p>
			</div>

			<div class="grid gap-4">
				{#each storySteps as step (step.label)}
					<article class="story-step rounded-lg border border-zinc-200 bg-zinc-50 p-5 shadow-sm dark:border-white/10 dark:bg-white/5">
						<p class="text-sm font-semibold text-zinc-500 dark:text-zinc-400">{step.label}</p>
						<h3 class="mt-2 text-lg font-semibold text-zinc-950 dark:text-white">{step.title}</h3>
						<p class="mt-2 text-sm leading-6 text-zinc-600 dark:text-zinc-400">{step.body}</p>
					</article>
				{/each}
			</div>
		</div>
	</section>

	<section id="products" class="border-b border-zinc-200 bg-zinc-50 dark:border-white/10 dark:bg-zinc-950">
		<div class="mx-auto max-w-[90rem] px-4 py-16 sm:px-6 lg:px-8 lg:py-24">
			<div class="flex flex-col gap-5 lg:flex-row lg:items-end lg:justify-between">
				<div class="max-w-3xl">
					<p class="text-sm font-semibold text-emerald-700 dark:text-emerald-300">{m.section_products_kicker()}</p>
					<h2 class="mt-3 text-3xl font-semibold text-zinc-950 sm:text-4xl dark:text-white">{m.products_title()}</h2>
					<p class="mt-5 text-base leading-7 text-zinc-600 dark:text-zinc-300">{m.products_body()}</p>
				</div>
				<a href={loginHref} class="inline-flex items-center gap-2 text-sm font-medium text-zinc-950 hover:underline dark:text-white">
					{m.product_login_cta()}
					<ArrowRight class="size-4" />
				</a>
			</div>

			<div class="mt-10 grid gap-4 lg:grid-cols-3">
				{#each products as product (product.id)}
					<article class="flex min-h-72 flex-col justify-between rounded-lg border border-zinc-200 bg-white p-6 shadow-sm transition hover:-translate-y-1 hover:shadow-md dark:border-white/10 dark:bg-white/5">
						<div>
							<div class="flex items-start justify-between gap-4">
								<div>
									<p class={['inline-flex rounded-md border px-2.5 py-1 text-xs font-medium', productTone(product.id)]}>
										{product.id}
									</p>
									<h3 class="mt-4 text-xl font-semibold text-zinc-950 dark:text-white">{product.name}</h3>
								</div>
								<span class={['rounded-md px-2.5 py-1 text-xs font-medium', auth.hasLicense(product.id) ? 'bg-emerald-50 text-emerald-700 ring-1 ring-emerald-200 dark:bg-emerald-400/10 dark:text-emerald-200 dark:ring-emerald-400/30' : 'bg-zinc-100 text-zinc-600 ring-1 ring-zinc-200 dark:bg-white/10 dark:text-zinc-300 dark:ring-white/10']}>
									{auth.hasLicense(product.id) ? m.license_active() : m.license_missing()}
								</span>
							</div>
							<p class="mt-5 text-sm leading-6 text-zinc-600 dark:text-zinc-400">{productDescription(product.id)}</p>
						</div>

						<div class="mt-8">
							{#if auth.hasLicense(product.id)}
								<Button href={productAppUrl(product.id, product.appUrl)} class="w-full">
									<ExternalLink />
									{m.open_app()}
								</Button>
							{:else}
								<Button href={loginHref} variant="secondary" class="w-full">
									<LogIn />
									{m.product_login_cta()}
								</Button>
							{/if}
						</div>
					</article>
				{/each}
			</div>
		</div>
	</section>

	<section id="access" class="bg-white dark:bg-zinc-900/60">
		<div class="mx-auto grid max-w-[90rem] gap-8 px-4 py-16 sm:px-6 lg:grid-cols-[1fr_0.8fr] lg:px-8 lg:py-24">
			<div>
				<p class="text-sm font-semibold text-cyan-700 dark:text-cyan-300">{m.section_access_kicker()}</p>
				<h2 class="mt-3 text-3xl font-semibold text-zinc-950 sm:text-4xl dark:text-white">{m.access_title()}</h2>
				<p class="mt-5 max-w-3xl text-base leading-7 text-zinc-600 dark:text-zinc-300">{m.access_body()}</p>
				<div class="mt-8 flex flex-col gap-3 sm:flex-row">
					<Button href={loginHref} size="lg" class="w-full sm:w-auto">
						<LogIn />
						{m.access_primary_cta()}
					</Button>
					<Button href="#products" size="lg" variant="secondary" class="w-full sm:w-auto">
						<Globe2 />
						{m.access_secondary_cta()}
					</Button>
				</div>
			</div>

			<div class="rounded-lg border border-zinc-200 bg-zinc-50 p-6 shadow-sm dark:border-white/10 dark:bg-white/5">
				<UsersRound class="size-6 text-zinc-950 dark:text-white" />
				<h3 class="mt-5 text-lg font-semibold text-zinc-950 dark:text-white">{m.access_panel_title()}</h3>
				<p class="mt-3 text-sm leading-6 text-zinc-600 dark:text-zinc-400">{m.access_panel_body()}</p>
				<p class="mt-5 rounded-md border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-600 dark:border-white/10 dark:bg-zinc-950 dark:text-zinc-300">
					{heroStatus}
				</p>
			</div>
		</div>
	</section>
</main>

<style>
	.hero-scene {
		opacity: 0.9;
	}

	.scene-grid {
		position: absolute;
		inset: 0;
		background-image:
			linear-gradient(to right, rgba(63, 63, 70, 0.1) 1px, transparent 1px),
			linear-gradient(to bottom, rgba(63, 63, 70, 0.1) 1px, transparent 1px);
		background-size: 56px 56px;
		mask-image: linear-gradient(to bottom, transparent, black 12%, black 82%, transparent);
	}

	:global(.dark) .scene-grid {
		background-image:
			linear-gradient(to right, rgba(244, 244, 245, 0.08) 1px, transparent 1px),
			linear-gradient(to bottom, rgba(244, 244, 245, 0.08) 1px, transparent 1px);
	}

	.scene-panel {
		position: absolute;
		display: none;
		width: min(22rem, 28vw);
		border: 1px solid rgba(212, 212, 216, 0.9);
		border-radius: 8px;
		background: rgba(255, 255, 255, 0.86);
		box-shadow: 0 22px 70px rgba(15, 23, 42, 0.12);
		padding: 1rem;
		backdrop-filter: blur(18px);
		animation: panel-float 8s ease-in-out infinite;
	}

	:global(.dark) .scene-panel {
		border-color: rgba(255, 255, 255, 0.12);
		background: rgba(24, 24, 27, 0.78);
		box-shadow: 0 22px 70px rgba(0, 0, 0, 0.35);
	}

	.scene-panel p {
		color: rgb(82, 82, 91);
		font-size: 0.78rem;
		font-weight: 600;
		margin: 0;
	}

	.scene-panel strong {
		color: rgb(9, 9, 11);
		display: block;
		font-size: 1.1rem;
		margin-top: 0.35rem;
	}

	:global(.dark) .scene-panel p {
		color: rgb(212, 212, 216);
	}

	:global(.dark) .scene-panel strong {
		color: white;
	}

	.scene-panel span {
		display: block;
		height: 0.45rem;
		margin-top: 1rem;
		border-radius: 999px;
		background: linear-gradient(90deg, rgb(8, 145, 178), rgb(22, 163, 74), rgb(245, 158, 11));
		transform-origin: left;
		animation: signal-fill 4s ease-in-out infinite;
	}

	.scene-panel-a {
		top: 16%;
		right: 8%;
	}

	.scene-panel-b {
		top: 41%;
		right: 18%;
		animation-delay: -2s;
	}

	.scene-panel-c {
		right: 6%;
		bottom: 15%;
		animation-delay: -4s;
	}

	.scene-path {
		position: absolute;
		display: none;
		height: 2px;
		width: 18rem;
		background: linear-gradient(90deg, transparent, rgba(8, 145, 178, 0.8), transparent);
		animation: path-pulse 4s ease-in-out infinite;
	}

	.scene-path-one {
		top: 35%;
		right: 14%;
		transform: rotate(24deg);
	}

	.scene-path-two {
		right: 17%;
		bottom: 31%;
		transform: rotate(-18deg);
		animation-delay: -1.5s;
	}

	.story-step {
		position: relative;
		overflow: hidden;
	}

	.story-step::before {
		content: '';
		position: absolute;
		inset: 0 auto 0 0;
		width: 3px;
		background: rgb(8, 145, 178);
	}

	.product-os-mobile {
		animation: panel-float 8s ease-in-out infinite;
	}

	@media (min-width: 768px) {
		.scene-panel,
		.scene-path {
			display: block;
		}
	}

	@keyframes panel-float {
		0%,
		100% {
			transform: translateY(0);
		}
		50% {
			transform: translateY(-10px);
		}
	}

	@keyframes signal-fill {
		0%,
		100% {
			transform: scaleX(0.42);
			opacity: 0.62;
		}
		50% {
			transform: scaleX(1);
			opacity: 1;
		}
	}

	@keyframes path-pulse {
		0%,
		100% {
			opacity: 0.2;
		}
		50% {
			opacity: 0.72;
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.scene-panel,
		.scene-panel span,
		.scene-path,
		.product-os-mobile {
			animation: none;
		}

		article {
			transition: none;
		}
	}
</style>
