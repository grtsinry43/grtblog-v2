<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Sidebar from '$lib/components/sidebar/Sidebar.svelte';
	import { initTheme, startThemeSync, themeManager } from '$lib/shared/theme.svelte';
	import { onMount } from 'svelte';
	import { consoleLogInfo } from '$lib/features/console-info/index';

	let { children, data } = $props();

	// Initialize theme on mount
	const theme = themeManager;

	onMount(() => {
		initTheme(theme);
		consoleLogInfo();
	});

	startThemeSync(theme);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>GRTBlog</title>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<meta name="description" content="GRTBlog - A personal blog about programming, technology, and software development." />
	<meta name="keywords" content="blog, programming, technology, software development, web development, coding" />
	<meta name="author" content="GRTinry43" />
	<meta property="og:title" content="GRTBlog" />
	<meta property="og:description" content="GRTBlog - A personal blog about programming, technology, and software development." />
	<meta property="og:type" content="website" />
	<meta property="og:url" content="" />
	<meta property="og:image" content="" />
	<meta name="twitter:card" content="summary_large_image" />
	<meta name="twitter:title" content="GRTBlog" />
	<meta name="twitter:description" content="GRTBlog - A personal blog about programming, technology, and software development." />
	<meta name="twitter:image" content="" />
	<script>
		// Inline script to prevent theme flash (fallback before Svelte hydrates)
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

<Sidebar menuTree={data.navMenus ?? []} />
<!-- noise background -->
<div class="bg-noise" aria-hidden="true"></div>

<main class="page-wrapper max-w-[960px] mx-auto px-4 sm:px-6 lg:px-8 py-10 md:py-16">
	<div class="content-container">
		{@render children()}
	</div>
</main>

<style lang="postcss">
	@reference "./layout.css";

	:global(html) {
		scroll-behavior: smooth;
	}
</style>
