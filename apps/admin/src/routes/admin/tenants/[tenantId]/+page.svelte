<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import AuditTimeline from '$lib/admin/components/AuditTimeline.svelte';
	import StatusBadge from '$lib/admin/components/StatusBadge.svelte';
	import * as Card from '@codelinks/ui/card';
	import * as Tabs from '@codelinks/ui/tabs';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title={data.tenant.name} description={`Tenant ${data.tenant.id}`} />
<Tabs.Root value="overview">
	<Tabs.List>
		<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
		<Tabs.Trigger value="audit">Audit</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="overview" class="mt-4 grid gap-3 md:grid-cols-2">
		<Card.Root>
			<Card.Header>
				<Card.Title>Status</Card.Title>
				<Card.Description>{data.tenant.tenant_type}</Card.Description>
			</Card.Header>
			<Card.Content class="space-y-2 text-sm">
				<StatusBadge label={data.tenant.status} tone={data.tenant.status === 'active' ? 'success' : 'warning'} />
				<p>Owner: {data.tenant.owner_user_id}</p>
				<p>Billing: {data.tenant.billing_email ?? 'not set'}</p>
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>Access</Card.Title>
				<Card.Description>{data.tenant.subscription_status}</Card.Description>
			</Card.Header>
			<Card.Content class="text-sm">
				<p>{data.tenant.active_products.join(', ') || 'No active products'}</p>
			</Card.Content>
		</Card.Root>
	</Tabs.Content>
	<Tabs.Content value="audit" class="mt-4">
		<AuditTimeline entries={data.audit.items} />
	</Tabs.Content>
</Tabs.Root>
