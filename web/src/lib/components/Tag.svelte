<script lang="ts">
	import { type Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		variant?: 'jade' | 'ink' | 'amber' | 'cinnabar' | 'outline';
		class?: string;
	}

	let { children, variant = 'ink', class: className = '' }: Props = $props();
</script>

<span class="tag {variant} {className}">
	{@render children()}
</span>

<style>
	@reference "../../routes/layout.css";

	.tag {
		@apply inline-flex items-center rounded px-1.5 py-0.5 font-mono text-[9.5px] font-bold tracking-[0.1em] uppercase transition-colors;
	}

	.tag.jade {
		@apply border border-jade-500/10 bg-jade-500/10 text-jade-700 dark:text-jade-400;
	}

	.tag.ink {
		@apply border border-ink-100 bg-ink-500/5 text-ink-500 dark:border-ink-800/50 dark:text-ink-400;
	}

	.tag.amber {
		@apply border border-amber-500/20 bg-amber-500/10 text-amber-700 dark:text-amber-400;
	}

	.tag.cinnabar {
		@apply text-cinnabar-700 dark:text-cinnabar-400 border border-cinnabar-500/20 bg-cinnabar-500/10;
	}

	.tag.outline {
		@apply border border-ink-100 bg-transparent text-ink-400 dark:border-ink-800/50;
	}
</style>
