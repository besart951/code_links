<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import DataTableShell from '$lib/admin/components/DataTableShell.svelte';
	import MetricGrid from '$lib/admin/components/MetricGrid.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="Dashboard" description="Platform overview, product health and security signals." />
<div class="space-y-4">
	<MetricGrid metrics={data.dashboard.metrics} />
	<section>
		<h2 class="mb-2 text-sm font-semibold">Products</h2>
		<DataTableShell
			items={data.dashboard.products}
			rowHref={(item) => `/admin/products/${item.product_key}`}
			columns={[
				{ key: 'name', label: 'Product', cell: (item) => item.name },
				{ key: 'tenants', label: 'Active tenants', cell: (item) => item.active_tenants },
				{ key: 'users', label: 'Active users', cell: (item) => item.active_users },
				{ key: 'subscriptions', label: 'Subscriptions', cell: (item) => item.active_subscriptions },
				{ key: 'warnings', label: 'Warnings', cell: (item) => item.warning_count }
			]}
		/>
	</section>
</div>
