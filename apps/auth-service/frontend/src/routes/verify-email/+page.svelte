<script lang="ts">
	import AuthShell from '$lib/components/auth-shell.svelte';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>E-Mail bestätigen | CodeLinks</title>
</svelte:head>

<AuthShell title="CodeLinks Konto" subtitle="E-Mail-Bestätigung für neue Accounts.">
	<Card.Root class="w-full max-w-sm">
		<Card.Header>
			<Card.Title>{data.verified ? 'E-Mail bestätigt' : 'Bestätigung fehlgeschlagen'}</Card.Title>
			<Card.Description>
				{data.verified ? data.email : data.error}
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<Button class="w-full" href="/login">Zum Login</Button>
		</Card.Content>
	</Card.Root>
</AuthShell>
