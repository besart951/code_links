<script lang="ts">
	import UserTable from '$lib/components/admin/users/UserTable.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Benutzerverwaltung | CodeLinks Admin</title>
</svelte:head>

<section class="space-y-1">
	<h2 class="text-2xl font-semibold tracking-normal">Benutzerverwaltung</h2>
	<p class="text-muted-foreground text-sm">
		Suche, Filter, Statusänderungen und Detailzugriff für Plattform-Benutzer.
	</p>
</section>

{#key JSON.stringify(data.query)}
	<UserTable users={data.items} total={data.total} query={data.query} admin={data.admin} maskIp={data.maskIp} />
{/key}
