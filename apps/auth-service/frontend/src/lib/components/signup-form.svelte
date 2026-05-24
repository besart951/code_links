<script lang="ts">
	import { Button } from '@codelinks/ui-library/components/ui/button';
	import * as Card from '@codelinks/ui-library/components/ui/card';
	import { Checkbox } from '@codelinks/ui-library/components/ui/checkbox';
	import * as Field from '@codelinks/ui-library/components/ui/field';
	import { Input } from '@codelinks/ui-library/components/ui/input';

	let {
		form
	}: {
		form?: {
			name?: string;
			email?: string;
			error?: string;
			success?: boolean;
			message?: string;
			verificationUrl?: string;
		};
	} = $props();
</script>

<Card.Root class="w-full max-w-sm">
	<Card.Header>
		<Card.Title>Account erstellen</Card.Title>
		<Card.Description>Registrierung für CodeLinks</Card.Description>
	</Card.Header>
	<Card.Content>
		{#if form?.success}
			<div class="rounded-md border border-border bg-muted/40 p-4 text-sm">
				<p>{form.message}</p>
				{#if form.verificationUrl}
					<a class="mt-3 block underline" href={form.verificationUrl}>E-Mail jetzt bestätigen</a>
				{/if}
				<a class="mt-3 block underline" href="/login">Zum Login</a>
			</div>
		{:else}
			<form method="POST">
				<Field.Group>
					<Field.Field>
						<Field.Label for="name">Name</Field.Label>
						<Input id="name" name="name" type="text" autocomplete="name" value={form?.name ?? ''} required />
					</Field.Field>
					<Field.Field>
						<Field.Label for="email">E-Mail</Field.Label>
						<Input id="email" name="email" type="email" autocomplete="email" value={form?.email ?? ''} required />
					</Field.Field>
					<Field.Field>
						<Field.Label for="password">Passwort</Field.Label>
						<Input id="password" name="password" type="password" autocomplete="new-password" required />
						<Field.Description>Mindestens 10 Zeichen.</Field.Description>
					</Field.Field>
					<Field.Field>
						<Field.Label for="confirm-password">Passwort bestätigen</Field.Label>
						<Input id="confirm-password" name="confirmPassword" type="password" autocomplete="new-password" required />
					</Field.Field>
					<Field.Field orientation="horizontal">
						<Checkbox id="accepted-terms" name="acceptedTerms" required />
						<Field.Label for="accepted-terms" class="font-normal">
							Ich akzeptiere die Nutzungsbedingungen
						</Field.Label>
					</Field.Field>
					<Field.Field>
						<Button type="submit" class="w-full">Account erstellen</Button>
						<Field.Description class="text-center">
							Bereits registriert? <a href="/login">Einloggen</a>
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
