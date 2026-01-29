<script lang="ts">
	import { onMount } from 'svelte';
	import type { Snippet } from 'svelte';

	let { children, fallback } = $props<{ children?: Snippet; fallback?: Snippet }>();
	let mounted = $state(false);

	onMount(() => {
		mounted = true;
	});
</script>

{#if mounted}
	{@render children?.()}
{:else if fallback}
	{@render fallback()}
{/if}
