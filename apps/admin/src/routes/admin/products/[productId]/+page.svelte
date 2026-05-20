<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import MetricGrid from '$lib/admin/components/MetricGrid.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title={data.product.name} description={`Product ${data.product.product_key}`} />
<MetricGrid
	metrics={[
		{ key: 'tenants', label: 'Active tenants', value: data.product.active_tenants, tone: 'neutral' },
		{ key: 'users', label: 'Active users', value: data.product.active_users, tone: 'neutral' },
		{ key: 'subscriptions', label: 'Subscriptions', value: data.product.active_subscriptions, tone: 'neutral' },
		{ key: 'warnings', label: 'Warnings', value: data.product.warning_count, tone: data.product.warning_count > 0 ? 'warning' : 'success' }
	]}
/>
