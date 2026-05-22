<script lang="ts">
	import NotificationsTable from '$lib/components/admin/notifications/NotificationsTable.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Benachrichtigungen | CodeLinks Admin</title>
</svelte:head>

<section class="space-y-1">
	<h2 class="text-2xl font-semibold tracking-normal">Benachrichtigungen</h2>
	<p class="text-muted-foreground text-sm">E-Mail-Benachrichtigungen und Versandstatus.</p>
</section>

<NotificationsTable notifications={data.notifications} />
