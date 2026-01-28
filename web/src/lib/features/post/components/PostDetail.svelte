<script lang="ts">
	import { onDestroy, tick } from 'svelte';
	import type { PostDetail } from '$lib/features/post/types';
	import { renderMarkdown } from '$lib/shared/markdown';
	import type { TOCNode } from '$lib/features/post/types';
	import { mountMarkdownComponents } from '$lib/shared/markdown/components';
	import { Calendar, Clock, Share2, ArrowLeft } from 'lucide-svelte';
	import Button from '$lib/components/Button.svelte';
	import Badge from '$lib/components/Badge.svelte';
	import Tag from '$lib/components/Tag.svelte';
	import Divider from '$lib/components/Divider.svelte';
	import '$lib/shared/markdown/components/register';

	let { post, updated = false } = $props<{ post: PostDetail | null; updated?: boolean }>();
	let contentRoot: HTMLElement | null = $state(null);
	let activeAnchor = $state<string | null>(null);
	let observer: IntersectionObserver | null = null;
	let cleanupComponents: (() => void) | null = null;

	const flattenTOC = (nodes?: TOCNode[]) => {
		if (!nodes?.length) return [];
		const anchors: string[] = [];
		const walk = (items: TOCNode[]) => {
			for (const item of items) {
				anchors.push(item.anchor);
				if (item.children?.length) walk(item.children);
			}
		};
		walk(nodes);
		return anchors;
	};

	let contentHtml = $derived(post ? renderMarkdown(post.content ?? '', flattenTOC(post.toc)) : '');

	const setupObserver = () => {
		if (!contentRoot || typeof IntersectionObserver === 'undefined') return;
		observer?.disconnect();
		const headings = contentRoot.querySelectorAll('h1, h2, h3, h4, h5, h6');
		if (!headings.length) {
			activeAnchor = null;
			return;
		}
		observer = new IntersectionObserver(
			(entries) => {
				const visible = entries.filter((entry) => entry.isIntersecting);
				if (!visible.length) return;
				visible.sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top);
				const target = visible[0]?.target as HTMLElement | undefined;
				if (target?.id) activeAnchor = target.id;
			},
			{ rootMargin: '0px 0px -70% 0px', threshold: 0 }
		);
		headings.forEach((heading) => observer?.observe(heading));
	};

	const refreshObserver = async () => {
		await tick();
		setupObserver();
		cleanupComponents?.();
		if (contentRoot) cleanupComponents = mountMarkdownComponents(contentRoot);
	};

	const scrollToAnchor = (anchor: string, event: MouseEvent) => {
		event.preventDefault();
		if (!contentRoot) return;
		const target = contentRoot.querySelector(`#${CSS.escape(anchor)}`) as HTMLElement | null;
		if (!target) return;
		target.scrollIntoView({ behavior: 'smooth', block: 'start' });
		activeAnchor = anchor;
		if (typeof history !== 'undefined') history.replaceState(null, '', `#${anchor}`);
	};

	$effect(() => {
		void refreshObserver();
	});

	onDestroy(() => {
		observer?.disconnect();
		cleanupComponents?.();
	});

	const formatDate = (value?: string) => {
		if (!value) return '';
		const date = new Date(value);
		if (Number.isNaN(date.getTime())) return value;
		return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`;
	};
</script>

{#if post}
	{#snippet backContent()}
		<ArrowLeft size={14} class="group-hover:-translate-x-1 transition-transform" />
		<span>返回</span>
	{/snippet}

	{#snippet shareContent()}
		<Share2 size={14} />
	{/snippet}

	{#snippet topContent()}
		返回顶部
	{/snippet}

	<article class="article-container">
		<!-- Header -->
		<header class="article-header">
			<div class="back-nav">
				<Button
					variant="ghost"
					class="back-btn group"
					onclick={() => history.back()}
					content={backContent}
				/>
				<div class="nav-divider"></div>
			</div>

			<div class="title-section">
				<div class="meta-badge-row">
					<Badge variant="soft">专题</Badge>
					<span class="meta-label">技术与设计</span>
				</div>

				<h1 class="article-title">{post.title}</h1>

				<div class="meta-info-row">
					<span class="meta-item"><Calendar size={12} /> {formatDate(post.createdAt)}</span>
					<span class="meta-item"><Clock size={12} /> 12 分钟阅读</span>
					{#if updated}
						<Badge variant="soft" class="animate-pulse">内容已更新</Badge>
					{/if}
				</div>
			</div>

			{#if post.leadIn}
				<p class="lead-in">{post.leadIn}</p>
			{/if}
		</header>

		<!-- Content Grid -->
		<div class="content-layout">
			<!-- Main Content -->
			<main class="main-content">
				<div
					class="markdown-preview markdown-body prose prose-ink dark:prose-invert"
					bind:this={contentRoot}
				>
					{@html contentHtml}
				</div>

				<footer class="article-footer">
					<div class="share-row">
						<span class="share-label">分享此文</span>
						<Button variant="ghost" class="share-btn" title="分享" content={shareContent} />
					</div>
					<Button
						variant="ghost"
						class="top-btn"
						onclick={() => window.scrollTo({ top: 0, behavior: 'smooth' })}
						content={topContent}
					/>
				</footer>
			</main>

			<!-- Sidebar / TOC -->
			{#if post.toc?.length}
				<aside class="desktop-sidebar">
					<div class="sidebar-sticky">
						<div class="toc-container">
							<span class="toc-header">本页目录</span>
							<ul class="toc-list">
								{#each post.toc as item}
									<li class="toc-item">
										<a
											class="toc-link {activeAnchor === item.anchor ? 'active' : ''}"
											href={'#' + item.anchor}
											onclick={(event) => scrollToAnchor(item.anchor, event)}
										>
											{item.name}
										</a>
										{#if item.children?.length}
											<ul class="toc-sublist">
												{#each item.children as child}
													<li>
														<a
															class="toc-sublink {activeAnchor === child.anchor ? 'active' : ''}"
															href={'#' + child.anchor}
															onclick={(event) => scrollToAnchor(child.anchor, event)}
														>
															{child.name}
														</a>
													</li>
												{/each}
											</ul>
										{/if}
									</li>
								{/each}
							</ul>
						</div>

						<div class="sidebar-card">
							<Tag variant="jade" class="border-none px-0">感悟</Tag>
							<p class="sidebar-card-text">
								每一篇文章都是漫长探索中的一小步。如果这些文字能引起共鸣，欢迎留下你的思考。
							</p>
						</div>
					</div>
				</aside>
			{/if}
		</div>
	</article>
{:else}
	<div class="empty-article">
		<p>请求的内容未能呈现。</p>
	</div>
{/if}

<style>
	@reference "../../../../routes/layout.css";

	.article-container {
		@apply article-enter space-y-10;
	}

	.article-header {
		@apply max-w-4xl space-y-6;
	}

	.back-nav {
		@apply flex items-center gap-4;
	}

	.back-btn {
		@apply !h-auto !p-0 font-mono text-[10px] font-semibold tracking-[0.2em] text-ink-400 uppercase hover:!bg-transparent hover:text-ink-900;
	}

	.nav-divider {
		@apply h-px w-6 bg-ink-200/50 dark:bg-ink-800/50;
	}

	.title-section {
		@apply space-y-4;
	}

	.meta-badge-row {
		@apply flex items-center gap-3;
	}

	.meta-label {
		@apply font-mono text-[9px] tracking-[0.3em] text-ink-400 uppercase;
	}

	.article-title {
		@apply font-serif text-2xl leading-[1.2] font-medium tracking-tight text-ink-950 md:text-3xl lg:text-4xl dark:text-ink-50;
	}

	.meta-info-row {
		@apply flex flex-wrap items-center gap-5 font-mono text-[9px] tracking-widest text-ink-400 uppercase;
	}

	.meta-item {
		@apply flex items-center gap-1.5;
	}

	.lead-in {
		@apply border-l-[1px] border-jade-500/20 py-0.5 pl-5 font-serif text-base leading-relaxed font-normal text-ink-600 italic opacity-90 md:text-lg dark:text-ink-400;
	}

	.content-layout {
		@apply grid gap-10 lg:grid-cols-[1fr_220px] lg:gap-16;
	}

	.main-content {
		@apply min-w-0;
	}

	.markdown-body {
		@apply max-w-none text-[15px] leading-[1.8] font-normal text-ink-800 md:text-base dark:text-ink-200;
	}

	.article-footer {
		@apply mt-16 flex flex-col items-start justify-between gap-4 border-t border-ink-50 pt-8 md:flex-row md:items-center dark:border-ink-800/30;
	}

	.share-row {
		@apply flex items-center gap-3;
	}

	.share-label {
		@apply font-mono text-[9px] tracking-widest text-ink-400 uppercase;
	}

	.share-btn {
		@apply h-auto p-1.5 text-ink-400 hover:text-jade-600;
	}

	.top-btn {
		@apply !h-auto !p-0 font-mono text-[9px] tracking-[0.2em] text-ink-400 uppercase hover:!bg-transparent hover:text-ink-900;
	}

	.desktop-sidebar {
		@apply hidden lg:block;
	}

	.sidebar-sticky {
		@apply sticky top-20 space-y-10;
	}

	.toc-container {
		@apply space-y-5;
	}

	.toc-header {
		@apply block border-b border-ink-50 pb-2 font-mono text-[8px] font-bold tracking-[0.4em] text-ink-300 uppercase dark:border-ink-800/30;
	}

	.toc-list {
		@apply space-y-3 font-sans;
	}

	.toc-item {
		@apply space-y-2;
	}

	.toc-link {
		@apply block text-[12px] text-ink-500 transition-all hover:translate-x-0.5 hover:text-jade-600 dark:text-ink-400 dark:hover:text-jade-400;
	}

	.toc-link.active {
		@apply font-bold text-jade-700 dark:text-jade-400;
	}

	.toc-sublist {
		@apply space-y-1.5 border-l border-ink-50 pl-3 dark:border-ink-800/30;
	}

	.toc-sublink {
		@apply block text-[11px] text-ink-400 transition-all hover:translate-x-0.5 hover:text-jade-500 dark:text-ink-500;
	}

	.toc-sublink.active {
		@apply font-bold text-jade-600 dark:text-jade-300;
	}

	.sidebar-card {
		@apply space-y-2.5 rounded-lg border border-jade-500/10 bg-jade-500/5 p-5;
	}

	.sidebar-card-text {
		@apply text-[10px] leading-relaxed font-normal text-ink-500 dark:text-ink-400;
	}

	.empty-article {
		@apply py-24 text-center font-serif text-sm text-ink-400 italic;
	}

	.article-enter {
		animation: article-enter 0.8s cubic-bezier(0.16, 1, 0.3, 1) both;
	}

	@keyframes article-enter {
		from {
			opacity: 0;
			transform: translateY(8px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	:global(.markdown-body h1, .markdown-body h2, .markdown-body h3) {
		@apply mt-10 mb-4 font-serif font-medium tracking-tight text-ink-950 dark:text-ink-50;
	}

	:global(.markdown-body h2) {
		@apply text-xl md:text-2xl;
	}
	:global(.markdown-body h3) {
		@apply text-lg md:text-xl;
	}

	:global(.markdown-body p) {
		@apply mb-6;
	}

	:global(.markdown-body blockquote) {
		@apply my-8 border-l-[1px] border-jade-500/40 py-0.5 pl-5 text-[0.95em] text-ink-600 italic opacity-90 dark:text-ink-400;
	}
</style>
