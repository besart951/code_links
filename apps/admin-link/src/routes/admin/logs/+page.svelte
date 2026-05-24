<script lang="ts">
	import ActivityFeed from '$lib/components/admin/activity/ActivityFeed.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Alle Logs | CodeLinks Admin</title>
</svelte:head>

<section class="space-y-1">
	<h2 class="text-2xl font-semibold tracking-normal">Alle Logs</h2>
	<p class="text-muted-foreground text-sm">Auth, Security, Notifications und Audit als chronologische Activity Events.</p>
</section>

<ActivityFeed events={data.events} />
