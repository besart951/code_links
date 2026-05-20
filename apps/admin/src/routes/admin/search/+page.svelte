<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import DataTableShell from '$lib/admin/components/DataTableShell.svelte';
	import FilterToolbar from '$lib/admin/components/FilterToolbar.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="Search" description="Global Platform search across safe admin indexes." />
<FilterToolbar query={data.query} placeholder="Search tenants, users, audit logs" action="/admin/search" />
<DataTableShell
	items={data.search.results}
	columns={[
		{ key: 'type', label: 'Type', cell: (item) => item.type },
		{ key: 'title', label: 'Title', cell: (item) => item.title },
		{ key: 'subtitle', label: 'Subtitle', cell: (item) => item.subtitle },
		{ key: 'rank', label: 'Rank', cell: (item) => item.rank.toFixed(2) }
	]}
/>
