<script lang="ts">
	interface Props {
		vertical?: boolean;
		class?: string;
		label?: string;
	}

	let { vertical = false, class: className = '', label = '' }: Props = $props();
</script>

<div class="divider-container {vertical ? 'vertical' : 'horizontal'} {className}">
	<div class="line"></div>
	{#if label}
		<span class="label">{label}</span>
		<div class="line"></div>
	{/if}
</div>

<style>
	@reference "../../routes/layout.css";

	.divider-container {
		@apply flex items-center gap-4;
	}

	.divider-container.horizontal {
		@apply w-full;
	}

	.divider-container.vertical {
		@apply h-full flex-col;
	}

	.line {
		@apply bg-ink-100 transition-colors dark:bg-ink-800/50;
	}

	.horizontal .line {
		@apply h-px flex-1;
	}

	.vertical .line {
		@apply w-px flex-1;
	}

	.label {
		@apply font-mono text-[9px] tracking-[0.3em] whitespace-nowrap text-ink-400 uppercase;
	}
</style>
