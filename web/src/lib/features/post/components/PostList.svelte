<script lang="ts">
	import type { PostSummary } from '$lib/features/post/types';
	import Card from '$lib/components/Card.svelte';
	import Pagination from '$lib/components/Pagination.svelte';
	import Tag from '$lib/components/Tag.svelte';
	import { ArrowRight, FileText } from 'lucide-svelte';

	type PaginationData = {
		total: number;
		page: number;
		size: number;
	};

	let { posts, pagination } = $props<{ posts: PostSummary[]; pagination: PaginationData }>();

	let totalPages = $derived(
		pagination.size > 0 ? Math.max(1, Math.ceil(pagination.total / pagination.size)) : 1
	);

	const onPageChange = (page: number) => {
		const safePage = Number.isFinite(page) && page > 1 ? page : 1;
		window.location.href = safePage === 1 ? '/posts/' : `/posts/page/${safePage}/`;
	};

	const formatDate = (dateStr: string) => {
		const date = new Date(dateStr);
		return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`;
	};
</script>

<div class="archive-container">
	<header class="archive-header">
		<h1 class="archive-title">文章归档</h1>
		<p class="archive-desc">按时间顺序排布的思考、笔记与技术沉淀。</p>
	</header>

	{#if posts && posts.length > 0}
		<div class="archive-list">
			{#each posts as post}
				<a href="/posts/{post.shortUrl}" class="archive-link group">
					<Card variant="seamless" hover={true} class="archive-card">
						<div class="archive-content">
							<div class="archive-meta">
								{#if post.createdAt}
									<span>{formatDate(post.createdAt)}</span>
								{/if}
								<span class="meta-sep">/</span>
								<Tag variant="jade" class="border-none px-0 tracking-widest !h-auto !text-[9px]"
									>专栏</Tag
								>
							</div>

							<h3 class="post-title">
								{post.title || '无标题文章'}
							</h3>

							<p class="post-summary">
								{post.summary || '这篇文章还没有摘要。'}
							</p>
						</div>

						<div class="arrow-box">
							<ArrowRight size={14} />
						</div>
					</Card>
				</a>
			{/each}
		</div>

		{#if totalPages > 1}
			<div class="pagination-wrapper">
				<Pagination current={pagination.page} total={totalPages} {onPageChange} />
			</div>
		{/if}
	{:else}
		<div class="empty-state">
			<FileText size={24} class="mx-auto text-ink-200 dark:text-ink-800" />
			<h3 class="empty-title">未发现内容</h3>
			<p class="empty-desc">文章归档列表目前为空。</p>
		</div>
	{/if}
</div>

<style>
	@reference "../../../../routes/layout.css";

	.archive-container {
		@apply space-y-10;
	}

	.archive-header {
		@apply space-y-3;
	}

	.archive-title {
		@apply font-serif text-2xl font-medium tracking-tight text-ink-950 md:text-3xl dark:text-ink-50;
	}

	.archive-desc {
		@apply max-w-lg text-[13px] leading-relaxed font-normal text-ink-500 opacity-80 dark:text-ink-400;
	}

	.archive-list {
		@apply grid gap-px border-y border-ink-50 bg-ink-50 dark:border-ink-900/50 dark:bg-ink-900/50;
	}

	.archive-link {
		@apply block;
	}

	.archive-card {
		@apply flex flex-col gap-5 bg-white p-5 md:flex-row md:items-center md:p-6 dark:bg-ink-950;
	}

	.archive-content {
		@apply flex-1 space-y-1.5;
	}

	.archive-meta {
		@apply flex items-center gap-2.5 font-mono text-[9px] tracking-widest text-ink-400 uppercase;
	}

	.meta-sep {
		@apply text-ink-100 dark:text-ink-800/50;
	}

	.post-title {
		@apply font-serif text-lg leading-tight font-medium text-ink-900 transition-colors group-hover:text-jade-800 md:text-xl dark:text-ink-100 dark:group-hover:text-jade-400;
	}

	.post-summary {
		@apply max-w-2xl text-[11px] leading-relaxed font-normal text-ink-500 opacity-80 md:text-xs dark:text-ink-400;
	}

	.arrow-box {
		@apply flex h-8 w-8 shrink-0 items-center justify-center rounded-full border border-ink-100 text-ink-400 transition-all duration-300 group-hover:border-jade-500 group-hover:bg-jade-500 group-hover:text-white dark:border-ink-800;
	}

	.pagination-wrapper {
		@apply pt-6;
	}

	.empty-state {
		@apply space-y-3 py-20 text-center;
	}

	.empty-title {
		@apply font-serif text-base text-ink-900 dark:text-ink-100;
	}

	.empty-desc {
		@apply text-xs font-normal text-ink-500 opacity-70;
	}
</style>
