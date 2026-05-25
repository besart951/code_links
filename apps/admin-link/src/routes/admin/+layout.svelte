<script lang="ts">
	import { page } from '$app/state';
	import AppSidebar from '$lib/components/app-sidebar.svelte';
	import SiteHeader from '$lib/components/site-header.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { LayoutData } from './$types';

	let { data, children }: { data: LayoutData; children: import('svelte').Snippet } = $props();

	const title = $derived.by(() => {
		if (page.url.pathname === '/admin') return 'Dashboard';
		if (page.url.pathname.startsWith('/admin/users')) return 'Benutzerverwaltung';
		if (page.url.pathname.startsWith('/admin/failed-logins')) return 'Fehlgeschlagene Logins';
		if (page.url.pathname.startsWith('/admin/logs')) return 'Alle Logs';
		if (page.url.pathname.startsWith('/admin/login-history')) return 'Login-Historie';
		if (page.url.pathname.startsWith('/admin/security')) return 'Security Events';
		if (page.url.pathname.startsWith('/admin/notifications')) return 'Benachrichtigungen';
		if (page.url.pathname.startsWith('/admin/settings/smtp')) return 'SMTP-Konfiguration';
		if (page.url.pathname.startsWith('/admin/audit')) return 'Audit Log';
		return 'CodeLinks Admin';
	});
</script>

<Sidebar.Provider>
	<AppSidebar admin={data.admin} pathname={page.url.pathname} />
	<Sidebar.Inset>
		<SiteHeader {title} adminMode={data.adminMode} />
		<main class="@container/main flex flex-1 flex-col gap-6 p-4 lg:p-6">
			{@render children()}
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
