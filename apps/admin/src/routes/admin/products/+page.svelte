<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import DataTableShell from '$lib/admin/components/DataTableShell.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="Products" description="Product access and health overview." />
<DataTableShell
	items={data.products}
	rowHref={(item) => `/admin/products/${item.product_key}`}
	columns={[
		{ key: 'name', label: 'Product', cell: (item) => item.name },
		{ key: 'tenants', label: 'Tenants', cell: (item) => item.active_tenants },
		{ key: 'users', label: 'Users', cell: (item) => item.active_users },
		{ key: 'subs', label: 'Subscriptions', cell: (item) => item.active_subscriptions },
		{ key: 'last', label: 'Last access', cell: (item) => item.last_access_at ?? 'none' }
	]}
/>
