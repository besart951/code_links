<script lang="ts">
	import { Badge } from '@codelinks/ui/badge';

	let {
		label,
		tone = 'neutral'
	}: {
		label: string;
		tone?: 'neutral' | 'success' | 'warning' | 'danger' | 'info';
	} = $props();

	const classes = $derived(
		({
			neutral: 'border-border bg-background text-foreground',
			success: 'border-success-border bg-success-muted text-success-muted-foreground',
			warning: 'border-warning-border bg-warning-muted text-warning-muted-foreground',
			danger: 'border-destructive/30 bg-destructive/10 text-destructive',
			info: 'border-info-border bg-info-muted text-info-muted-foreground'
		})[tone]
	);
</script>

<Badge variant="outline" class={classes}>{label}</Badge>
