<script lang="ts">
	import LoginHistoryTable from '$lib/components/admin/login-history/LoginHistoryTable.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Fehlgeschlagene Logins | CodeLinks Admin</title>
</svelte:head>

<section class="space-y-1">
	<h2 class="text-2xl font-semibold tracking-normal">Fehlgeschlagene Login-Versuche</h2>
	<p class="text-muted-foreground text-sm">Falsche Passwörter, gesperrte Accounts und unbekannte E-Mail-Adressen.</p>
</section>

<LoginHistoryTable attempts={data.items} maskIp={data.maskIp} />
