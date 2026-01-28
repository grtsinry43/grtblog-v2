<script lang="ts">
	import { type Snippet } from 'svelte';

	interface Props {
		value?: string;
		placeholder?: string;
		type?: string;
		icon?: Snippet;
		class?: string;
		oninput?: (e: Event) => void;
	}

	let {
		value = $bindable(''),
		placeholder = '',
		type = 'text',
		icon,
		class: className = '',
		oninput
	}: Props = $props();
</script>

<div class="input-wrapper group {className}">
	{#if icon}
		<div class="icon-box">
			{@render icon()}
		</div>
	{/if}

	<input bind:value {type} {placeholder} {oninput} class="base-input" class:has-icon={!!icon} />
</div>

<style>
	@reference "../../routes/layout.css";

	.input-wrapper {
		@apply relative;
	}

	.icon-box {
		@apply pointer-events-none absolute top-1/2 left-3.5 -translate-y-1/2 text-ink-300 transition-colors group-focus-within:text-jade-600 dark:text-ink-500 dark:group-focus-within:text-jade-400;
	}

	:global(.icon-box svg) {
		@apply h-3.5 w-3.5;
	}

	.base-input {
		@apply h-9 w-full border border-ink-100/50 bg-ink-50/50 dark:border-ink-800/30 dark:bg-ink-900/40;
		@apply hover:border-ink-200 hover:bg-white dark:hover:border-ink-700 dark:hover:bg-ink-950/60;
		@apply focus:border-jade-500/40 focus:bg-white focus:ring-4 focus:ring-jade-500/5 dark:focus:border-jade-500/40 dark:focus:bg-ink-950;
		@apply rounded-md px-3.5 text-[13px] text-ink-900 transition-all duration-300 outline-none dark:text-ink-100;
		@apply font-normal placeholder:text-ink-300 dark:placeholder:text-ink-600;
	}

	.base-input.has-icon {
		@apply pl-10;
	}
</style>
