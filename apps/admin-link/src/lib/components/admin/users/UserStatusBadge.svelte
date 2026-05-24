<script lang="ts">
	import { Badge } from '@codelinks/ui-library/components/ui/badge';
	import type { UserStatus } from '$lib/domain/users/types';

	let { status }: { status: UserStatus } = $props();
</script>

<Badge variant={status === 'locked' ? 'destructive' : 'secondary'}>
	{status === 'active' ? 'Aktiv' : status === 'disabled' ? 'Deaktiviert' : 'Gesperrt'}
</Badge>
