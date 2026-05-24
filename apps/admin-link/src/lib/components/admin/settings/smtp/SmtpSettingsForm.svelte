<script lang="ts">
	import { Button } from '@codelinks/ui-library/components/ui/button';
	import * as Card from '@codelinks/ui-library/components/ui/card';
	import { Checkbox } from '@codelinks/ui-library/components/ui/checkbox';
	import { Input } from '@codelinks/ui-library/components/ui/input';
	import { Label } from '@codelinks/ui-library/components/ui/label';
	import type { SmtpSettings } from '$lib/domain/smtp/types';

	let {
		settings,
		canUpdate,
		form
	}: {
		settings: SmtpSettings;
		canUpdate: boolean;
		form?: { ok?: boolean; error?: boolean; message?: string };
	} = $props();
</script>

<div class="grid gap-4 xl:grid-cols-[1fr_22rem]">
	<Card.Root>
		<Card.Header>
			<Card.Title>SMTP-Konfiguration</Card.Title>
			<Card.Description>Aktualisiert am {new Date(settings.updatedAt).toLocaleString('de-DE')}</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" action="?/save" class="grid gap-4 md:grid-cols-2">
				<div class="grid gap-2">
					<Label for="host">SMTP Host</Label>
					<Input id="host" name="host" value={settings.host} disabled={!canUpdate} required />
				</div>
				<div class="grid gap-2">
					<Label for="port">SMTP Port</Label>
					<Input id="port" name="port" type="number" min="1" max="65535" value={settings.port} disabled={!canUpdate} required />
				</div>
				<div class="grid gap-2">
					<Label for="username">SMTP Username</Label>
					<Input id="username" name="username" value={settings.username} disabled={!canUpdate} />
				</div>
				<div class="grid gap-2">
					<Label for="password">SMTP Password</Label>
					<Input
						id="password"
						name="password"
						type="password"
						placeholder={settings.hasPassword ? 'Vorhandenes Passwort bleibt erhalten' : ''}
						disabled={!canUpdate}
					/>
				</div>
				<div class="grid gap-2">
					<Label for="encryption">Encryption Type</Label>
					<select
						id="encryption"
						name="encryption"
						class="border-input bg-background h-9 rounded-md border px-3 text-sm"
						disabled={!canUpdate}
					>
						<option selected={settings.encryption === 'none'} value="none">none</option>
						<option selected={settings.encryption === 'ssl'} value="ssl">SSL</option>
						<option selected={settings.encryption === 'tls'} value="tls">TLS</option>
						<option selected={settings.encryption === 'starttls'} value="starttls">STARTTLS</option>
					</select>
				</div>
				<div class="grid gap-2">
					<Label for="fromEmail">From E-Mail</Label>
					<Input id="fromEmail" name="fromEmail" type="email" value={settings.fromEmail} disabled={!canUpdate} required />
				</div>
				<div class="grid gap-2">
					<Label for="fromName">From Name</Label>
					<Input id="fromName" name="fromName" value={settings.fromName} disabled={!canUpdate} required />
				</div>
				<div class="grid gap-2">
					<Label for="replyToEmail">Reply-To E-Mail</Label>
					<Input id="replyToEmail" name="replyToEmail" type="email" value={settings.replyToEmail} disabled={!canUpdate} required />
				</div>
				<div class="flex items-center gap-2 md:col-span-2">
					<Checkbox id="active" name="active" checked={settings.active} disabled={!canUpdate} />
					<Label for="active" class="font-normal">SMTP aktiv</Label>
				</div>
				<div class="flex items-center gap-3 md:col-span-2">
					<Button type="submit" disabled={!canUpdate}>Speichern</Button>
					{#if form?.message}
						<p class={form.error ? 'text-sm text-destructive' : 'text-sm text-muted-foreground'}>
							{form.message}
						</p>
					{/if}
				</div>
			</form>
		</Card.Content>
	</Card.Root>

	<Card.Root>
		<Card.Header>
			<Card.Title>Test-E-Mail</Card.Title>
			<Card.Description>Versand mit aktueller Konfiguration prüfen</Card.Description>
		</Card.Header>
		<Card.Content>
			<form method="POST" action="?/testEmail" class="grid gap-3">
				<Label for="recipient">Empfänger</Label>
				<Input id="recipient" name="recipient" type="email" value="admin@example.com" disabled={!canUpdate} required />
				<Button type="submit" variant="secondary" disabled={!canUpdate}>Test-E-Mail senden</Button>
			</form>
		</Card.Content>
	</Card.Root>
</div>
