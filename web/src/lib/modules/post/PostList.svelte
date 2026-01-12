<script lang="ts">
	import type { PostSummary } from '$lib/models/post';

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

	const buildPageLink = (page: number) => `?page=${page}&pageSize=${pagination.size}`;
</script>

<h1>Articles</h1>

{#if posts && posts.length > 0}
	<ul>
		{#each posts as post}
			<li>
				<div class="p-2">
					<div>
						<a href="/posts/{post.shortUrl}">{post.title === "" ? "未命名": post.title}</a>
					</div>
					<a href="/posts/{post.shortUrl}">{post.summary}</a>
				</div>
			</li>
		{/each}
	</ul>
{:else}
	<p>No articles found.</p>
{/if}

<nav class="pagination">
	<a
		class="page-btn {hasPrev ? '' : 'is-disabled'}"
		href={buildPageLink(Math.max(1, pagination.page - 1))}
		aria-disabled={!hasPrev}
	>
		Prev
	</a>
	<span class="page-info">Page {pagination.page} / {totalPages}</span>
	<a
		class="page-btn {hasNext ? '' : 'is-disabled'}"
		href={buildPageLink(Math.min(totalPages, pagination.page + 1))}
		aria-disabled={!hasNext}
	>
		Next
	</a>
</nav>

<style>
	.pagination {
		margin-top: 16px;
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.page-btn {
		border: 1px solid #ddd;
		border-radius: 6px;
		padding: 6px 10px;
		background: #fff;
	}

	.page-btn.is-disabled {
		opacity: 0.5;
		pointer-events: none;
	}

	.page-info {
		font-size: 14px;
		color: #666;
	}
</style>
