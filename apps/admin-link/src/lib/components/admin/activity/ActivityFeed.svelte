<script lang="ts">
	import { onMount, tick } from 'svelte';
	import AlertTriangleIcon from '@tabler/icons-svelte/icons/alert-triangle';
	import BellIcon from '@tabler/icons-svelte/icons/bell';
	import CheckIcon from '@tabler/icons-svelte/icons/check';
	import ClipboardListIcon from '@tabler/icons-svelte/icons/clipboard-list';
	import CodeIcon from '@tabler/icons-svelte/icons/code';
	import LoginIcon from '@tabler/icons-svelte/icons/login';
	import MoonIcon from '@tabler/icons-svelte/icons/moon';
	import ShieldLockIcon from '@tabler/icons-svelte/icons/shield-lock';
	import SunIcon from '@tabler/icons-svelte/icons/sun';
	import { Badge } from '@codelinks/ui-library/components/ui/badge';
	import { Button } from '@codelinks/ui-library/components/ui/button';
	import type { ActivityEvent, ActivityEventSource, ActivityEventTone } from '$lib/domain/activity/types';

	type SourceFilter = 'all' | ActivityEventSource;

	let { events }: { events: ActivityEvent[] } = $props();

	let selectedSource = $state<SourceFilter>('all');
	let darkTheme = $state(false);
	let viewport: HTMLDivElement;

	const filteredEvents = $derived(
		selectedSource === 'all' ? events : events.filter((event) => event.source === selectedSource)
	);

	onMount(() => {
		tick().then(() => viewport?.scrollTo({ top: viewport.scrollHeight }));
	});

	function sourceLabel(source: SourceFilter) {
		const labels: Record<SourceFilter, string> = {
			all: 'Alle',
			auth: 'Auth',
			security: 'Security',
			notification: 'Mail',
			audit: 'Audit',
			runtime: 'Runtime'
		};
		return labels[source];
	}

	function toneClasses(tone: ActivityEventTone) {
		const classes: Record<ActivityEventTone, string> = {
			neutral: 'border-border bg-card text-card-foreground',
			success: 'border-emerald-300 bg-emerald-50 text-emerald-950 dark:border-emerald-900 dark:bg-emerald-950/40 dark:text-emerald-100',
			warning: 'border-amber-300 bg-amber-50 text-amber-950 dark:border-amber-900 dark:bg-amber-950/40 dark:text-amber-100',
			danger: 'border-red-300 bg-red-50 text-red-950 dark:border-red-900 dark:bg-red-950/40 dark:text-red-100',
			info: 'border-sky-300 bg-sky-50 text-sky-950 dark:border-sky-900 dark:bg-sky-950/40 dark:text-sky-100'
		};
		return classes[tone];
	}

	function toneIcon(tone: ActivityEventTone) {
		return tone === 'success' ? CheckIcon : tone === 'danger' || tone === 'warning' ? AlertTriangleIcon : BellIcon;
	}

	function sourceIcon(source: ActivityEventSource) {
		const icons = {
			auth: LoginIcon,
			security: ShieldLockIcon,
			notification: BellIcon,
			audit: ClipboardListIcon,
			runtime: CodeIcon
		};
		return icons[source];
	}

	function formatDate(value: string) {
		return new Intl.DateTimeFormat('de-DE', {
			day: '2-digit',
			month: '2-digit',
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit'
		}).format(new Date(value));
	}
</script>

<section class={['rounded-lg border bg-background text-foreground', darkTheme && 'dark']}>
	<div class="flex flex-col gap-3 border-b bg-muted/40 p-3 md:flex-row md:items-center md:justify-between">
		<div class="flex flex-wrap gap-1">
			{#each ['all', 'runtime', 'auth', 'security', 'notification', 'audit'] as source (source)}
				<Button
					type="button"
					size="sm"
					variant={selectedSource === source ? 'default' : 'ghost'}
					onclick={() => (selectedSource = source as SourceFilter)}
				>
					{sourceLabel(source as SourceFilter)}
				</Button>
			{/each}
		</div>
		<Button type="button" size="icon" variant="outline" aria-label="Theme wechseln" onclick={() => (darkTheme = !darkTheme)}>
			{#if darkTheme}
				<SunIcon class="size-4" />
			{:else}
				<MoonIcon class="size-4" />
			{/if}
		</Button>
	</div>

	<div bind:this={viewport} class="max-h-[calc(100vh-16rem)] min-h-[32rem] overflow-y-auto p-4">
		<div class="mx-auto flex max-w-5xl flex-col gap-4">
			{#each filteredEvents as event (event.id)}
				{@const ToneIcon = toneIcon(event.tone)}
				{@const SourceIcon = sourceIcon(event.source)}
				<article class={['rounded-lg border p-4 shadow-sm', toneClasses(event.tone)]}>
					<div class="flex items-start gap-3">
						<div class="mt-0.5 flex size-9 shrink-0 items-center justify-center rounded-md bg-background/80">
							<ToneIcon class="size-4" />
						</div>
						<div class="min-w-0 flex-1 space-y-3">
							<div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
								<div class="min-w-0">
									<div class="flex flex-wrap items-center gap-2">
										<SourceIcon class="size-4 shrink-0 opacity-70" />
										<h3 class="truncate text-sm font-semibold">{event.title}</h3>
										<Badge variant="outline">{sourceLabel(event.source)}</Badge>
									</div>
									<p class="mt-1 break-words text-sm leading-6">{event.message}</p>
								</div>
								<time class="shrink-0 font-mono text-xs text-muted-foreground" datetime={event.occurredAt}>
									{formatDate(event.occurredAt)}
								</time>
							</div>

							{#if event.details.length > 0}
								<dl class="grid gap-2 text-xs sm:grid-cols-2 lg:grid-cols-4">
									{#each event.details as detail (detail.label)}
										<div class="rounded-md bg-background/70 px-2.5 py-2">
											<dt class="text-muted-foreground">{detail.label}</dt>
											<dd class="mt-1 break-words font-mono">{detail.value}</dd>
										</div>
									{/each}
								</dl>
							{/if}
						</div>
					</div>
				</article>
			{:else}
				<div class="rounded-lg border border-dashed p-8 text-center text-sm text-muted-foreground">Keine Logs gefunden.</div>
			{/each}
		</div>
	</div>
</section>
