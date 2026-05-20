<script lang="ts">
	import { Button } from '@codelinks/ui/button';
	import { Input } from '@codelinks/ui/input';
	import { Search, X } from '@lucide/svelte';

	let {
		query = '',
		placeholder = 'Search',
		action = ''
	}: {
		query?: string;
		placeholder?: string;
		action?: string;
	} = $props();
</script>

<form method="GET" {action} class="mb-3 flex min-w-0 flex-wrap items-center gap-2">
	<div class="relative min-w-64 flex-1">
		<Search class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
		<Input name="q" value={query} class="pl-8" {placeholder} />
	</div>
	<Button type="submit" variant="secondary">Apply</Button>
	{#if query}
		<Button href={action || '?'} variant="ghost" size="icon" aria-label="Clear filters">
			<X />
		</Button>
	{/if}
</form>
