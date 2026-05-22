<script lang="ts">
	import { Button } from '$lib/components/ui/button/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import * as Field from '$lib/components/ui/field/index.js';
	import { Input } from '$lib/components/ui/input/index.js';

	let {
		form,
		redirectTo = '/'
	}: {
		form?: { email?: string; error?: string };
		redirectTo?: string;
	} = $props();
</script>

<Card.Root class="w-full max-w-sm">
	<Card.Header>
		<Card.Title class="text-2xl">Login</Card.Title>
		<Card.Description>Mit deinem CodeLinks Konto anmelden</Card.Description>
	</Card.Header>
	<Card.Content>
		<form method="POST" data-sveltekit-reload>
			<input type="hidden" name="redirectTo" value={redirectTo} />
			<Field.Group>
				<Field.Field>
					<Field.Label for="email">E-Mail</Field.Label>
					<Input
						id="email"
						name="email"
						type="email"
						autocomplete="email"
						value={form?.email ?? 'demo@codelinks.dev'}
						required
					/>
				</Field.Field>
				<Field.Field>
					<div class="flex items-center">
						<Field.Label for="password">Passwort</Field.Label>
						<a href="/forgot-password" class="ms-auto inline-block text-sm underline">
							Passwort vergessen?
						</a>
					</div>
					<Input
						id="password"
						name="password"
						type="password"
						autocomplete="current-password"
						value="password"
						required
					/>
				</Field.Field>
				<Field.Field orientation="horizontal">
					<Checkbox id="remember" name="remember" />
					<Field.Label for="remember" class="font-normal">Angemeldet bleiben</Field.Label>
				</Field.Field>
				<Field.Field>
					<Button type="submit" class="w-full">Login</Button>
					<Field.Description class="text-center">
						Noch kein Account? <a href="/signup">Account erstellen</a>
					</Field.Description>
					{#if form?.error}
						<Field.Error>{form.error}</Field.Error>
					{/if}
				</Field.Field>
			</Field.Group>
		</form>
	</Card.Content>
</Card.Root>
