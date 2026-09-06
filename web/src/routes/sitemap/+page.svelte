<script lang="ts">
	import SitemapDirectory from '$lib/features/discovery/components/SitemapDirectory.svelte';
	import type { PageData } from './$types';
	let { data }: { data: PageData } = $props();
</script>

<div class="mx-auto w-full max-w-6xl px-4 py-10 sm:px-6 sm:py-16">
	<SitemapDirectory catalog={data.catalog} />
</div>
