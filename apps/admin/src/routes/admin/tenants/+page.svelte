<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import DataTableShell from '$lib/admin/components/DataTableShell.svelte';
	import FilterToolbar from '$lib/admin/components/FilterToolbar.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="Tenants" description="Tenant records are listed with minimized ownership and billing context." />
<FilterToolbar placeholder="Filter tenants" action="/admin/tenants" />
<DataTableShell
	items={data.tenants.items}
	rowHref={(item) => `/admin/tenants/${item.id}`}
	columns={[
		{ key: 'name', label: 'Name', cell: (item) => item.name },
		{ key: 'type', label: 'Type', cell: (item) => item.tenant_type },
		{ key: 'status', label: 'Status', cell: (item) => item.status },
		{ key: 'products', label: 'Products', cell: (item) => item.active_products.join(', ') || 'none' },
		{ key: 'risk', label: 'Risk', cell: (item) => item.risk_status }
	]}
/>
