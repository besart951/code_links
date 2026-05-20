<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import DataTableShell from '$lib/admin/components/DataTableShell.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="Subscriptions" description="Manual and provider-backed subscription state." />
<DataTableShell
	items={data.subscriptions.items}
	columns={[
		{ key: 'tenant', label: 'Tenant', cell: (item) => item.tenant_name },
		{ key: 'product', label: 'Product', cell: (item) => item.product_key },
		{ key: 'plan', label: 'Plan', cell: (item) => item.plan_name },
		{ key: 'status', label: 'Status', cell: (item) => item.status },
		{ key: 'period', label: 'Period end', cell: (item) => item.current_period_end ?? 'none' }
	]}
/>
