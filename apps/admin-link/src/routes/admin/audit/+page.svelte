<script lang="ts">
	import * as Table from '@codelinks/ui-library/components/ui/table';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>Audit Log | CodeLinks Admin</title>
</svelte:head>

<section class="space-y-1">
	<h2 class="text-2xl font-semibold tracking-normal">Audit Log</h2>
	<p class="text-muted-foreground text-sm">Nachvollziehbare Admin-Aktionen mit Actor, Ziel und Zeitpunkt.</p>
</section>

<div class="rounded-lg border">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Actor</Table.Head>
				<Table.Head>Aktion</Table.Head>
				<Table.Head>Ziel</Table.Head>
				<Table.Head>Zeitpunkt</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each data.entries as entry (entry.id)}
				<Table.Row>
					<Table.Cell>{entry.actor}</Table.Cell>
					<Table.Cell>{entry.action}</Table.Cell>
					<Table.Cell>{entry.target}</Table.Cell>
					<Table.Cell>{new Date(entry.createdAt).toLocaleString('de-DE')}</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
