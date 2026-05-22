<script lang="ts">
	let {
		title,
		subtitle,
		children
	}: {
		title: string;
		subtitle: string;
		children: import('svelte').Snippet;
	} = $props();
</script>

<main class="grid min-h-screen bg-muted/30 px-4 py-10 lg:grid-cols-[0.95fr_1.05fr] lg:px-8">
	<section class="hidden min-h-[calc(100vh-5rem)] flex-col justify-between rounded-lg bg-primary p-10 text-primary-foreground lg:flex">
		<a class="text-lg font-semibold tracking-normal" href="/">CodeLinks</a>
		<div class="max-w-lg">
			<h1 class="text-4xl font-semibold tracking-normal">{title}</h1>
			<p class="mt-4 max-w-md text-base leading-7 text-primary-foreground/80">{subtitle}</p>
		</div>
		<p class="text-sm text-primary-foreground/70">auth.codelinks.localhost</p>
	</section>

	<section class="flex min-h-[calc(100vh-5rem)] items-center justify-center">
		{@render children()}
	</section>
</main>
