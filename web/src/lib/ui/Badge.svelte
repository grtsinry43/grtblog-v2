<script lang="ts">
	import { type Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		variant?: 'ghost' | 'soft' | 'dot';
		class?: string;
	}

	let { children, variant = 'soft', class: className = '' }: Props = $props();
</script>

<div class="badge {variant} {className}">
	{#if variant === 'dot'}
		<span class="dot"></span>
	{/if}
	<span class="content">{@render children()}</span>
</div>

<style>
	@reference "../../routes/layout.css";

	.badge {
		@apply inline-flex items-center gap-1.5 rounded-full px-2 py-0.5;
	}

	.content {
		@apply font-mono text-[9px] font-bold tracking-[0.1em] uppercase;
	}

	.badge.soft {
		@apply border border-jade-500/10 bg-jade-500/5 text-jade-700 dark:text-jade-400;
	}

	.badge.ghost {
		@apply border border-ink-100 bg-transparent text-ink-400 dark:border-ink-800/50;
	}

	.badge.dot {
		@apply bg-transparent px-0 text-ink-500 dark:text-ink-400;
	}

	.dot {
		@apply h-1 w-1 animate-pulse rounded-full bg-jade-500 dark:bg-jade-400;
	}

	:global(.badge svg) {
		@apply h-2.5 w-2.5;
	}
</style>
