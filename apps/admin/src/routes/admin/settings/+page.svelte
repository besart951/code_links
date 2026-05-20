<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import DataTableShell from '$lib/admin/components/DataTableShell.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="Settings" description="Typed Platform settings with audit-aware mutations." />
<DataTableShell
	items={data.settings.items}
	columns={[
		{ key: 'label', label: 'Setting', cell: (item) => item.label },
		{ key: 'key', label: 'Key', cell: (item) => item.key },
		{ key: 'value', label: 'Value', cell: (item) => (item.sensitive ? 'hidden' : String(item.value)) },
		{ key: 'reason', label: 'Reason', cell: (item) => (item.requires_reason ? 'required' : 'not required') }
	]}
/>
