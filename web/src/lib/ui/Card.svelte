<script lang="ts">
	import { type Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		class?: string;
		variant?: 'glass' | 'solid' | 'seamless';
		hover?: boolean;
	}

	let { children, class: className = '', variant = 'glass', hover = false }: Props = $props();
</script>

<div class="card {variant} {hover ? 'hoverable' : ''} {className}">
	{@render children()}
</div>

<style>
	@reference "../../routes/layout.css";

	.card {
		@apply relative overflow-hidden transition-all duration-500 ease-[cubic-bezier(0.23,1,0.32,1)];
	}

	.card.glass {
		@apply rounded-lg border border-white/30 bg-white/40 shadow-sm backdrop-blur-md dark:border-white/5 dark:bg-ink-950/40;
	}

	.card.solid {
		@apply rounded-lg border border-ink-50 bg-white shadow-sm dark:border-ink-800 dark:bg-ink-900;
	}

	.card.seamless {
		@apply rounded-none border-none bg-transparent shadow-none;
	}

	.card.hoverable.glass:hover {
		@apply -translate-y-0.5 border-jade-200/50 shadow-md dark:border-jade-800/20;
	}

	.card.hoverable.solid:hover {
		@apply -translate-y-0.5 border-ink-100 shadow-md dark:border-ink-700;
	}

	.card.hoverable.seamless:hover {
		@apply bg-ink-50/50 dark:bg-white/5;
	}
</style>
