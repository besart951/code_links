<script lang="ts">
	import { Button } from '@codelinks/ui-library/components/ui/button';
	import * as Card from '@codelinks/ui-library/components/ui/card';
	import * as Field from '@codelinks/ui-library/components/ui/field';
	import { Input } from '@codelinks/ui-library/components/ui/input';

	let {
		form
	}: {
		form?: {
			email?: string;
			error?: string;
			success?: boolean;
			message?: string;
			debugResetUrl?: string;
		};
	} = $props();
</script>

<Card.Root class="w-full max-w-sm">
	<Card.Header>
		<Card.Title>Passwort vergessen</Card.Title>
		<Card.Description>Reset-Link anfordern</Card.Description>
	</Card.Header>
	<Card.Content>
		<form method="POST">
			<Field.Group>
				<Field.Field>
					<Field.Label for="email">E-Mail</Field.Label>
					<Input id="email" name="email" type="email" autocomplete="email" value={form?.email ?? ''} required />
				</Field.Field>
				<Field.Field>
					<Button type="submit" class="w-full">Reset-Link senden</Button>
					<Field.Description class="text-center">
						<a href="/login">Zurück zum Login</a>
					</Field.Description>
				</Field.Field>
			</Field.Group>
		</form>

		{#if form?.success}
			<div class="mt-4 rounded-md border border-border bg-muted/40 p-3 text-sm">
				<p>{form.message}</p>
				{#if form.debugResetUrl}
					<a class="mt-2 block underline" href={form.debugResetUrl}>Reset-Link öffnen</a>
				{/if}
			</div>
		{/if}

		{#if form?.error}
			<p class="mt-4 text-sm text-destructive">{form.error}</p>
		{/if}
	</Card.Content>
</Card.Root>
