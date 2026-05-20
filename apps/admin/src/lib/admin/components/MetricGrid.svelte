<script lang="ts">
	import StatusBadge from './StatusBadge.svelte';
	import type { AdminDashboardMetric } from '@codelinks/contracts';

	let { metrics }: { metrics: AdminDashboardMetric[] } = $props();
	const toneLabels = {
		neutral: '',
		success: 'OK',
		warning: 'Watch',
		danger: 'Critical',
		info: 'Info'
	};
</script>

<div class="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
	{#each metrics as metric}
		<section class="rounded-md border bg-background p-4 shadow-xs">
			<div class="flex items-start justify-between gap-3">
				<p class="text-sm text-muted-foreground">{metric.label}</p>
				{#if metric.tone !== 'neutral'}
					<StatusBadge label={toneLabels[metric.tone]} tone={metric.tone} />
				{/if}
			</div>
			<p class="mt-3 text-2xl font-semibold tabular-nums">{metric.value.toLocaleString('de-CH')}</p>
		</section>
	{/each}
</div>
