<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import AmbientBackground from '$lib/ui/background/AmbientBackground.svelte';
	import Sidebar from '$lib/ui/sidebar/Sidebar.svelte';
	import { SectionId } from '$lib/ui/sidebar/types';

	let { children } = $props();
</script>

<AmbientBackground />

<Sidebar activeSection={SectionId.Home} onNavigate={() => {}} />

<svelte:head>
	<link rel="icon" href={favicon} />
	<link rel="preconnect" href="https://fonts.googleapis.com" />
	<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin="anonymous" />
	<link
		href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600&family=Noto+Sans+SC:wght@300;400;500;700&family=Noto+Serif+SC:wght@400;500;600;700&family=JetBrains+Mono:wght@400;500&display=swap"
		rel="stylesheet"
	/>
</svelte:head>

<main class="relative z-10 w-full min-h-screen">
	{@render children()}
</main>

<style>
	main {
		/* Content will have its own padding/max-width */
	}
</style>
