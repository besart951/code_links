<script lang="ts">
	import LoginHistoryTable from '$lib/components/admin/login-history/LoginHistoryTable.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Login-Historie | CodeLinks Admin</title>
</svelte:head>

<section class="space-y-1">
	<h2 class="text-2xl font-semibold tracking-normal">Login-Historie</h2>
	<p class="text-muted-foreground text-sm">Alle erfolgreichen und fehlgeschlagenen Authentifizierungsereignisse.</p>
</section>

<LoginHistoryTable attempts={data.items} maskIp={data.maskIp} />
