<script lang="ts">
	import SecurityEventsTable from '$lib/components/admin/security/SecurityEventsTable.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Security Events | CodeLinks Admin</title>
</svelte:head>

<section class="space-y-1">
	<h2 class="text-2xl font-semibold tracking-normal">Security Events</h2>
	<p class="text-muted-foreground text-sm">Verdächtige Aktivitäten aus Login-Mustern und Account-Statusereignissen.</p>
</section>

<SecurityEventsTable events={data.events} maskIp={data.maskIp} />
