<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import AuditTimeline from '$lib/admin/components/AuditTimeline.svelte';
	import StatusBadge from '$lib/admin/components/StatusBadge.svelte';
	import { maskEmail } from '$lib/admin/security.js';
	import * as Card from '@codelinks/ui/card';
	import * as Tabs from '@codelinks/ui/tabs';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title={data.user.display_name} description={maskEmail(data.user.email)} />
<Tabs.Root value="overview">
	<Tabs.List>
		<Tabs.Trigger value="overview">Overview</Tabs.Trigger>
		<Tabs.Trigger value="audit">Audit</Tabs.Trigger>
	</Tabs.List>
	<Tabs.Content value="overview" class="mt-4 grid gap-3 md:grid-cols-3">
		<Card.Root>
			<Card.Header>
				<Card.Title>Status</Card.Title>
			</Card.Header>
			<Card.Content class="space-y-2 text-sm">
				<StatusBadge label={data.user.status} tone={data.user.status === 'active' ? 'success' : 'warning'} />
				<p>Email verified: {data.user.email_verified ? 'yes' : 'no'}</p>
				<p>MFA: {data.user.mfa_enabled ? 'enabled' : 'off'}</p>
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>Sessions</Card.Title>
			</Card.Header>
			<Card.Content class="text-sm">
				<p>{data.user.active_sessions} active</p>
				<p>{data.user.failed_login_count} failed logins</p>
			</Card.Content>
		</Card.Root>
		<Card.Root>
			<Card.Header>
				<Card.Title>Tenancy</Card.Title>
			</Card.Header>
			<Card.Content class="text-sm">
				<p>{data.user.tenant_count} memberships</p>
			</Card.Content>
		</Card.Root>
	</Tabs.Content>
	<Tabs.Content value="audit" class="mt-4">
		<AuditTimeline entries={data.audit.items} />
	</Tabs.Content>
</Tabs.Root>
