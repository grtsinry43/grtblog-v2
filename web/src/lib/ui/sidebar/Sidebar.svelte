<script lang="ts">
	import {
		Home,
		BookOpen,
		PenTool,
		Archive,
		Image,
		User,
		Menu,
		X,
		ChevronDown,
		Terminal,
		Coffee,
		Sparkles,
		Code,
		List,
		Sun,
		Moon,
		Monitor
	} from 'lucide-svelte';
	import { slide, fade, fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { SectionId } from './types';
	import TableOfContents from './TableOfContents.svelte';
	import { themeManager } from '$lib/shared/theme.svelte';

	type NavItem = {
		id: SectionId;
		label: string;
		icon: any;
		desc: string;
		children?: { id: SectionId; label: string; icon?: any }[];
	};

	let { activeSection, onNavigate } = $props<{
		activeSection: SectionId;
		onNavigate: (id: SectionId) => void;
	}>();

	let isMobileMenuOpen = $state(false);
	let hoveredItemId = $state<SectionId | null>(null);
	let expandedMobileItems = $state<Set<SectionId>>(new Set());
	let scrollY = $state(0);
	let navProgress = $derived(Math.max(0, Math.min(scrollY / 60, 1)));
	let isTocOpen = $state(false);

	let isMenuAnimating = $state(false);
	$effect(() => {
		isMobileMenuOpen;
		isMenuAnimating = true;
		const timer = setTimeout(() => (isMenuAnimating = false), 500);
		return () => clearTimeout(timer);
	});

	const navItems: NavItem[] = [
		{ id: SectionId.Home, label: '首页', icon: Home, desc: '回到开始的地方' },
		{
			id: SectionId.Articles,
			label: '文章',
			icon: BookOpen,
			desc: '深度思考与技术沉淀',
			children: [
				{ id: SectionId.ArtTech, label: '技术', icon: Terminal },
				{ id: SectionId.ArtLife, label: '生活', icon: Coffee }
			]
		},
		{
			id: SectionId.Notes,
			label: '手记',
			icon: PenTool,
			desc: '碎片化的灵感记录',
			children: [
				{ id: SectionId.NoteCode, label: '代码片段', icon: Code },
				{ id: SectionId.NoteThought, label: '随想', icon: Sparkles }
			]
		},
		{ id: SectionId.Archives, label: '归档', icon: Archive, desc: '时间的足迹' },
		{ id: SectionId.Gallery, label: '相册', icon: Image, desc: '光影的瞬间' },
		{ id: SectionId.About, label: '关于', icon: User, desc: '我是谁' }
	];

	function handleNavigate(id: SectionId) {
		onNavigate(id);
		isMobileMenuOpen = false;
	}

	function toggleMobileSubmenu(e: Event, id: SectionId) {
		e.stopPropagation();
		const newSet = new Set(expandedMobileItems);
		if (newSet.has(id)) newSet.delete(id);
		else newSet.add(id);
		expandedMobileItems = newSet;
	}

	let activeLabel = $derived(
		navItems.find((n) => n.id === activeSection || n.children?.some((c) => c.id === activeSection))
			?.label || '博客'
	);

	let displayTitle = $derived.by(() => {
		const post = $page.data?.post;
		const article = $page.data?.article;
		return post?.title || article?.title || activeLabel;
	});

	// Get real post data for TOC
	let postData = $derived($page.data?.post);

	// Theme logic
	const toggleTheme = () => {
		const themes: ('light' | 'dark' | 'system')[] = ['light', 'dark', 'system'];
		const nextIndex = (themes.indexOf(themeManager.current) + 1) % themes.length;
		themeManager.set(themes[nextIndex]);
	};

	let ThemeIcon = $derived.by(() => {
		if (themeManager.current === 'light') return Sun;
		if (themeManager.current === 'dark') return Moon;
		return Monitor;
	});
</script>

<!-- ================= Desktop Dock ================= -->
<nav class="desktop-dock-container hidden lg:block">
	<div class="glass-dock" onmouseleave={() => (hoveredItemId = null)} role="presentation">
		<!-- Avatar -->
		<button class="avatar-btn" onclick={() => handleNavigate(SectionId.Home)}>
			<img src="https://dogeoss.grtsinry43.com/img/author.jpeg" alt="头像" class="avatar-img" />
		</button>

		<div class="dock-divider"></div>

		{#each navItems as item (item.id)}
			{@const isActive =
				activeSection === item.id || item.children?.some((c) => c.id === activeSection)}

			<div
				class="nav-item-wrapper"
				onmouseenter={() => (hoveredItemId = item.id)}
				role="presentation"
			>
				<button onclick={() => handleNavigate(item.id)} class="nav-btn {isActive ? 'active' : ''}">
					<item.icon size={16} strokeWidth={isActive ? 2.5 : 2} class="icon-transition" />
				</button>

				{#if item.children && hoveredItemId === item.id}
					<div transition:fly={{ x: -4, duration: 200, opacity: 0 }} class="popover-wrapper">
						<div class="glass-popover">
							{#each item.children as subItem}
								{@const isSubActive = activeSection === subItem.id}
								<button
									onclick={(e) => {
										e.stopPropagation();
										handleNavigate(subItem.id);
									}}
									class="popover-item {isSubActive ? 'active' : ''}"
								>
									{subItem.label}
								</button>
							{/each}
						</div>
					</div>
				{:else if !item.children && hoveredItemId === item.id}
					<span transition:fade={{ duration: 150 }} class="dock-tooltip">
						{item.label}
					</span>
				{/if}
			</div>
		{/each}

		<div class="dock-divider mt-1"></div>

		<!-- Theme Toggle -->
		<button onclick={toggleTheme} class="nav-btn mt-1" title="切换主题 ({themeManager.current})">
			<ThemeIcon size={16} />
		</button>
	</div>
</nav>

<!-- ================= Mobile Floating Capsule ================= -->
<svelte:window bind:scrollY />

<div
	class="mobile-nav-container lg:hidden"
	class:is-open={isMobileMenuOpen}
	class:is-animating={isMenuAnimating}
	style:top={isMobileMenuOpen ? '0' : `${12 * (1 - navProgress)}px`}
	style:left={isMobileMenuOpen ? '0' : `${12 * (1 - navProgress)}px`}
	style:right={isMobileMenuOpen ? '0' : `${12 * (1 - navProgress)}px`}
>
	<div
		class="capsule-wrapper"
		class:shadow-xl={isMobileMenuOpen}
		class:rounded-none={isMobileMenuOpen}
		class:h-screen={isMobileMenuOpen}
		style:border-radius={isMobileMenuOpen ? '0' : `${12 * (1 - navProgress)}px`}
	>
		<!-- Background -->
		<div class="capsule-bg" style:height={isMobileMenuOpen ? '100vh' : '2.5rem'}></div>

		<!-- Header Bar -->
		<div class="capsule-header">
			<div class="header-left">
				<button
					onclick={(e) => {
						e.stopPropagation();
						isMobileMenuOpen = !isMobileMenuOpen;
					}}
					class="avatar-trigger"
				>
					<div class="avatar-capsule">
						<img
							src="https://dogeoss.grtsinry43.com/img/author.jpeg"
							alt="头像"
							class="avatar-img-mobile"
						/>
					</div>
				</button>

				<div class="title-wrapper" class:hidden={isMobileMenuOpen}>
					<span class="capsule-title">{displayTitle}</span>
				</div>
			</div>

			<div class="header-right">
				{#if postData?.toc?.length}
					<button
						onclick={(e) => {
							e.stopPropagation();
							isTocOpen = true;
						}}
						class="icon-btn"
					>
						<List size={16} />
					</button>
				{/if}

				<button
					onclick={(e) => {
						e.stopPropagation();
						isMobileMenuOpen = !isMobileMenuOpen;
					}}
					class="icon-btn"
				>
					{#if isMobileMenuOpen}
						<X size={16} />
					{:else}
						<Menu size={16} />
					{/if}
				</button>
			</div>
		</div>

		<!-- Expanded Menu -->
		{#if isMobileMenuOpen}
			<div transition:slide={{ duration: 400 }} class="mobile-menu-content no-scrollbar">
				<div class="menu-divider-row">
					<span class="menu-divider-label">主导航</span>
					<div class="flex items-center gap-4">
						<button
							onclick={toggleTheme}
							class="text-[10px] text-ink-400 flex items-center gap-1.5 active:text-jade-600 transition-colors uppercase tracking-widest font-mono"
						>
							<ThemeIcon size={12} />
							{themeManager.current}
						</button>
						<span class="menu-divider-tag">MENU</span>
					</div>
				</div>

				<div class="menu-list">
					{#each navItems as item}
						{@const isActive =
							activeSection === item.id || item.children?.some((c) => c.id === activeSection)}
						{@const hasChildren = !!item.children}
						{@const isExpanded = expandedMobileItems.has(item.id)}

						<div class="menu-item-group">
							<div class="menu-item {isActive ? 'active' : ''}">
								<button type="button" onclick={() => handleNavigate(item.id)} class="menu-item-btn">
									<div class="menu-icon-box {isActive ? 'active' : ''}">
										<item.icon size={13} />
									</div>
									<span class="menu-label {isActive ? 'active' : ''}">{item.label}</span>
								</button>

								{#if hasChildren}
									<button
										type="button"
										onclick={(e) => toggleMobileSubmenu(e, item.id)}
										class="submenu-toggle"
									>
										<ChevronDown size={14} class="toggle-icon {isExpanded ? 'expanded' : ''}" />
									</button>
								{/if}
							</div>

							{#if hasChildren && isExpanded}
								<div transition:slide={{ duration: 300 }} class="submenu-list">
									<div class="tree-line"></div>
									{#each item.children as subItem}
										{@const isSubActive = activeSection === subItem.id}
										<button
											onclick={() => handleNavigate(subItem.id)}
											class="submenu-item {isSubActive ? 'active' : ''}"
										>
											<div class="branch-line"></div>
											<span>{subItem.label}</span>
										</button>
									{/each}
								</div>
							{/if}
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>

	{#if isMobileMenuOpen}
		<button
			type="button"
			transition:fade={{ duration: 300 }}
			class="mobile-nav-backdrop"
			onclick={() => (isMobileMenuOpen = false)}
		></button>
	{/if}
</div>

<!-- TOC Component with Real Data -->
<TableOfContents
	isOpen={isTocOpen}
	onClose={() => (isTocOpen = false)}
	toc={postData?.toc}
	activeAnchor={$page.url.hash.slice(1)}
/>

<style>
	@reference "../../../routes/layout.css";

	/* --- Desktop Dock --- */
	.desktop-dock-container {
		@apply fixed top-1/2 left-4 z-50 hidden -translate-y-1/2 lg:block;
	}

	.glass-dock {
		@apply flex flex-col items-center gap-2 rounded-full border border-ink-100 bg-white/80 px-1.5 py-2.5 shadow-sm backdrop-blur-2xl dark:border-ink-800/50 dark:bg-ink-900/80;
	}

	.avatar-btn {
		@apply mb-1 h-8 w-8 cursor-pointer transition-transform hover:scale-105;
	}

	.avatar-img {
		@apply h-full w-full rounded-full border border-white/40 shadow-sm;
	}

	.dock-divider {
		@apply mb-1 h-px w-4 bg-ink-100 dark:bg-ink-800/50;
	}

	.nav-item-wrapper {
		@apply relative flex items-center;
	}

	.nav-btn {
		@apply flex h-8 w-8 items-center justify-center rounded-full text-ink-400 transition-all hover:bg-jade-50 hover:text-jade-600 dark:text-ink-500 dark:hover:bg-white/5 dark:hover:text-jade-400;
	}

	.nav-btn.active {
		@apply scale-105 bg-jade-800 text-white shadow-md dark:bg-jade-600;
	}

	.popover-wrapper {
		@apply absolute top-1/2 left-10 z-10 -translate-y-1/2 pl-2;
	}

	.glass-popover {
		@apply flex min-w-[100px] flex-col rounded-lg border border-ink-50 bg-white/95 py-1 shadow-xl backdrop-blur-xl dark:border-ink-800 dark:bg-ink-900/95;
	}

	.popover-item {
		@apply px-4 py-1.5 text-left text-[12px] font-normal text-ink-600 transition-colors hover:bg-jade-50 hover:text-jade-700 dark:text-ink-400 dark:hover:bg-white/5 dark:hover:text-jade-300;
	}

	.popover-item.active {
		@apply bg-jade-50/50 font-bold text-jade-700 dark:bg-jade-900/20 dark:text-jade-400;
	}

	.dock-tooltip {
		@apply pointer-events-none absolute top-1/2 left-12 -translate-y-1/2 rounded bg-ink-900 px-2 py-0.5 font-serif text-[10px] tracking-widest whitespace-nowrap text-white shadow-lg dark:bg-white dark:text-ink-900;
	}

	/* --- Mobile Capsule --- */
	.mobile-nav-container {
		@apply fixed inset-x-0 z-[90] flex justify-center transition-all ease-[cubic-bezier(0.23,1,0.32,1)] lg:hidden;
	}

	.mobile-nav-container.is-open {
		@apply inset-x-0;
	}

	.mobile-nav-container.is-animating {
		@apply duration-500;
	}

	.capsule-wrapper {
		@apply relative mx-auto w-full overflow-hidden transition-all ease-[cubic-bezier(0.23,1,0.32,1)];
	}

	.capsule-bg {
		@apply absolute inset-0 border border-ink-50 bg-white/90 shadow-sm backdrop-blur-xl transition-all duration-500 dark:border-ink-800/40 dark:bg-ink-900/95;
	}

	.capsule-header {
		@apply relative z-10 flex h-10 items-center justify-between px-3;
	}

	.header-left {
		@apply flex items-center gap-2;
	}

	.avatar-trigger {
		@apply flex h-7 w-7 items-center justify-center rounded-full transition-transform active:scale-95;
	}

	.avatar-capsule {
		@apply h-6 w-6 shrink-0 overflow-hidden rounded-full border border-ink-50 dark:border-ink-800;
	}

	.avatar-img-mobile {
		@apply h-full w-full object-cover;
	}

	.title-wrapper {
		@apply flex flex-col justify-center transition-all duration-300;
	}

	.capsule-title {
		@apply max-w-[150px] truncate font-serif text-[12px] leading-none font-bold text-ink-900 dark:text-jade-100;
	}

	.header-right {
		@apply flex items-center gap-0.5;
	}

	.icon-btn {
		@apply flex h-8 w-8 items-center justify-center rounded-full text-ink-500 active:bg-ink-50 dark:text-ink-400 dark:active:bg-white/5;
	}

	.mobile-menu-content {
		@apply relative z-10 no-scrollbar flex max-h-[70vh] flex-col overflow-y-auto px-2 pb-6;
	}

	.menu-divider-row {
		@apply mb-2 flex items-center justify-between border-b border-ink-50 px-4 pt-1 pb-3 opacity-80 dark:border-ink-800/30;
	}

	.menu-divider-label {
		@apply text-[10px] font-bold tracking-widest text-ink-400 uppercase;
	}

	.menu-divider-tag {
		@apply font-mono text-[8px] text-ink-300;
	}

	.menu-list {
		@apply flex flex-col gap-0.5;
	}

	.menu-item-group {
		@apply flex flex-col;
	}

	.menu-item {
		@apply relative flex items-center gap-2 overflow-hidden rounded-lg px-3 py-1.5 transition-all;
	}

	.menu-item.active {
		@apply border border-jade-100/50 bg-jade-50/50 dark:border-jade-800/20 dark:bg-jade-900/10;
	}

	.menu-item-btn {
		@apply flex min-w-0 flex-1 items-center gap-3 text-left;
	}

	.menu-icon-box {
		@apply flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-ink-50 text-ink-400 transition-colors dark:bg-ink-800;
	}

	.menu-icon-box.active {
		@apply bg-jade-600 text-white shadow-sm;
	}

	.menu-label {
		@apply font-serif text-[13px] font-medium text-ink-700 dark:text-ink-300;
	}

	.menu-label.active {
		@apply text-jade-800 dark:text-jade-100;
	}

	.submenu-toggle {
		@apply -mr-1 rounded-full p-1.5 text-ink-300;
	}

	.toggle-icon {
		@apply transition-transform duration-300;
	}

	.toggle-icon.expanded {
		@apply rotate-180 text-jade-600;
	}

	.submenu-list {
		@apply relative mt-0.5 mb-1.5 flex flex-col gap-0.5;
	}

	.tree-line {
		@apply absolute top-0 bottom-3 left-[31px] w-[1px] bg-ink-100 dark:bg-ink-800/50;
	}

	.submenu-item {
		@apply relative mx-1 flex items-center gap-2 rounded-md py-1.5 pr-4 pl-11 text-left text-[12px] font-normal text-ink-500 transition-colors hover:text-jade-600 dark:text-ink-400;
	}

	.submenu-item.active {
		@apply font-bold text-jade-700 dark:text-jade-300;
	}

	.branch-line {
		@apply absolute top-1/2 left-[24px] h-[1px] w-3 bg-ink-100 dark:bg-ink-800/50;
	}

	.mobile-nav-backdrop {
		@apply fixed inset-0 -z-10 bg-ink-950/20 backdrop-blur-[1px];
	}

	.empty-state {
		@apply p-8 text-center text-xs font-light text-ink-300 italic;
	}
</style>
