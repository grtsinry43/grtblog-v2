<script lang="ts">
	import type { PostSummary } from '$lib/domain/post';
	import Card from '$lib/ui/shared/Card.svelte';
	import Button from '$lib/ui/shared/Button.svelte';

	type Pagination = {
		total: number;
		page: number;
		size: number;
	};

	let { posts, pagination } = $props<{ posts: PostSummary[]; pagination: Pagination }>();

	let totalPages = $derived(
		pagination.size > 0 ? Math.max(1, Math.ceil(pagination.total / pagination.size)) : 1
	);
	let hasPrev = $derived(pagination.page > 1);
	let hasNext = $derived(pagination.page < totalPages);

	const buildPageLink = (page: number) => {
		const safePage = Number.isFinite(page) && page > 1 ? page : 1;
		return safePage === 1 ? '/posts/' : `/posts/page/${safePage}/`;
	};

	const formatDate = (value?: string) => {
		if (!value) return '';
		const date = new Date(value);
		return new Intl.DateTimeFormat('zh-CN', {
			year: 'numeric',
			month: 'short',
			day: '2-digit'
		}).format(date);
	};
</script>

<div class="space-y-8">
	{#if posts && posts.length > 0}
		<div class="grid gap-8 md:grid-cols-2 lg:grid-cols-2">
			{#each posts as post}
				<a href="/posts/{post.shortUrl}" class="block group h-full">
					<Card
						className="h-full flex flex-col justify-between hover:-translate-y-1 transition-all duration-500"
					>
						<div>
							<div class="flex items-center gap-3 mb-4">
								<span
									class="text-[10px] font-mono font-bold text-jade-600 dark:text-jade-400 uppercase tracking-widest px-2 py-0.5 bg-jade-50 dark:bg-jade-900/30 rounded-full"
									>Article</span
								>
								<span class="text-[10px] font-mono text-ink-400 dark:text-ink-500"
									>{formatDate(post.createdAt)}</span
								>
							</div>

							<h3
								class="font-serif text-2xl text-ink-900 dark:text-ink-100 mb-3 group-hover:text-jade-700 dark:group-hover:text-jade-400 transition-colors leading-tight"
							>
								{post.title || '未命名'}
							</h3>

							<p
								class="text-sm text-ink-500 dark:text-ink-400 line-clamp-3 leading-relaxed mb-6 font-light"
							>
								{post.summary || '暂无摘要...'}
							</p>
						</div>

						<div
							class="flex items-center text-xs font-mono text-jade-700 dark:text-jade-400 opacity-0 group-hover:opacity-100 transition-all duration-500 transform translate-x-[-10px] group-hover:translate-x-0"
						>
							READ MORE <span class="ml-2">→</span>
						</div>
					</Card>
				</a>
			{/each}
		</div>
	{:else}
		<div class="py-24 text-center">
			<p class="text-ink-400 font-serif italic text-lg">此处静谧无声...</p>
			<p class="text-xs text-ink-300 font-mono mt-2 uppercase tracking-widest">No articles found</p>
		</div>
	{/if}

	{#if totalPages > 1}
		<nav
			class="flex items-center justify-center gap-6 mt-16 py-8 border-t border-ink-200 dark:border-ink-800"
		>
			<Button
				variant="secondary"
				className={!hasPrev ? 'opacity-30 pointer-events-none' : ''}
				onclick={() => (window.location.href = buildPageLink(pagination.page - 1))}
			>
				<span class="mr-2">←</span> Previous
			</Button>

			<div class="text-xs font-mono text-ink-400 dark:text-ink-500">
				{pagination.page} / {totalPages}
			</div>

			<Button
				variant="secondary"
				className={!hasNext ? 'opacity-30 pointer-events-none' : ''}
				onclick={() => (window.location.href = buildPageLink(pagination.page + 1))}
			>
				Next <span class="ml-2">→</span>
			</Button>
		</nav>
	{/if}
</div>
