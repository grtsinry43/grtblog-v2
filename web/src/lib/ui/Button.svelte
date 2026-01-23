<script lang="ts">
	import { type Snippet } from 'svelte';

	interface Props {
		variant?: 'primary' | 'secondary' | 'ghost' | 'icon';
		children: Snippet;
		onclick?: (e: MouseEvent) => void;
		class?: string;
		type?: 'button' | 'submit' | 'reset';
		disabled?: boolean;
		title?: string;
	}

	let {
		variant = 'primary',
		children,
		onclick,
		class: className = '',
		type = 'button',
		disabled = false,
		title
	}: Props = $props();
</script>

<button {type} {disabled} {title} class="btn {variant} {className}" {onclick}>
	{@render children()}
</button>

<style>
	@reference "../../routes/layout.css";

	.btn {
		@apply inline-flex items-center justify-center transition-all duration-300 outline-none;
		@apply active:scale-[0.98] disabled:cursor-not-allowed disabled:opacity-50;
		@apply text-[13px] font-normal tracking-wide;
	}

	.btn.primary {
		@apply rounded-md bg-jade-800 px-3.5 py-1.5 font-medium text-white shadow-sm hover:bg-jade-700 dark:bg-jade-600 dark:hover:bg-jade-500;
	}

	.btn.secondary {
		@apply dark:hover:bg-ink-750 rounded-md border border-ink-100 bg-white px-3.5 py-1.5 text-ink-900 shadow-sm hover:bg-ink-50 dark:border-ink-700 dark:bg-ink-800 dark:text-ink-100;
	}

	.btn.ghost {
		@apply rounded-md px-3 py-1.5 text-ink-600 hover:bg-jade-50 dark:text-ink-400 dark:hover:bg-white/5;
	}

	.btn.icon {
		@apply h-8 w-8 rounded-full p-1.5 text-ink-400 hover:bg-ink-50 dark:text-ink-500 dark:hover:bg-white/5;
	}

	:global(.btn svg) {
		@apply h-4 w-4;
	}
</style>
