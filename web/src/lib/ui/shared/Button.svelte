<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		onclick?: () => void;
		variant?: 'primary' | 'secondary' | 'ghost';
		className?: string;
		type?: 'button' | 'submit' | 'reset';
	}

	let { children, onclick, variant = 'primary', className = '', type = 'button' }: Props = $props();

	const variants = {
		primary: 'bg-jade-500 hover:bg-jade-600 text-white shadow-glow',
		secondary:
			'bg-white/80 dark:bg-ink-900/80 backdrop-blur-xl border border-white/50 dark:border-ink-700 text-ink-900 dark:text-ink-100 shadow-glass',
		ghost:
			'bg-transparent hover:bg-jade-50/50 dark:hover:bg-jade-900/20 text-ink-600 dark:text-ink-400'
	};
</script>

<button
	{type}
	on:click={onclick}
	class="px-5 py-2.5 rounded-full font-medium transition-all duration-300 active:scale-95 disabled:opacity-50 {variants[
		variant
	]} {className}"
>
	{@render children()}
</button>
