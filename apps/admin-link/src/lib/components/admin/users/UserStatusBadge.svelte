<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index.js';
	import type { UserStatus } from '$lib/domain/users/types';

	let { status }: { status: UserStatus } = $props();
</script>

<Badge variant={status === 'locked' ? 'destructive' : 'secondary'}>
	{status === 'active' ? 'Aktiv' : status === 'disabled' ? 'Deaktiviert' : 'Gesperrt'}
</Badge>
