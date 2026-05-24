<script lang="ts">
	import { Badge } from '@codelinks/ui-library/components/ui/badge';
	import * as Table from '@codelinks/ui-library/components/ui/table';
	import type { LoginAttempt } from '$lib/domain/auth-logs/types';

	let { attempts, maskIp = false }: { attempts: LoginAttempt[]; maskIp?: boolean } = $props();

	function formatIp(ipAddress: string) {
		if (!maskIp) return ipAddress;
		return ipAddress.replace(/\.\d+$/, '.xxx');
	}
</script>

<div class="rounded-lg border">
	<Table.Root>
		<Table.Header>
			<Table.Row>
				<Table.Head>Zeitpunkt</Table.Head>
				<Table.Head>Ergebnis</Table.Head>
				<Table.Head>E-Mail</Table.Head>
				<Table.Head>IP-Adresse</Table.Head>
				<Table.Head>Ort</Table.Head>
				<Table.Head>Gerät</Table.Head>
				<Table.Head>Grund</Table.Head>
			</Table.Row>
		</Table.Header>
		<Table.Body>
			{#each attempts as attempt (attempt.id)}
				<Table.Row>
					<Table.Cell>{new Date(attempt.occurredAt).toLocaleString('de-DE')}</Table.Cell>
					<Table.Cell>
						<Badge variant={attempt.success ? 'secondary' : 'destructive'}>
							{attempt.success ? 'Erfolgreich' : 'Fehlgeschlagen'}
						</Badge>
					</Table.Cell>
					<Table.Cell>{attempt.emailAttempted}</Table.Cell>
					<Table.Cell class="font-mono text-xs">{formatIp(attempt.ipAddress)}</Table.Cell>
					<Table.Cell>{attempt.city ? `${attempt.city}, ${attempt.countryCode}` : attempt.countryCode}</Table.Cell>
					<Table.Cell>{attempt.device.browser} / {attempt.device.os}</Table.Cell>
					<Table.Cell>{attempt.failureReason ?? '-'}</Table.Cell>
				</Table.Row>
			{/each}
		</Table.Body>
	</Table.Root>
</div>
