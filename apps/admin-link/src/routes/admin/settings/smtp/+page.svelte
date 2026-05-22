<script lang="ts">
	import SmtpSettingsForm from '$lib/components/admin/settings/smtp/SmtpSettingsForm.svelte';
	import type { ActionData, PageData } from './$types';

	let { data, form }: { data: PageData; form: ActionData } = $props();
	const actionForm = $derived(form ?? undefined);
</script>

<svelte:head>
	<title>SMTP-Konfiguration | CodeLinks Admin</title>
</svelte:head>

<section class="space-y-1">
	<h2 class="text-2xl font-semibold tracking-normal">SMTP-Konfiguration</h2>
	<p class="text-muted-foreground text-sm">E-Mail-Versand für Signup, Passwort-Reset und Systemmeldungen.</p>
</section>

<SmtpSettingsForm settings={data.settings} canUpdate={data.canUpdate} form={actionForm} />
