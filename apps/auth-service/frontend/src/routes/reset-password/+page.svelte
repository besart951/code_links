<script lang="ts">
	import { page } from '$app/state';
	import AuthShell from '$lib/components/auth-shell.svelte';
	import ResetPasswordForm from '$lib/components/reset-password-form.svelte';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();
	const token = $derived(page.url.searchParams.get('token') ?? '');
	const actionForm = $derived(form ?? undefined);
</script>

<svelte:head>
	<title>Passwort ändern | CodeLinks</title>
</svelte:head>

<AuthShell title="CodeLinks Passwort" subtitle="Reset-Tokens sind zeitlich begrenzt und werden nach Nutzung ungültig.">
	<ResetPasswordForm form={actionForm} {token} />
</AuthShell>
