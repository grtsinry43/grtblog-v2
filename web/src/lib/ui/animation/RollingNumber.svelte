<script lang="ts">
	import { onMount } from 'svelte';
	import RollingDigit from './RollingDigit.svelte';

	interface Props {
		/** The numeric value to display */
		value: number;
		/** Additional CSS classes */
		class?: string;
		/** Per-digit animation duration in ms */
		duration?: number;
		/** Stagger delay between digits (right-to-left) in ms */
		stagger?: number;
	}

	let { value, class: className = '', duration = 600, stagger = 50 }: Props = $props();

	let animate = $state(false);

	onMount(() => {
		// Enable transitions after first paint to avoid animating SSR hydration
		requestAnimationFrame(() => {
			animate = true;
		});
	});

	const safeValue = $derived(Number.isFinite(value) ? Math.max(0, Math.floor(value)) : 0);
	const digits = $derived(String(safeValue).split('').map(Number));
</script>

<span class="rolling-number {className}">
	<span class="sr-only">{safeValue}</span>
	{#each digits as digit, i (digits.length - 1 - i)}
		<RollingDigit {digit} {duration} delay={(digits.length - 1 - i) * stagger} {animate} />
	{/each}
</span>

<style>
	.rolling-number {
		display: inline-flex;
		font-variant-numeric: tabular-nums;
		line-height: 1;
	}
</style>
