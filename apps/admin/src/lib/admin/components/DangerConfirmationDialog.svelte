<script lang="ts">
	import { Button } from '@codelinks/ui/button';
	import * as Dialog from '@codelinks/ui/dialog';

	let {
		title,
		description,
		open = $bindable(false),
		onConfirm
	}: {
		title: string;
		description: string;
		open?: boolean;
		onConfirm?: () => void;
	} = $props();
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Header>
			<Dialog.Title>{title}</Dialog.Title>
			<Dialog.Description>{description}</Dialog.Description>
		</Dialog.Header>
		<Dialog.Footer>
			<Button variant="ghost" onclick={() => (open = false)}>Cancel</Button>
			<Button
				variant="destructive"
				onclick={() => {
					onConfirm?.();
					open = false;
				}}
			>
				Confirm
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>
