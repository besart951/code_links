<script lang="ts">
	import { Button } from '@codelinks/ui/button';
	import * as Card from '@codelinks/ui/card';
	import { Input } from '@codelinks/ui/input';
	import { Label } from '@codelinks/ui/label';
	import { page } from '$app/state';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();
	const returnTo = $derived(page.url.searchParams.get('returnTo') ?? '/admin');
	const loginError = $derived(page.url.searchParams.get('error'));
</script>

<main class="grid min-h-screen place-items-center bg-muted/30 p-4">
	<Card.Root class="w-full max-w-sm">
		<Card.Header>
			<Card.Title>Superadmin</Card.Title>
			<Card.Description>Platform session required</Card.Description>
		</Card.Header>
		<Card.Content>
			{#if loginError === 'platform_unavailable' || form?.error === 'platform_unavailable'}
				<div class="rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive">
					Platform API is not reachable. Start the Go Platform API or set PLATFORM_API_URL to the running API.
				</div>
			{:else if form?.error}
				<div class="rounded-md border border-destructive/30 bg-destructive/10 px-3 py-2 text-sm text-destructive">
					Login failed. Check the Superadmin credentials.
				</div>
			{/if}
			<form method="POST" class="mt-4 space-y-4">
				<input type="hidden" name="returnTo" value={returnTo} />
				<div class="space-y-2">
					<Label for="email">Email</Label>
					<Input id="email" name="email" type="email" autocomplete="email" value={form?.email ?? ''} required />
				</div>
				<div class="space-y-2">
					<Label for="password">Password</Label>
					<Input id="password" name="password" type="password" autocomplete="current-password" required />
				</div>
				<Button type="submit" class="w-full">Sign in</Button>
			</form>
		</Card.Content>
	</Card.Root>
</main>
