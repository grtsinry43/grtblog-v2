<script lang="ts">
	import { onMount } from 'svelte';
	import type { Snippet } from 'svelte';
	import ClientOnly from '$lib/components/ClientOnly.svelte';
	import type { QueryClientConfig } from '@tanstack/svelte-query';

	let { children, options, fallback } = $props<{
		children?: Snippet;
		options?: QueryClientConfig;
		fallback?: Snippet;
	}>();
	let Provider = $state<null | typeof import('@tanstack/svelte-query').QueryClientProvider>(null);
	let client = $state<null | import('@tanstack/svelte-query').QueryClient>(null);
	let ready = $state(false);

	onMount(async () => {
		// Dynamically import for minimal bundle size
		const mod = await import('@tanstack/svelte-query');
		const { QueryClient, QueryClientProvider } = mod;
		client = new QueryClient(options);
		Provider = QueryClientProvider;
		ready = true;
	});
</script>

<ClientOnly fallback={fallback}>
	{#if ready && Provider && client}
		<Provider client={client}>
			{@render children?.()}
		</Provider>
	{/if}
</ClientOnly>
