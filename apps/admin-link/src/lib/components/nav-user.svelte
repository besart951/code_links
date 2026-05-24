<script lang="ts">
  import DotsVerticalIcon from "@tabler/icons-svelte/icons/dots-vertical";
  import LogoutIcon from "@tabler/icons-svelte/icons/logout";
  import MoonIcon from "@tabler/icons-svelte/icons/moon";
  import ShieldIcon from "@tabler/icons-svelte/icons/shield";
  import SunIcon from "@tabler/icons-svelte/icons/sun";
  import { browser } from "$app/environment";
  import { onMount } from "svelte";
  import * as Avatar from "@codelinks/ui-library/components/ui/avatar";
  import * as DropdownMenu from "@codelinks/ui-library/components/ui/dropdown-menu";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import type { AdminActor } from "$lib/domain/admin-access/types";

  let { user }: { user: AdminActor } = $props();

  const sidebar = Sidebar.useSidebar();
  let isDark = $state(false);
  const initials = $derived(
    user.name
      .split(" ")
      .map((part) => part[0])
      .join("")
      .slice(0, 2)
      .toUpperCase(),
  );

  onMount(() => {
    if (!browser) return;

    const savedTheme = localStorage.getItem("theme");
    if (savedTheme === "dark" || savedTheme === "light") {
      isDark = savedTheme === "dark";
      document.documentElement.classList.toggle("dark", isDark);
      return;
    }

    isDark = document.documentElement.classList.contains("dark");
  });

  function toggleTheme() {
    if (!browser) return;

    isDark = !isDark;
    document.documentElement.classList.toggle("dark", isDark);
    localStorage.setItem("theme", isDark ? "dark" : "light");
  }
</script>

<Sidebar.Menu>
  <Sidebar.MenuItem>
    <DropdownMenu.Root>
      <DropdownMenu.Trigger>
        {#snippet child({ props })}
          <Sidebar.MenuButton
            {...props}
            size="lg"
            class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <Avatar.Root class="size-8 rounded-lg">
              <Avatar.Fallback class="rounded-lg">{initials}</Avatar.Fallback>
            </Avatar.Root>
            <div class="grid flex-1 text-start text-sm leading-tight">
              <span class="truncate font-medium">{user.name}</span>
              <span class="text-muted-foreground truncate text-xs"
                >{user.email}</span
              >
            </div>
            <DotsVerticalIcon class="ms-auto size-4" />
          </Sidebar.MenuButton>
        {/snippet}
      </DropdownMenu.Trigger>
      <DropdownMenu.Content
        class="w-(--bits-dropdown-menu-anchor-width) min-w-56 rounded-lg"
        side={sidebar.isMobile ? "bottom" : "right"}
        align="end"
        sideOffset={4}
      >
        <DropdownMenu.Label class="p-0 font-normal">
          <div class="flex items-center gap-2 px-1 py-1.5 text-start text-sm">
            <Avatar.Root class="size-8 rounded-lg">
              <Avatar.Fallback class="rounded-lg">{initials}</Avatar.Fallback>
            </Avatar.Root>
            <div class="grid flex-1 text-start text-sm leading-tight">
              <span class="truncate font-medium">{user.name}</span>
              <span class="text-muted-foreground truncate text-xs"
                >{user.roles.join(", ")}</span
              >
            </div>
          </div>
        </DropdownMenu.Label>
        <DropdownMenu.Separator />
        <DropdownMenu.Item>
          <ShieldIcon />
          {user.permissions.length} Berechtigungen
        </DropdownMenu.Item>
        <DropdownMenu.Separator />
        <DropdownMenu.Item onclick={toggleTheme}>
          {#if isDark}
            <SunIcon />
            Helles Theme
          {:else}
            <MoonIcon />
            Dunkles Theme
          {/if}
        </DropdownMenu.Item>
        <DropdownMenu.Separator />
        <DropdownMenu.Item disabled>
          <LogoutIcon />
          Logout via Auth-Service folgt
        </DropdownMenu.Item>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  </Sidebar.MenuItem>
</Sidebar.Menu>
