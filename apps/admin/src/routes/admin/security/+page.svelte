<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import DataTableShell from '$lib/admin/components/DataTableShell.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="Security" description="Failed logins, session risk, MFA and key rotation signals." />
<DataTableShell
	items={data.security.items}
	columns={[
		{ key: 'type', label: 'Type', cell: (item) => item.event_type },
		{ key: 'severity', label: 'Severity', cell: (item) => item.severity },
		{ key: 'summary', label: 'Summary', cell: (item) => item.summary },
		{ key: 'created', label: 'Created', cell: (item) => item.created_at }
	]}
/>
