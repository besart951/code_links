<script lang="ts">
	import { Button } from '@codelinks/ui-library/components/ui/button';
	import * as Card from '@codelinks/ui-library/components/ui/card';
	import * as Field from '@codelinks/ui-library/components/ui/field';
	import { Input } from '@codelinks/ui-library/components/ui/input';

	let {
		form,
		token
	}: {
		form?: { token?: string; error?: string; success?: boolean; message?: string };
		token?: string;
	} = $props();
</script>

<Card.Root class="w-full max-w-sm">
	<Card.Header>
		<Card.Title>Passwort ändern</Card.Title>
		<Card.Description>Neues Passwort festlegen</Card.Description>
	</Card.Header>
	<Card.Content>
		{#if form?.success}
			<div class="rounded-md border border-border bg-muted/40 p-4 text-sm">
				<p>{form.message}</p>
				<a class="mt-3 block underline" href="/login">Zum Login</a>
			</div>
		{:else}
			<form method="POST">
				<input type="hidden" name="token" value={form?.token ?? token ?? ''} />
				<Field.Group>
					<Field.Field>
						<Field.Label for="password">Neues Passwort</Field.Label>
						<Input id="password" name="password" type="password" autocomplete="new-password" required />
						<Field.Description>Mindestens 10 Zeichen.</Field.Description>
					</Field.Field>
					<Field.Field>
						<Field.Label for="confirm-password">Neues Passwort bestätigen</Field.Label>
						<Input id="confirm-password" name="confirmPassword" type="password" autocomplete="new-password" required />
					</Field.Field>
					<Field.Field>
						<Button type="submit" class="w-full">Passwort ändern</Button>
						<Field.Description class="text-center">
							<a href="/login">Zurück zum Login</a>
						</Field.Description>
						{#if form?.error}
							<Field.Error>{form.error}</Field.Error>
						{/if}
					</Field.Field>
				</Field.Group>
			</form>
		{/if}
	</Card.Content>
</Card.Root>
