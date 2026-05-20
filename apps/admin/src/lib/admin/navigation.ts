import {
	Activity,
	Bell,
	Building2,
	LayoutDashboard,
	Package,
	Search,
	Settings,
	Shield,
	Users
} from '@lucide/svelte';
import type { Component } from 'svelte';

export type AdminNavItem = {
	href: string;
	label: string;
	icon: Component;
};

export const adminNavItems: AdminNavItem[] = [
	{ href: '/admin', label: 'Dashboard', icon: LayoutDashboard },
	{ href: '/admin/search', label: 'Search', icon: Search },
	{ href: '/admin/tenants', label: 'Tenants', icon: Building2 },
	{ href: '/admin/users', label: 'Users', icon: Users },
	{ href: '/admin/products', label: 'Products', icon: Package },
	{ href: '/admin/notifications', label: 'Notifications', icon: Bell },
	{ href: '/admin/audit', label: 'Audit', icon: Activity },
	{ href: '/admin/security', label: 'Security', icon: Shield },
	{ href: '/admin/settings', label: 'Settings', icon: Settings }
];
