<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import DataTableShell from '$lib/admin/components/DataTableShell.svelte';
	import FilterToolbar from '$lib/admin/components/FilterToolbar.svelte';
	import { maskEmail } from '$lib/admin/security.js';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="Users" description="User list views avoid unnecessary profile detail." />
<FilterToolbar placeholder="Filter users" action="/admin/users" />
<DataTableShell
	items={data.users.items}
	rowHref={(item) => `/admin/users/${item.id}`}
	columns={[
		{ key: 'email', label: 'Email', cell: (item) => maskEmail(item.email) },
		{ key: 'display', label: 'Name', cell: (item) => item.display_name },
		{ key: 'status', label: 'Status', cell: (item) => item.status },
		{ key: 'mfa', label: 'MFA', cell: (item) => (item.mfa_enabled ? 'enabled' : 'off') },
		{ key: 'sessions', label: 'Sessions', cell: (item) => item.active_sessions }
	]}
/>
