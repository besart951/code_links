<script lang="ts">
	import type { AuditLogEntry } from '@codelinks/contracts';

	let { entries }: { entries: AuditLogEntry[] } = $props();
</script>

<div class="rounded-md border bg-background">
	{#if entries.length === 0}
		<p class="p-6 text-center text-sm text-muted-foreground">No audit events</p>
	{:else}
		<ol class="divide-y">
			{#each entries as entry}
				<li class="grid gap-2 p-4 sm:grid-cols-[11rem_1fr]">
					<time class="text-xs text-muted-foreground">{new Date(entry.created_at).toLocaleString('de-CH')}</time>
					<div class="min-w-0">
						<p class="truncate text-sm font-medium">{entry.action}</p>
						<p class="truncate text-xs text-muted-foreground">
							{entry.actor_user_id} -> {entry.target_type}:{entry.target_id}
						</p>
						{#if entry.reason}
							<p class="mt-2 rounded-md bg-muted px-2 py-1 text-xs">{entry.reason}</p>
						{/if}
					</div>
				</li>
			{/each}
		</ol>
	{/if}
</div>
