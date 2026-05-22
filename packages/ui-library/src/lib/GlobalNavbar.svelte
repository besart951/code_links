<script lang="ts" module>
	export interface NavbarLink {
		label: string;
		href: string;
		active?: boolean;
	}
</script>

<script lang="ts">
	import Button from './Button.svelte';

	interface Props {
		brand?: string;
		links?: NavbarLink[];
		isAuthenticated?: boolean;
		userName?: string | null;
		loginLabel?: string;
		logoutLabel?: string;
		onLogin?: () => void;
		onLogout?: () => void;
	}

	let {
		brand = 'CodeLinks',
		links = [],
		isAuthenticated = false,
		userName = null,
		loginLabel = 'Login',
		logoutLabel = 'Logout',
		onLogin,
		onLogout
	}: Props = $props();
</script>

<header class="border-b border-zinc-200 bg-white/95 backdrop-blur">
	<nav class="mx-auto flex min-h-16 max-w-7xl items-center justify-between gap-4 px-4 sm:px-6 lg:px-8">
		<a href="/" class="text-base font-semibold tracking-normal text-zinc-950">{brand}</a>

		<div class="hidden items-center gap-1 md:flex">
			{#each links as link (link.href)}
				<a
					href={link.href}
					aria-current={link.active ? 'page' : undefined}
					class="rounded-md px-3 py-2 text-sm font-medium text-zinc-600 transition hover:bg-zinc-100 hover:text-zinc-950 aria-[current=page]:bg-zinc-100 aria-[current=page]:text-zinc-950"
				>
					{link.label}
				</a>
			{/each}
		</div>

		<div class="flex items-center gap-3">
			{#if isAuthenticated}
				{#if userName}
					<span class="hidden max-w-40 truncate text-sm text-zinc-600 sm:inline">{userName}</span>
				{/if}
				<Button variant="secondary" onclick={onLogout}>{logoutLabel}</Button>
			{:else}
				<Button onclick={onLogin}>{loginLabel}</Button>
			{/if}
		</div>
	</nav>
</header>
