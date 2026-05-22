<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';
	import type { DashboardSummary } from '$lib/domain/statistics/types';

	let { trend }: { trend: DashboardSummary['trend'] } = $props();

	const maxValue = $derived(Math.max(...trend.map((item) => item.successful + item.failed), 1));
</script>

<Card.Root>
	<Card.Header>
		<Card.Title>Login-Trend</Card.Title>
		<Card.Description>Erfolgreiche und fehlgeschlagene Logins der letzten 7 Tage</Card.Description>
	</Card.Header>
	<Card.Content>
		<div class="grid h-64 grid-cols-7 items-end gap-3">
			{#each trend as item (item.date)}
				<div class="flex h-full flex-col justify-end gap-1">
					<div
						class="rounded-t-sm bg-destructive/70"
						style={`height: ${Math.max((item.failed / maxValue) * 100, item.failed > 0 ? 4 : 0)}%`}
						title={`${item.failed} fehlgeschlagen`}
					></div>
					<div
						class="rounded-t-sm bg-primary/80"
						style={`height: ${Math.max((item.successful / maxValue) * 100, item.successful > 0 ? 4 : 0)}%`}
						title={`${item.successful} erfolgreich`}
					></div>
					<div class="text-muted-foreground truncate text-center text-xs">{item.date.slice(5)}</div>
				</div>
			{/each}
		</div>
	</Card.Content>
</Card.Root>
