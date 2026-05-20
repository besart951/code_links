<script lang="ts">
	import { Button } from '@codelinks/ui/button';
	import { Input } from '@codelinks/ui/input';
	import { Search } from '@lucide/svelte';

	let query = $state('');

	function submit() {
		const value = query.trim();
		if (value) {
			window.location.href = `/admin/search?q=${encodeURIComponent(value)}`;
		}
	}
</script>

<div class="flex min-w-0 flex-1 items-center gap-2">
	<div class="relative max-w-xl flex-1">
		<Search class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
		<Input
			bind:value={query}
			class="h-9 bg-background pl-8"
			placeholder="Global search"
			onkeydown={(event) => {
				if (event.key === 'Enter') submit();
			}}
		/>
	</div>
	<Button type="button" variant="secondary" size="sm" onclick={submit}>Search</Button>
</div>
