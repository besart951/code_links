<script lang="ts">
	import AlertTriangleIcon from '@tabler/icons-svelte/icons/alert-triangle';
	import ChartBarIcon from '@tabler/icons-svelte/icons/chart-bar';
	import DashboardIcon from '@tabler/icons-svelte/icons/dashboard';
	import HistoryIcon from '@tabler/icons-svelte/icons/history';
	import LockIcon from '@tabler/icons-svelte/icons/lock';
	import MailIcon from '@tabler/icons-svelte/icons/mail';
	import SettingsIcon from '@tabler/icons-svelte/icons/settings';
	import ShieldLockIcon from '@tabler/icons-svelte/icons/shield-lock';
	import UsersIcon from '@tabler/icons-svelte/icons/users';
	import NavUser from './nav-user.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import { hasPermission } from '$lib/domain/admin-access/permissions';
	import type { AdminActor, AdminPermission } from '$lib/domain/admin-access/types';
	import type { ComponentProps } from 'svelte';
	import type { Icon } from '@tabler/icons-svelte';

	type NavItem = {
		title: string;
		url: string;
		icon: Icon;
		permission: AdminPermission;
	};

	let {
		admin,
		pathname,
		...restProps
	}: { admin: AdminActor; pathname: string } & ComponentProps<typeof Sidebar.Root> = $props();

	const items: NavItem[] = [
		{ title: 'Übersicht', url: '/admin', icon: DashboardIcon, permission: 'admin.dashboard.read' },
		{ title: 'Benutzer', url: '/admin/users', icon: UsersIcon, permission: 'admin.users.read' },
		{ title: 'Login-Historie', url: '/admin/login-history', icon: HistoryIcon, permission: 'admin.auth_logs.read' },
		{ title: 'Fehlgeschlagen', url: '/admin/failed-logins', icon: AlertTriangleIcon, permission: 'admin.auth_logs.read' },
		{ title: 'Security', url: '/admin/security', icon: ShieldLockIcon, permission: 'admin.security_events.read' },
		{ title: 'Benachrichtigungen', url: '/admin/notifications', icon: MailIcon, permission: 'admin.notifications.read' },
		{ title: 'SMTP', url: '/admin/settings/smtp', icon: SettingsIcon, permission: 'admin.smtp_settings.read' },
		{ title: 'Audit Log', url: '/admin/audit', icon: LockIcon, permission: 'admin.audit_entries.read' }
	];

	const visibleItems = $derived(items.filter((item) => hasPermission(admin, item.permission)));

	function isActive(url: string) {
		return url === '/admin' ? pathname === url : pathname.startsWith(url);
	}
</script>

<Sidebar.Root collapsible="offcanvas" {...restProps}>
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton class="data-[slot=sidebar-menu-button]:!p-1.5">
					{#snippet child({ props })}
						<a href="/admin" {...props}>
							<ChartBarIcon class="!size-5" />
							<span class="text-base font-semibold">CodeLinks Admin</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>
	<Sidebar.Content>
		<Sidebar.Group>
			<Sidebar.GroupLabel>Verwaltung</Sidebar.GroupLabel>
			<Sidebar.GroupContent>
				<Sidebar.Menu>
					{#each visibleItems as item (item.url)}
						<Sidebar.MenuItem>
							<Sidebar.MenuButton tooltipContent={item.title} isActive={isActive(item.url)}>
								{#snippet child({ props })}
									<a href={item.url} {...props}>
										<item.icon />
										<span>{item.title}</span>
									</a>
								{/snippet}
							</Sidebar.MenuButton>
						</Sidebar.MenuItem>
					{/each}
				</Sidebar.Menu>
			</Sidebar.GroupContent>
		</Sidebar.Group>
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser user={admin} />
	</Sidebar.Footer>
</Sidebar.Root>
