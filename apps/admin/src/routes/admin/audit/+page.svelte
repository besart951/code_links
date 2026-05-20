<script lang="ts">
	import AdminPageHeader from '$lib/admin/components/AdminPageHeader.svelte';
	import AuditTimeline from '$lib/admin/components/AuditTimeline.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<AdminPageHeader title="Audit" description="Append-only admin and security activity." />
<AuditTimeline entries={data.audit.items} />
