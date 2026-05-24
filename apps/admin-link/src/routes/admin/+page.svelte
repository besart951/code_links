<script lang="ts">
	import LoginTrendChart from '$lib/components/admin/dashboard/LoginTrendChart.svelte';
	import MetricCards from '$lib/components/admin/dashboard/MetricCards.svelte';
	import LoginHistoryTable from '$lib/components/admin/login-history/LoginHistoryTable.svelte';
	import SecurityEventsTable from '$lib/components/admin/security/SecurityEventsTable.svelte';
	import * as Card from '@codelinks/ui-library/components/ui/card';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>CodeLinks Admin Dashboard</title>
</svelte:head>

<MetricCards summary={data.summary} />

<div class="grid gap-6 xl:grid-cols-[minmax(0,1fr)_360px]">
	<LoginTrendChart trend={data.summary.trend} />

	<Card.Root>
		<Card.Header>
			<Card.Title>Top Login-Länder</Card.Title>
			<Card.Description>Aggregiert aus den letzten Login-Versuchen</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-3">
			{#each data.summary.topCountries as country (country.countryCode)}
				<div class="flex items-center justify-between text-sm">
					<span>{country.countryCode}</span>
					<span class="font-medium">{country.count}</span>
				</div>
			{/each}
		</Card.Content>
	</Card.Root>
</div>

<section class="space-y-3">
	<h2 class="text-lg font-semibold">Letzte Aktivitäten</h2>
	<LoginHistoryTable attempts={data.summary.recentActivity} maskIp={data.maskIp} />
</section>

<section class="space-y-3">
	<h2 class="text-lg font-semibold">Hervorgehobene Security Events</h2>
	<SecurityEventsTable events={data.summary.highlightedEvents} maskIp={data.maskIp} />
</section>
