<script lang="ts">
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Table from '$lib/components/ui/table/index.js';
	import type { SecurityEvent } from '$lib/domain/security/types';

	let { events, maskIp = false }: { events: SecurityEvent[]; maskIp?: boolean } = $props();

	function formatIp(ipAddress: string | null) {
		if (!ipAddress) return '-';
		return maskIp ? ipAddress.replace(/\.\d+$/, '.xxx') : ipAddress;
	}
</script>

<div class="rounded-lg border">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Schweregrad</Table.Head>
				<Table.Head>Ereignis</Table.Head>
				<Table.Head>Quelle</Table.Head>
				<Table.Head>Status</Table.Head>
				<Table.Head>Erkannt am</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each events as event (event.id)}
				<Table.Row>
					<Table.Cell>
						<Badge variant={event.severity === 'high' || event.severity === 'critical' ? 'destructive' : 'secondary'}>
							{event.severity}
						</Badge>
					</Table.Cell>
					<Table.Cell>{event.summary}</Table.Cell>
					<Table.Cell class="font-mono text-xs">{formatIp(event.sourceIpAddress)}</Table.Cell>
					<Table.Cell>{event.status === 'open' ? 'Offen' : 'Erledigt'}</Table.Cell>
					<Table.Cell>{new Date(event.detectedAt).toLocaleString('de-DE')}</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
