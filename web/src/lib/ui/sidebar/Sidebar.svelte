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
		ChevronRight,
		ChevronDown,
		Hash,
		Code,
		Coffee,
		Sparkles,
		Terminal,
		Camera,
		Zap,
		List
	} from 'lucide-svelte';
	import { slide, fade, fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { SectionId } from './types';
	import TableOfContents from './TableOfContents.svelte';

	type NavItem = {
		id: SectionId;
		label: string;
		icon: any;
		desc: string;
		children?: { id: SectionId; label: string; icon?: any }[];
	};

	// --- 2. Props ---
	let { activeSection, onNavigate } = $props<{
		activeSection: SectionId;
		onNavigate: (id: SectionId) => void;
	}>();

	// --- 3. 状态管理 (Runes) ---
	let isMobileMenuOpen = $state(false);
	let scrollProgress = $state(0);
	let hoveredItemId = $state<SectionId | null>(null); // Desktop Hover 状态
	let expandedMobileItems = $state<Set<SectionId>>(new Set()); // Mobile Accordion 状态
	// let isScrolled = $state(false); // Removed: using direct scrollY binding
	let scrollY = $state(0);
	// Interpolation progress: 0 (top/capsule) -> 1 (scrolled/full)
	// Threshold: 50px
	let navProgress = $derived(Math.max(0, Math.min(scrollY / 50, 1)));
	let isTocOpen = $state(false); // Mobile TOC 状态

	// Animation state to fix "snap to close" bug
	let isMenuAnimating = $state(false);
	$effect(() => {
		// Trigger animation flag when menu state changes
		isMobileMenuOpen; // Dependency
		isMenuAnimating = true;
		const timer = setTimeout(() => (isMenuAnimating = false), 500); // 500ms matches duration-500
		return () => clearTimeout(timer);
	});

	// --- 4. 导航数据配置 ---
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

	// --- 5. 辅助逻辑 ---

	// 滚动监听
	onMount(() => {
		const handleScroll = () => {
			const h = document.documentElement;
			const st = h.scrollTop || document.body.scrollTop;
			const sh = h.scrollHeight || document.body.scrollHeight;
			const ch = h.clientHeight;
			// 防止除以0
			const percent = sh - ch > 0 ? st / (sh - ch) : 0;
			scrollProgress = Math.min(Math.max(percent, 0), 1);

			// 更新滚动状态 (阈值 20px)
		};
		window.addEventListener('scroll', handleScroll);
		return () => window.removeEventListener('scroll', handleScroll);
	});

	// 导航处理
	function handleNavigate(id: SectionId) {
		onNavigate(id);
		isMobileMenuOpen = false;
	}

	// 移动端折叠切换
	function toggleMobileSubmenu(e: Event, id: SectionId) {
		e.stopPropagation();
		const newSet = new Set(expandedMobileItems);
		if (newSet.has(id)) newSet.delete(id);
		else newSet.add(id);
		expandedMobileItems = newSet;
	}

	// 计算当前显示的标题
	let activeLabel = $derived(
		navItems.find((n) => n.id === activeSection || n.children?.some((c) => c.id === activeSection))
			?.label || '博客'
	);

	// 优先显示文章标题
	let displayTitle = $derived.by(() => {
		// Debug logging to help identify why title might be missing
		// console.log('Sidebar Page Data:', $page.data);
		const post = $page.data?.post;
		const article = $page.data?.article;
		return post?.title || article?.title || activeLabel;
	});
</script>

<!-- ================= Desktop Dock (左侧悬浮) ================= -->
<nav class="fixed left-4 top-1/2 -translate-y-1/2 z-50 hidden lg:block">
	<div
		class="glass-dock py-3 px-1.5 flex flex-col gap-3 items-center hover:shadow-float transition-all duration-500"
		onmouseleave={() => (hoveredItemId = null)}
		role="presentation"
	>
		<!-- Avatar -->
		<button
			class="w-10 h-10 mb-2 cursor-pointer hover:scale-105 transition-transform"
			onclick={() => handleNavigate(SectionId.Home)}
		>
			<img
				src="https://dogeoss.grtsinry43.com/img/author.jpeg"
				alt="Avatar"
				class="rounded-full w-full h-full border-2 border-white/30"
			/>
		</button>

		<div class="w-full h-[1px] bg-ink-200/20 dark:bg-white/10 mb-1"></div>

		{#each navItems as item (item.id)}
			{@const isActive =
				activeSection === item.id || item.children?.some((c) => c.id === activeSection)}

			<!-- 菜单项容器 (处理 Hover) -->
			<div
				class="relative flex items-center"
				onmouseenter={() => (hoveredItemId = item.id)}
				role="presentation"
			>
				<button
					onclick={() => handleNavigate(item.id)}
					class="nav-btn-desktop relative z-20 {isActive ? 'active' : ''}"
				>
					<item.icon
						size={20}
						strokeWidth={isActive ? 2.5 : 2}
						class="transition-transform duration-300 group-hover:scale-110"
					/>
				</button>

				<!-- A. 二级菜单气泡 (Glass Popover) -->
				{#if item.children && hoveredItemId === item.id}
					<div
						transition:fly={{ x: -5, duration: 200, opacity: 0 }}
						class="absolute left-12 top-1/2 -translate-y-1/2 z-10 pl-2"
					>
						<div class="glass-popover p-1 min-w-[140px] flex flex-col gap-0.5">
							{#each item.children as subItem}
								{@const isSubActive = activeSection === subItem.id}
								<button
									onclick={(e) => {
										e.stopPropagation();
										handleNavigate(subItem.id);
									}}
									class="relative z-10 flex items-center gap-2.5 px-3 py-1.5 rounded-md text-sm font-medium transition-colors text-left
                   {isSubActive
										? 'bg-jade-50 text-jade-700 dark:bg-jade-900/30 dark:text-jade-300'
										: 'text-ink-600 dark:text-ink-400 hover:bg-ink-100 dark:hover:bg-white/5'}"
								>
									{#if subItem.icon}
										<subItem.icon size={14} class={isSubActive ? 'opacity-100' : 'opacity-70'} />
									{/if}
									{subItem.label}
								</button>
							{/each}
						</div>
					</div>

					<!-- B. 普通 Tooltip -->
				{:else if !item.children && hoveredItemId === item.id}
					<span transition:fade={{ duration: 150 }} class="tooltip">
						{item.label}
					</span>
				{/if}
			</div>
		{/each}
	</div>
</nav>

<!-- ================= Mobile Floating Capsule (顶部悬浮) ================= -->
<svelte:window bind:scrollY />

<div
	class="fixed z-50 lg:hidden flex justify-center transition-all ease-[cubic-bezier(0.23,1,0.32,1)]"
	class:duration-0={!isMobileMenuOpen && !isMenuAnimating}
	class:duration-500={isMobileMenuOpen || isMenuAnimating}
	class:top-0={isMobileMenuOpen}
	class:inset-x-0={isMobileMenuOpen}
	style:top={isMobileMenuOpen ? undefined : `${16 * (1 - navProgress)}px`}
	style:left={isMobileMenuOpen ? undefined : `${16 * (1 - navProgress)}px`}
	style:right={isMobileMenuOpen ? undefined : `${16 * (1 - navProgress)}px`}
>
	<div
		class="relative w-full mx-auto overflow-hidden transition-all ease-[cubic-bezier(0.23,1,0.32,1)]"
		class:duration-0={!isMobileMenuOpen && !isMenuAnimating}
		class:duration-500={isMobileMenuOpen || isMenuAnimating}
		class:shadow-glass-lg={isMobileMenuOpen}
		class:rounded-none={isMobileMenuOpen}
		class:h-screen={isMobileMenuOpen}
		style:border-radius={isMobileMenuOpen ? undefined : `${24 * (1 - navProgress)}px`}
	>
		<!-- 背景层 (变形动画) -->
		<div
			class="absolute inset-0 bg-white/90 dark:bg-ink-900/90 backdrop-blur-xl border-white/40 dark:border-ink-700 shadow-glass transition-all ease-[cubic-bezier(0.23,1,0.32,1)]"
			class:duration-0={!isMobileMenuOpen && !isMenuAnimating}
			class:duration-500={isMobileMenuOpen || isMenuAnimating}
			style:opacity={1}
			style:border-width="1px"
			style:height={isMobileMenuOpen ? '100vh' : '3rem'}
			style:min-height={isMobileMenuOpen ? '100vh' : '3rem'}
		></div>

		<!-- 1. 收起状态栏 (Collapsed Header) -->
		<div class="relative z-10 flex items-center justify-between px-3 h-12">
			<!-- 左侧信息 & 菜单开关 -->
			<div class="flex items-center gap-3">
				<button
					onclick={(e) => {
						e.stopPropagation();
						isMobileMenuOpen = !isMobileMenuOpen;
					}}
					class="w-9 h-9 flex items-center justify-center rounded-full active:scale-90 transition-transform"
				>
					<div
						class="w-8 h-8 rounded-full border border-ink-100 dark:border-ink-700 overflow-hidden shrink-0"
					>
						<img
							src="https://dogeoss.grtsinry43.com/img/author.jpeg"
							alt="A"
							class="w-full h-full object-cover"
						/>
					</div>
				</button>

				<div
					class="flex flex-col justify-center transition-all duration-300"
					class:opacity-0={isMobileMenuOpen}
				>
					<span
						class="text-sm font-serif font-bold text-ink-900 dark:text-jade-100 leading-none truncate max-w-[200px]"
						>{displayTitle}</span
					>
				</div>
			</div>

			<!-- 右侧按钮组 (TOC & Menu Icon for visual balance if needed, but we used Avatar as Menu Trigger above? Wait, let's keep Menu button)
           Alignment Fix: The previous logic had a separate Menu button. The user wants "vertical centering".
           Let's standardise: Left = Avatar/Title, Right = TOC + Menu Toggle.
      -->

			<div class="flex items-center gap-1">
				<!-- TOC Trigger -->
				<button
					onclick={(e) => {
						e.stopPropagation();
						isTocOpen = true;
					}}
					class="w-9 h-9 flex items-center justify-center rounded-full hover:bg-black/5 dark:hover:bg-white/10 transition-colors text-ink-600 dark:text-ink-300"
				>
					<List size={20} />
				</button>

				<!-- Main Menu Trigger -->
				<button
					onclick={(e) => {
						e.stopPropagation();
						isMobileMenuOpen = !isMobileMenuOpen;
					}}
					class="w-9 h-9 flex items-center justify-center rounded-full hover:bg-black/5 dark:hover:bg-white/10 transition-colors"
				>
					{#if isMobileMenuOpen}
						<X size={20} class="text-ink-600 dark:text-ink-300" />
					{:else}
						<!-- Use a generic menu icon or maybe a different one since Avatar can also be Home?
                 Let's stick to Menu icon for clarity as requested. -->
						<Menu size={20} class="text-ink-600 dark:text-ink-300" />
					{/if}
				</button>
			</div>
		</div>

		<!-- 2. 展开菜单内容 (Expanded Content) -->
		{#if isMobileMenuOpen}
			<div
				transition:slide={{ duration: 400, axis: 'y' }}
				class="relative z-10 px-2 pb-6 pt-0 flex flex-col max-h-[75vh] overflow-y-auto no-scrollbar"
			>
				<!-- 装饰标题 -->
				<div
					class="px-4 pb-4 pt-2 border-b border-ink-200/50 dark:border-ink-700/50 mb-3 flex items-center justify-between"
				>
					<span class="text-xs font-bold text-ink-400 uppercase tracking-widest">Navigation</span>
					<span class="text-[10px] text-ink-300 font-mono">MENU</span>
				</div>

				<div class="flex flex-col gap-1">
					{#each navItems as item}
						{@const isActive =
							activeSection === item.id || item.children?.some((c) => c.id === activeSection)}
						{@const hasChildren = !!item.children}
						{@const isExpanded = expandedMobileItems.has(item.id)}

						<div class="flex flex-col">
							<!-- 一级菜单项 -->
							<div
								class="group flex items-center gap-3 px-3 py-2 rounded-xl transition-all duration-300 relative overflow-hidden cursor-pointer select-none
                {isActive
									? 'bg-white dark:bg-ink-800'
									: 'hover:bg-white/50 dark:hover:bg-ink-800/50'}"
							>
								<!-- 激活边框 -->
								{#if isActive}
									<div
										class="absolute inset-0 border border-jade-200 dark:border-jade-800 rounded-xl pointer-events-none"
									></div>
								{/if}

								<button
									type="button"
									onclick={() => handleNavigate(item.id)}
									class="flex flex-1 min-w-0 items-center gap-3 text-left"
								>
									<!-- 图标 -->
									<div
										class="w-8 h-8 flex items-center justify-center rounded-full transition-colors duration-300 shrink-0
                  {isActive
											? 'bg-jade-100 text-jade-700 dark:bg-jade-900 dark:text-jade-300'
											: 'bg-ink-100 text-ink-500 dark:bg-ink-950 dark:text-ink-400'}"
									>
										<item.icon size={16} />
									</div>

									<!-- 文本 -->
									<div class="flex-1 min-w-0">
										<div
											class="font-serif text-[15px] font-medium truncate {isActive
												? 'text-jade-800 dark:text-jade-100'
												: 'text-ink-700 dark:text-ink-300'}"
										>
											{item.label}
										</div>
										<div
											class="text-[10px] text-ink-400 dark:text-ink-500 line-clamp-1 font-sans mt-0.5"
										>
											{item.desc}
										</div>
									</div>
								</button>

								<!-- 展开/收起按钮 (独立交互) -->
								{#if hasChildren}
									<button
										type="button"
										onclick={(e) => toggleMobileSubmenu(e, item.id)}
										class="p-2 -mr-2 rounded-full hover:bg-ink-100 dark:hover:bg-white/10 text-ink-400 transition-colors active:scale-90"
									>
										<ChevronDown
											size={16}
											class="transition-transform duration-300 {isExpanded
												? 'rotate-180 text-jade-600'
												: ''}"
										/>
									</button>
								{/if}
							</div>

							<!-- 二级菜单 (Tree Accordion) -->
							{#if hasChildren && isExpanded}
								<div
									transition:slide={{ duration: 300, easing: cubicOut }}
									class="flex flex-col gap-1 mt-1 mb-2 relative"
								>
									<!-- 垂直树线 -->
									<div
										class="absolute left-[39px] top-0 bottom-4 w-[1px] bg-ink-200 dark:bg-ink-700"
									></div>

									{#each item.children as subItem}
										{@const isSubActive = activeSection === subItem.id}
										<button
											onclick={() => handleNavigate(subItem.id)}
											class="relative flex items-center gap-3 pl-[54px] pr-4 py-2.5 rounded-lg ml-2 mr-2 transition-colors text-left group/sub
                       {isSubActive
												? 'bg-jade-50/50 dark:bg-jade-900/20'
												: 'hover:bg-white/60 dark:hover:bg-white/5'}"
										>
											<!-- 水平分支线 -->
											<div
												class="absolute left-[31px] top-1/2 w-4 h-[1px] bg-ink-200 dark:bg-ink-700"
											></div>

											{#if subItem.icon}
												<subItem.icon
													size={14}
													class="{isSubActive
														? 'text-jade-600 dark:text-jade-400'
														: 'text-ink-400'} transition-colors"
												/>
											{/if}
											<span
												class="text-sm font-medium {isSubActive
													? 'text-jade-700 dark:text-jade-300'
													: 'text-ink-600 dark:text-ink-400'}"
											>
												{subItem.label}
											</span>
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

	<!-- 全局遮罩 (点击外部关闭) -->
	{#if isMobileMenuOpen}
		<button
			type="button"
			aria-label="Close menu"
			transition:fade={{ duration: 300 }}
			class="fixed inset-0 bg-ink-900/20 backdrop-blur-[2px] -z-10"
			onclick={() => (isMobileMenuOpen = false)}
		></button>
	{/if}
</div>

<!-- ================= Mobile TOC Drawer ================= -->
<TableOfContents isOpen={isTocOpen} onClose={() => (isTocOpen = false)} />

<style>
	@reference "../../../app.css";

	/* --- 通用：隐藏滚动条 --- */
	.no-scrollbar::-webkit-scrollbar {
		display: none;
	}
	.no-scrollbar {
		-ms-overflow-style: none;
		scrollbar-width: none;
	}

	/* --- Desktop Styles --- */
	.glass-dock {
		@apply dark:bg-ink-950/40 bg-white/40;
		@apply border border-white/50 backdrop-blur-2xl dark:border-white/10;
		@apply shadow-glass;
		@apply rounded-full;
	}

	.glass-popover {
		@apply dark:bg-ink-900/60 bg-white/60;
		@apply border border-white/40 backdrop-blur-xl dark:border-white/10;
		@apply shadow-float rounded-2xl;
	}

	.nav-btn-desktop {
		@apply flex h-10 w-10 items-center justify-center rounded-full;
		@apply text-ink-500 dark:text-ink-400;
		@apply hover:bg-jade-50/50 dark:hover:bg-jade-900/20;
		@apply hover:text-jade-600 dark:hover:text-jade-400;
		@apply transition-all duration-300;
	}

	.nav-btn-desktop.active {
		@apply bg-jade-500 dark:bg-jade-500 text-white;
		@apply shadow-glow;
		@apply scale-110;
	}

	/* Tooltip */
	.tooltip {
		@apply bg-ink-900 dark:bg-ink-100 dark:text-ink-900 absolute left-14 text-white;
		@apply rounded-full px-3 py-1.5 text-[12px] whitespace-nowrap;
		@apply shadow-float pointer-events-none font-serif tracking-wide;
		@apply top-1/2 -translate-y-1/2;
	}

	/* --- Shadows --- */
	.shadow-glass {
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
	}
	.shadow-glass-lg {
		box-shadow: 0 15px 30px rgba(0, 0, 0, 0.08);
	}
	:global(.dark) .shadow-glass {
		box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
	}
</style>
