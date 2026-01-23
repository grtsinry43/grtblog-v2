<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Sidebar from '$lib/ui/sidebar/Sidebar.svelte';
	import { SectionId } from '$lib/ui/sidebar/types';
	import { themeManager } from '$lib/shared/theme.svelte';

	let { children } = $props();

	// Initialize theme on mount
	const theme = themeManager;
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<script>
		// Inline script to prevent theme flash
		(function () {
			try {
				const theme = localStorage.getItem('theme') || 'system';
				const isDark =
					theme === 'dark' ||
					(theme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches);
				document.documentElement.classList.toggle('dark', isDark);
			} catch (e) {}
		})();
	</script>
</svelte:head>

<Sidebar activeSection={SectionId.Home} onNavigate={() => {}} />

<main class="page-wrapper">
	<div class="content-container">
		{@render children()}
	</div>
</main>

<style>
	@reference "./layout.css";

	:global(html) {
		scroll-behavior: smooth;
	}

	.page-wrapper {
		@apply min-h-screen bg-white transition-colors duration-500 dark:bg-ink-950;
	}

	.content-container {
		@apply mx-auto max-w-4xl px-5 py-10 md:px-10 md:py-16 lg:py-20;
	}
</style>
