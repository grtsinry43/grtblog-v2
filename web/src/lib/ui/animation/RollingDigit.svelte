<script lang="ts">
	import { tick, untrack } from 'svelte';

	interface Props {
		digit: number;
		duration?: number;
		delay?: number;
		/** When true, digit changes animate; when false, they snap instantly */
		animate?: boolean;
	}

	let { digit, duration = 600, delay = 0, animate = false }: Props = $props();

	// 0-9 repeated 3 times; offset resets at 20 via transitionend to stay within range
	const CELLS = Array.from({ length: 30 }, (_, i) => i % 10);

	let offset = $state(untrack(() => digit));
	let prevDigit = untrack(() => digit);
	// Internal override: temporarily disables transition during offset reset
	let transitionLock = $state(false);

	$effect(() => {
		const d = digit;
		if (d !== prevDigit) {
			if (animate) {
				// Cumulative offset → always scrolls upward
				const diff = d >= prevDigit ? d - prevDigit : 10 - prevDigit + d;
				offset += diff;
			} else {
				// Snap: jump directly to digit position
				offset = d;
			}
			prevDigit = d;
		}
	});

	const effectiveAnimate = $derived(animate && !transitionLock);

	// After transition, snap offset back to low range if it grew too large
	const handleTransitionEnd = async () => {
		if (offset < 20) return;
		transitionLock = true;
		offset = offset % 10;
		await tick();
		requestAnimationFrame(() => {
			transitionLock = false;
		});
	};
</script>

<span class="rd-col" class:rd-ready={effectiveAnimate}>
	<span
		class="rd-strip"
		style:transform="translateY({-offset}em)"
		style:transition-duration="{duration}ms"
		style:transition-delay="{delay}ms"
		ontransitionend={handleTransitionEnd}
	>
		{#each CELLS as d, i (`${i}-${d}`)}
			<span class="rd-cell" aria-hidden="true">{d}</span>
		{/each}
	</span>
</span>

<style>
.rd-col {
	display: inline-block;
	height: 1em;
	line-height: 1;
	overflow: hidden;
	user-select: none;
}

	.rd-strip {
		display: flex;
		flex-direction: column;
		will-change: transform;
	}

	.rd-col.rd-ready .rd-strip {
		transition-property: transform;
		transition-timing-function: cubic-bezier(0.22, 1, 0.36, 1);
	}

	.rd-cell {
		height: 1em;
		line-height: 1;
		text-align: center;
	}
</style>
