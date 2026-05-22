<script lang="ts">
	import AuthShell from '$lib/components/auth-shell.svelte';
	import ForgotPasswordForm from '$lib/components/forgot-password-form.svelte';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();
	const actionForm = $derived(form ?? undefined);
</script>

<svelte:head>
	<title>Passwort vergessen | CodeLinks</title>
</svelte:head>

<AuthShell title="CodeLinks Passwort" subtitle="Reset-Anfragen werden neutral beantwortet und sicher protokolliert.">
	<ForgotPasswordForm form={actionForm} />
</AuthShell>
