<script lang="ts">
	import LoginHistoryTable from '$lib/components/admin/login-history/LoginHistoryTable.svelte';
	import UserStatusBadge from '$lib/components/admin/users/UserStatusBadge.svelte';
	import * as Card from '@codelinks/ui-library/components/ui/card';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>{data.user.name} | CodeLinks Admin</title>
</svelte:head>

<section class="flex flex-col gap-3 md:flex-row md:items-start md:justify-between">
	<div>
		<h2 class="text-2xl font-semibold tracking-normal">{data.user.name}</h2>
		<p class="text-muted-foreground text-sm">{data.user.email}</p>
	</div>
	<UserStatusBadge status={data.user.status} />
</section>

<div class="grid gap-4 md:grid-cols-3">
	<Card.Root>
		<Card.Header><Card.Title>Rollen</Card.Title></Card.Header>
		<Card.Content>{data.user.roles.map((role) => role.role).join(', ')}</Card.Content>
	</Card.Root>
	<Card.Root>
		<Card.Header><Card.Title>Lizenzen</Card.Title></Card.Header>
		<Card.Content>{data.user.productLicenses.length ? data.user.productLicenses.join(', ') : 'Keine'}</Card.Content>
	</Card.Root>
	<Card.Root>
		<Card.Header><Card.Title>Letzter Login</Card.Title></Card.Header>
		<Card.Content>{data.user.lastLoginAt ? new Date(data.user.lastLoginAt).toLocaleString('de-DE') : 'Nie'}</Card.Content>
	</Card.Root>
</div>

<section class="space-y-3">
	<h3 class="text-lg font-semibold">Login-Historie</h3>
	<LoginHistoryTable attempts={data.attempts} maskIp={data.maskIp} />
</section>
