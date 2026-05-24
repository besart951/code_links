<script lang="ts">
	import { Badge } from '@codelinks/ui-library/components/ui/badge';
	import * as Table from '@codelinks/ui-library/components/ui/table';
	import type { Notification } from '$lib/domain/notifications/types';

	let { notifications }: { notifications: Notification[] } = $props();

	function badgeVariant(status: Notification['status']) {
		return status === 'failed' ? 'destructive' : 'outline';
	}
</script>

<div class="overflow-x-auto rounded-lg border">
	<Table.Root class="min-w-[44rem]">
		<Table.Header>
			<Table.Row>
				<Table.Head>Typ</Table.Head>
				<Table.Head>Empfänger</Table.Head>
				<Table.Head>Betreff</Table.Head>
				<Table.Head>Status</Table.Head>
				<Table.Head>Erstellt</Table.Head>
				<Table.Head>Gesendet</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each notifications as notification (notification.id)}
				<Table.Row>
					<Table.Cell class="font-medium">{notification.type}</Table.Cell>
					<Table.Cell>{notification.recipient}</Table.Cell>
					<Table.Cell>{notification.subject}</Table.Cell>
					<Table.Cell>
						<Badge variant={badgeVariant(notification.status)}>{notification.status}</Badge>
					</Table.Cell>
					<Table.Cell>{new Date(notification.createdAt).toLocaleString('de-DE')}</Table.Cell>
					<Table.Cell>{notification.sentAt ? new Date(notification.sentAt).toLocaleString('de-DE') : '-'}</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
