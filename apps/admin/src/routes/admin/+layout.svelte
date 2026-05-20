<script lang="ts">
	import SearchCommand from '$lib/admin/components/SearchCommand.svelte';
	import { adminNavItems } from '$lib/admin/navigation.js';
	import * as Avatar from '@codelinks/ui/avatar';
	import { Button } from '@codelinks/ui/button';
	import { Separator } from '@codelinks/ui/separator';
	import * as Sidebar from '@codelinks/ui/sidebar';
	import { LogOut } from '@lucide/svelte';
	import { page } from '$app/state';
	import type { Snippet } from 'svelte';
	import type { LayoutData } from './$types';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();
</script>

<Sidebar.Provider>
	<Sidebar.Root class="border-r">
		<Sidebar.Header class="gap-3 px-3 py-4">
			<a href="/admin" class="flex min-w-0 items-center gap-2">
				<div class="grid size-8 shrink-0 place-items-center rounded-md bg-primary text-sm font-semibold text-primary-foreground">
					CL
				</div>
				<div class="min-w-0">
					<p class="truncate text-sm font-semibold">CodeLinks</p>
					<p class="truncate text-xs text-muted-foreground">Superadmin</p>
				</div>
			</a>
		</Sidebar.Header>
		<Sidebar.Content>
			<Sidebar.Group>
				<Sidebar.GroupContent>
					<Sidebar.Menu>
						{#each adminNavItems as item (item.href)}
							<Sidebar.MenuItem>
								<Sidebar.MenuButton isActive={page.url.pathname === item.href}>
									{#snippet child({ props })}
										<a href={item.href} {...props}>
											<item.icon />
											<span>{item.label}</span>
										</a>
									{/snippet}
								</Sidebar.MenuButton>
							</Sidebar.MenuItem>
						{/each}
					</Sidebar.Menu>
				</Sidebar.GroupContent>
			</Sidebar.Group>
		</Sidebar.Content>
		<Sidebar.Footer class="gap-3 p-3">
			<Separator />
			<div class="flex min-w-0 items-center gap-2">
				<Avatar.Root class="size-8">
					<Avatar.Fallback>{data.admin?.user.display_name?.slice(0, 2).toUpperCase() ?? 'SA'}</Avatar.Fallback>
				</Avatar.Root>
				<div class="min-w-0 flex-1">
					<p class="truncate text-sm font-medium">{data.admin?.user.display_name ?? 'Superadmin'}</p>
					<p class="truncate text-xs text-muted-foreground">{data.admin?.user.email}</p>
				</div>
				<Button variant="ghost" size="icon-sm" aria-label="Logout">
					<LogOut />
				</Button>
			</div>
		</Sidebar.Footer>
	</Sidebar.Root>

	<Sidebar.Inset>
		<header class="sticky top-0 z-20 flex h-14 items-center gap-3 border-b bg-background/95 px-4 backdrop-blur">
			<Sidebar.Trigger />
			<SearchCommand />
		</header>
		<main class="min-h-[calc(100vh-3.5rem)] bg-muted/25 px-4 py-4 sm:px-6">
			{@render children()}
		</main>
	</Sidebar.Inset>
</Sidebar.Provider>
