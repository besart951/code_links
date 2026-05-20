<script lang="ts">
	import { hasActionReason } from '$lib/admin/security.js';
	import { Button } from '@codelinks/ui/button';
	import * as Dialog from '@codelinks/ui/dialog';
	import { Textarea } from '@codelinks/ui/textarea';

	let {
		action,
		title,
		description,
		open = $bindable(false),
		onConfirm
	}: {
		action: string;
		title: string;
		description: string;
		open?: boolean;
		onConfirm?: (reason: string) => void;
	} = $props();

	let reason = $state('');
	const valid = $derived(hasActionReason(action, reason));
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>{title}</Dialog.Title>
			<Dialog.Description>{description}</Dialog.Description>
		</Dialog.Header>
		<Textarea bind:value={reason} rows={4} placeholder="Reason" />
		<Dialog.Footer>
			<Button variant="ghost" onclick={() => (open = false)}>Cancel</Button>
			<Button
				disabled={!valid}
				onclick={() => {
					if (!valid) return;
					onConfirm?.(reason.trim());
					open = false;
					reason = '';
				}}
			>
				Confirm
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
