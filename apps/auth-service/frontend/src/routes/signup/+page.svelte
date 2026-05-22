<script lang="ts">
	import AuthShell from '$lib/components/auth-shell.svelte';
	import SignupForm from '$lib/components/signup-form.svelte';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();
	const actionForm = $derived(form ?? undefined);
</script>

<svelte:head>
	<title>Signup | CodeLinks</title>
</svelte:head>

<AuthShell title="CodeLinks Konto" subtitle="Registrierung mit verpflichtender E-Mail-Bestätigung.">
	<SignupForm form={actionForm} />
</AuthShell>
