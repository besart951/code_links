<script lang="ts">
	import { page } from '$app/state';
	import AuthShell from '$lib/components/auth-shell.svelte';
	import LoginForm from '$lib/components/login-form.svelte';
	import type { ActionData } from './$types';

	let { form }: { form: ActionData } = $props();
	const redirectTo = $derived(page.url.searchParams.get('redirectTo') ?? '/');
	const actionForm = $derived(form ?? undefined);
</script>

<svelte:head>
	<title>Login | CodeLinks</title>
</svelte:head>

<AuthShell title="CodeLinks Admin & Produkte" subtitle="Ein Konto für Auth, Lizenzen und Admin-Konsole.">
	<LoginForm form={actionForm} {redirectTo} />
</AuthShell>
