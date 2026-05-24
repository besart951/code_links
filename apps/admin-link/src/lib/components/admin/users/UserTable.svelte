<script lang="ts">
	import { Button } from '@codelinks/ui-library/components/ui/button';
	import { Badge } from '@codelinks/ui-library/components/ui/badge';
	import { Input } from '@codelinks/ui-library/components/ui/input';
	import * as Table from '@codelinks/ui-library/components/ui/table';
	import { UserTableState } from '$lib/application/users/UserTableState.svelte';
	import { hasPermission } from '$lib/domain/admin-access/permissions';
	import type { AdminActor } from '$lib/domain/admin-access/types';
	import type { UserListItem, UserListQuery } from '$lib/domain/users/types';
	import UserStatusBadge from './UserStatusBadge.svelte';

	let {
		users,
		total,
		query,
		admin,
		maskIp = false
	}: {
		users: UserListItem[];
		total: number;
		query: UserListQuery;
		admin: AdminActor;
		maskIp?: boolean;
	} = $props();

	const state = new UserTableState();

	$effect(() => {
		state.query = query.query ?? '';
		state.role = query.role ?? 'all';
		state.status = query.status ?? 'all';
		state.page = query.page;
		state.pageSize = query.pageSize;
		state.sort = query.sort;
	});

	const canLockUser = $derived(hasPermission(admin, 'admin.users.lock'));
	const canChangeRole = $derived(hasPermission(admin, 'admin.users.change_role'));

	function formatIp(ipAddress: string | null) {
		if (!ipAddress) return '-';
		return maskIp ? ipAddress.replace(/\.\d+$/, '.xxx') : ipAddress;
	}
</script>

<div class="space-y-4">
	<form method="GET" class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
		<Input name="query" placeholder="Name oder E-Mail suchen" bind:value={state.query} class="md:max-w-80" />

		<div class="flex flex-wrap gap-2">
			<select name="role" bind:value={state.role} class="border-input bg-background h-9 rounded-md border px-3 text-sm">
				<option value="all">Alle Rollen</option>
				<option value="admin">Admin</option>
				<option value="support">Support</option>
				<option value="auditor">Auditor</option>
				<option value="user">User</option>
			</select>
			<select name="status" bind:value={state.status} class="border-input bg-background h-9 rounded-md border px-3 text-sm">
				<option value="all">Alle Status</option>
				<option value="active">Aktiv</option>
				<option value="disabled">Deaktiviert</option>
				<option value="locked">Gesperrt</option>
			</select>
			<Button type="submit" variant="secondary">Filtern</Button>
		</div>
	</form>

	<div class="rounded-lg border">
		<Table.Root>
			<Table.Header>
				<Table.Row>
					<Table.Head>Name</Table.Head>
					<Table.Head>E-Mail</Table.Head>
					<Table.Head>Rolle</Table.Head>
					<Table.Head>Status</Table.Head>
					<Table.Head>E-Mail</Table.Head>
					<Table.Head>Erstellt am</Table.Head>
					<Table.Head>Letzter Login</Table.Head>
					<Table.Head>Logins</Table.Head>
					<Table.Head>Letzte IP</Table.Head>
					<Table.Head>Land</Table.Head>
					<Table.Head class="text-right">Aktionen</Table.Head>
				</Table.Row>
			</Table.Header>
			<Table.Body>
				{#each users as user (user.id)}
					<Table.Row>
						<Table.Cell class="font-medium">{user.name}</Table.Cell>
						<Table.Cell>{user.email}</Table.Cell>
						<Table.Cell>{user.primaryRole}</Table.Cell>
						<Table.Cell><UserStatusBadge status={user.status} /></Table.Cell>
						<Table.Cell>
							<Badge variant={user.emailVerified ? 'outline' : 'destructive'}>
								{user.emailVerified ? 'Bestätigt' : 'Offen'}
							</Badge>
						</Table.Cell>
						<Table.Cell>{new Date(user.createdAt).toLocaleDateString('de-DE')}</Table.Cell>
						<Table.Cell>{user.lastLoginAt ? new Date(user.lastLoginAt).toLocaleString('de-DE') : 'Nie'}</Table.Cell>
						<Table.Cell>{user.successfulLoginCount} / {user.failedLoginCount}</Table.Cell>
						<Table.Cell class="font-mono text-xs">{formatIp(user.lastKnownIpAddress)}</Table.Cell>
						<Table.Cell>{user.lastLoginCountryCode ?? '-'}</Table.Cell>
						<Table.Cell class="space-x-2 text-right">
							<Button variant="ghost" size="sm" href={`/admin/users/${user.id}`}>Öffnen</Button>
							{#if canChangeRole}
								<form method="POST" action="?/setRole" class="inline-flex items-center gap-2">
									<input type="hidden" name="userId" value={user.id} />
									<select
										name="role"
										class="border-input bg-background h-8 rounded-md border px-2 text-xs"
										aria-label="Rolle ändern"
									>
										<option selected={user.primaryRole === 'admin'} value="admin">Admin</option>
										<option selected={user.primaryRole === 'support'} value="support">Support</option>
										<option selected={user.primaryRole === 'auditor'} value="auditor">Auditor</option>
										<option selected={user.primaryRole === 'user'} value="user">User</option>
									</select>
									<Button variant="outline" size="sm" type="submit">Rolle</Button>
								</form>
							{/if}
							{#if canLockUser}
								<form method="POST" action="?/setStatus" class="inline">
									<input type="hidden" name="userId" value={user.id} />
									<input type="hidden" name="status" value={user.status === 'locked' ? 'active' : 'locked'} />
									<Button variant="outline" size="sm" type="submit">
										{user.status === 'locked' ? 'Entsperren' : 'Sperren'}
									</Button>
								</form>
							{/if}
						</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>

	<p class="text-muted-foreground text-sm">{total} Benutzer</p>
</div>
