<script lang="ts" generics="Item">
	import * as Table from '@codelinks/ui/table';
	import type { Snippet } from 'svelte';

	type Column<Item> = {
		key: string;
		label: string;
		cell: (item: Item) => string | number | null | undefined;
	};

	let {
		items,
		columns,
		emptyLabel = 'No records',
		rowHref,
		rowActions
	}: {
		items: Item[];
		columns: Column<Item>[];
		emptyLabel?: string;
		rowHref?: (item: Item) => string;
		rowActions?: Snippet<[Item]>;
	} = $props();
</script>

<div class="overflow-hidden rounded-md border bg-background">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				{#each columns as column}
					<Table.Head>{column.label}</Table.Head>
				{/each}
				{#if rowActions}
					<Table.Head resizable={false} class="w-12"></Table.Head>
				{/if}
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#if items.length === 0}
				<Table.Row>
					<Table.Cell colspan={columns.length + (rowActions ? 1 : 0)} class="h-20 text-center text-muted-foreground">
						{emptyLabel}
					</Table.Cell>
				</Table.Row>
			{:else}
				{#each items as item}
					<Table.Row>
						{#each columns as column}
							<Table.Cell>
								{#if rowHref}
									<a class="block min-w-0 truncate hover:underline" href={rowHref(item)}>
										{column.cell(item) ?? ''}
									</a>
								{:else}
									<span class="block min-w-0 truncate">{column.cell(item) ?? ''}</span>
								{/if}
							</Table.Cell>
						{/each}
						{#if rowActions}
							<Table.Cell class="text-right">
								{@render rowActions(item)}
							</Table.Cell>
						{/if}
					</Table.Row>
				{/each}
			{/if}
		</Table.Body>
	</Table.Root>
</div>
