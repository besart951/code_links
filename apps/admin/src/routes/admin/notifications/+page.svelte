<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import DataTableShell from '$lib/admin/components/DataTableShell.svelte';
	import * as Tabs from '@codelinks/ui/tabs';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="Notifications" description="Templates, deliveries and retry queues." />
<Tabs.Root value="templates">
	<Tabs.List>
		<Tabs.Trigger value="templates">Templates</Tabs.Trigger>
		<Tabs.Trigger value="deliveries">Deliveries</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="templates" class="mt-4">
		<DataTableShell
			items={data.templates}
			columns={[
				{ key: 'key', label: 'Key', cell: (item) => item.key },
				{ key: 'channel', label: 'Channel', cell: (item) => item.channel },
				{ key: 'subject', label: 'Subject', cell: (item) => item.subject },
				{ key: 'enabled', label: 'Enabled', cell: (item) => (item.enabled ? 'yes' : 'no') }
			]}
		/>
	</Tabs.Content>
	<Tabs.Content value="deliveries" class="mt-4">
		<DataTableShell
			items={data.deliveries}
			columns={[
				{ key: 'event', label: 'Event', cell: (item) => item.event_key },
				{ key: 'channel', label: 'Channel', cell: (item) => item.channel },
				{ key: 'status', label: 'Status', cell: (item) => item.status },
				{ key: 'created', label: 'Created', cell: (item) => item.created_at }
			]}
		/>
	</Tabs.Content>
</Tabs.Root>
