<script lang="ts">
import { onDestroy, tick } from 'svelte';
import type { PostDetail } from '$lib/models/post';
import { renderMarkdown } from '$lib/shared/markdown';
import type { TOCNode } from '$lib/models/post';
import { mountMarkdownComponents } from '$lib/markdown/components';
import '$lib/markdown/components/register';

	let { post, updated = false } = $props<{ post: PostDetail | null; updated?: boolean }>();
	let contentRoot: HTMLElement | null = $state(null);
let activeAnchor = $state<string | null>(null);
let observer: IntersectionObserver | null = null;
let cleanupComponents: (() => void) | null = null;

	const flattenTOC = (nodes?: TOCNode[]) => {
		if (!nodes?.length) {
			return [];
		}
		const anchors: string[] = [];
		const walk = (items: TOCNode[]) => {
			for (const item of items) {
				anchors.push(item.anchor);
				if (item.children?.length) {
					walk(item.children);
				}
			}
		};
		walk(nodes);
		return anchors;
	};

	let contentHtml = $derived(post ? renderMarkdown(post.content ?? '', flattenTOC(post.toc)) : '');

	const setupObserver = () => {
		if (!contentRoot || typeof IntersectionObserver === 'undefined') {
			return;
		}
		observer?.disconnect();
		const headings = contentRoot.querySelectorAll('h1, h2, h3, h4, h5, h6');
		if (!headings.length) {
			activeAnchor = null;
			return;
		}
		observer = new IntersectionObserver(
			(entries) => {
				const visible = entries.filter((entry) => entry.isIntersecting);
				if (!visible.length) {
					return;
				}
				visible.sort((a, b) => a.boundingClientRect.top - b.boundingClientRect.top);
				const target = visible[0]?.target as HTMLElement | undefined;
				if (target?.id) {
					activeAnchor = target.id;
				}
			},
			{ rootMargin: '0px 0px -70% 0px', threshold: 0 }
		);
		headings.forEach((heading) => observer?.observe(heading));
	};

const refreshObserver = async () => {
	await tick();
	setupObserver();
	cleanupComponents?.();
	if (contentRoot) {
		cleanupComponents = mountMarkdownComponents(contentRoot);
	}
};

	const scrollToAnchor = (anchor: string, event: MouseEvent) => {
		event.preventDefault();
		if (!contentRoot) {
			return;
		}
		const target = contentRoot.querySelector(`#${CSS.escape(anchor)}`) as HTMLElement | null;
		if (!target) {
			return;
		}
		target.scrollIntoView({ behavior: 'smooth', block: 'start' });
		activeAnchor = anchor;
		if (typeof history !== 'undefined') {
			history.replaceState(null, '', `#${anchor}`);
		}
	};

	$effect(() => {
		void refreshObserver();
	});

onDestroy(() => {
	observer?.disconnect();
	cleanupComponents?.();
});

	const formatDate = (value?: string) => {
		if (!value) {
			return '';
		}
		const date = new Date(value);
		if (Number.isNaN(date.getTime())) {
			return value;
		}
		return new Intl.DateTimeFormat('zh-CN', {
			year: 'numeric',
			month: 'short',
			day: '2-digit'
		}).format(date);
	};
</script>

{#if post}
	<div class="relative px-6 py-10 md:px-10 md:py-14">
		<div class="absolute inset-0 -z-10 bg-[radial-gradient(circle_at_top,_rgba(45,212,191,0.15),_transparent_55%)]"></div>
		<article class="mx-auto max-w-5xl article-enter">
			<header class="space-y-5">
				<p class="text-xs uppercase tracking-[0.3em] text-ink-400 font-mono">Article</p>
				<h1 class="font-serif text-4xl leading-tight text-ink-900 md:text-5xl">
					{post.title}
				</h1>
				<div class="flex flex-wrap items-center gap-3 text-xs text-ink-500">
					<span>Published {formatDate(post.createdAt)}</span>
					<span class="h-1 w-1 rounded-full bg-ink-300"></span>
					<span>Updated {formatDate(post.updatedAt)}</span>
				</div>
				<div class="update-hint {updated ? 'update-hint--show' : ''}">
					<span class="update-dot"></span>
					内容已更新
				</div>
				{#if post.leadIn}
					<p class="rounded-3xl border border-ink-200 bg-white/70 px-6 py-5 text-lg leading-relaxed text-ink-700 shadow-subtle backdrop-blur">
						{post.leadIn}
					</p>
				{/if}
			</header>

			<div class="mt-10 grid gap-10 lg:grid-cols-[minmax(0,1fr)_280px]">
				<section class="space-y-8">
					<div class="rounded-[28px] border border-ink-200/70 bg-white/80 p-6 shadow-float backdrop-blur">
						<div class="text-sm font-semibold uppercase tracking-[0.25em] text-ink-400">Content</div>
						<div
							class="markdown-preview markdown-body prose mt-5 max-w-none text-[17px] leading-8 text-ink-800"
							bind:this={contentRoot}
						>
							{@html contentHtml}
						</div>
					</div>
				</section>

				{#if post.toc?.length}
					<aside class="lg:sticky lg:top-32 lg:max-h-[calc(100vh-10rem)] lg:overflow-auto">
						<div class="rounded-3xl border border-ink-200/80 bg-white/90 p-5 shadow-subtle backdrop-blur">
							<div class="flex items-center justify-between">
								<span class="text-xs font-semibold uppercase tracking-[0.25em] text-ink-400">
									Contents
								</span>
								<span class="text-[11px] text-ink-400 font-mono">TOC</span>
							</div>
							<ul class="mt-4 space-y-3 text-sm text-ink-700">
								{#each post.toc as item}
									<li class="space-y-2">
										<a
											class="transition-colors hover:text-ink-900 {activeAnchor === item.anchor ? 'font-semibold text-ink-900' : ''}"
											href={"#" + item.anchor}
											onclick={(event) => scrollToAnchor(item.anchor, event)}
										>
											{item.name}
										</a>
										{#if item.children?.length}
											<ul class="space-y-1 pl-4 text-xs text-ink-500">
												{#each item.children as child}
													<li>
														<a
															class="transition-colors hover:text-ink-800 {activeAnchor === child.anchor ? 'font-semibold text-ink-800' : ''}"
															href={"#" + child.anchor}
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
					</aside>
				{/if}
			</div>
		</article>
	</div>
{:else}
	<p>Article not found.</p>
{/if}

<style>
	.article-enter {
		animation: article-enter 0.7s cubic-bezier(0.16, 1, 0.3, 1) both;
	}

	.update-hint {
		display: inline-flex;
		align-items: center;
		gap: 0.5rem;
		padding: 0.4rem 0.75rem;
		border-radius: 999px;
		border: 1px solid rgba(45, 212, 191, 0.35);
		background: rgba(240, 253, 250, 0.85);
		color: #0f766e;
		font-size: 0.75rem;
		font-weight: 600;
		text-transform: uppercase;
		letter-spacing: 0.12em;
		box-shadow: 0 8px 20px -10px rgba(13, 148, 136, 0.35);
		opacity: 0;
		transform: translateY(-6px);
		transition: opacity 0.3s ease, transform 0.3s ease;
		pointer-events: none;
	}

	.update-hint--show {
		opacity: 1;
		transform: translateY(0);
	}

	.update-dot {
		width: 6px;
		height: 6px;
		border-radius: 999px;
		background: #14b8a6;
		box-shadow: 0 0 0 4px rgba(20, 184, 166, 0.25);
		animation: update-pulse 1.6s ease-in-out infinite;
	}

	@keyframes article-enter {
		from {
			opacity: 0;
			transform: translateY(18px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	@keyframes update-pulse {
		0%,
		100% {
			transform: scale(1);
			box-shadow: 0 0 0 4px rgba(20, 184, 166, 0.25);
		}
		50% {
			transform: scale(1.25);
			box-shadow: 0 0 0 8px rgba(20, 184, 166, 0.1);
		}
	}
</style>
