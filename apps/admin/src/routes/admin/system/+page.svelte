<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import DataTableShell from '$lib/admin/components/DataTableShell.svelte';
	import MetricGrid from '$lib/admin/components/MetricGrid.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="System" description="Operational health and pending platform messages." />
<div class="space-y-4">
	<MetricGrid
		metrics={[
			{ key: 'security', label: 'Security warnings', value: data.dashboard.security_warnings, tone: data.dashboard.security_warnings > 0 ? 'warning' : 'success' },
			{ key: 'messages', label: 'System messages', value: data.dashboard.open_system_messages, tone: data.dashboard.open_system_messages > 0 ? 'info' : 'success' }
		]}
	/>
	<DataTableShell
		items={data.security.items}
		columns={[
			{ key: 'type', label: 'Type', cell: (item) => item.event_type },
			{ key: 'severity', label: 'Severity', cell: (item) => item.severity },
			{ key: 'summary', label: 'Summary', cell: (item) => item.summary }
		]}
	/>
</div>
