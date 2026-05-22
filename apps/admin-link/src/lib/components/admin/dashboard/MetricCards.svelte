<script lang="ts">
	import AlertTriangleIcon from '@tabler/icons-svelte/icons/alert-triangle';
	import KeyIcon from '@tabler/icons-svelte/icons/key';
	import LockIcon from '@tabler/icons-svelte/icons/lock';
	import LoginIcon from '@tabler/icons-svelte/icons/login';
	import MailIcon from '@tabler/icons-svelte/icons/mail';
	import UsersIcon from '@tabler/icons-svelte/icons/users';
	import { Badge } from '$lib/components/ui/badge/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import type { DashboardSummary } from '$lib/domain/statistics/types';

	let { summary }: { summary: DashboardSummary } = $props();

	const failedRate = $derived(
		summary.loginAttempts.total === 0
			? 0
			: Math.round((summary.loginAttempts.failed / summary.loginAttempts.total) * 100)
	);
</script>

<div
	class="*:data-[slot=card]:from-primary/5 *:data-[slot=card]:to-card dark:*:data-[slot=card]:bg-card grid grid-cols-1 gap-4 *:data-[slot=card]:bg-gradient-to-t *:data-[slot=card]:shadow-xs md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-6"
>
	<Card.Root>
		<Card.Header>
			<Card.Description>Benutzer gesamt</Card.Description>
			<Card.Title class="flex items-center gap-2 text-3xl font-semibold tabular-nums">
				<UsersIcon class="size-6 text-primary" />
				{summary.users.total}
			</Card.Title>
		</Card.Header>
		<Card.Footer class="text-muted-foreground text-sm">
			{summary.users.newLast7Days} neu in den letzten 7 Tagen
		</Card.Footer>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Description>Aktive Benutzer</Card.Description>
			<Card.Title class="text-3xl font-semibold tabular-nums">{summary.users.active}</Card.Title>
		</Card.Header>
		<Card.Footer class="text-muted-foreground text-sm">
			{summary.users.newLast30Days} neue Accounts in 30 Tagen
		</Card.Footer>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Description>Fehlgeschlagene Logins</Card.Description>
			<Card.Title class="flex items-center gap-2 text-3xl font-semibold tabular-nums">
				<LoginIcon class="size-6 text-destructive" />
				{summary.loginAttempts.failed}
			</Card.Title>
			<Card.Action>
				<Badge variant={failedRate > 40 ? 'destructive' : 'outline'}>{failedRate}%</Badge>
			</Card.Action>
		</Card.Header>
		<Card.Footer class="text-muted-foreground text-sm">
			{summary.loginAttempts.successful} erfolgreiche Logins im Zeitraum
		</Card.Footer>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Description>Passwort-Resets</Card.Description>
			<Card.Title class="flex items-center gap-2 text-3xl font-semibold tabular-nums">
				<KeyIcon class="size-6 text-primary" />
				{summary.passwordResetRequests}
			</Card.Title>
		</Card.Header>
		<Card.Footer class="text-muted-foreground text-sm">Anfragen der letzten 24 Stunden</Card.Footer>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Description>Notifications</Card.Description>
			<Card.Title class="flex items-center gap-2 text-3xl font-semibold tabular-nums">
				<MailIcon class="size-6 text-primary" />
				{summary.notifications}
			</Card.Title>
		</Card.Header>
		<Card.Footer class="text-muted-foreground text-sm">E-Mail-Kanal in v1</Card.Footer>
	</Card.Root>

	<Card.Root class={summary.security.openEvents > 0 ? 'border-destructive/50' : ''}>
		<Card.Header>
			<Card.Description>Security Events</Card.Description>
			<Card.Title class="flex items-center gap-2 text-3xl font-semibold tabular-nums">
				{#if summary.security.openEvents > 0}
					<AlertTriangleIcon class="size-6 text-destructive" />
				{:else}
					<LockIcon class="size-6 text-primary" />
				{/if}
				{summary.security.openEvents}
			</Card.Title>
			<Card.Action>
				<Badge variant={summary.security.suspiciousAttempts > 0 ? 'destructive' : 'outline'}>
					{summary.security.suspiciousAttempts} auffällig
				</Badge>
			</Card.Action>
		</Card.Header>
		<Card.Footer class="text-muted-foreground text-sm">{summary.users.locked} gesperrte Benutzer</Card.Footer>
	</Card.Root>
</div>
