<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import { postDetailCtx } from '$routes/posts/[id]/post-detail-context';
	import { getPostDetail } from '$lib/features/post/api';

	const { updateModelData } = postDetailCtx.useModelActions();
	const shortUrlStore = postDetailCtx.selectModelData((data) => data?.shortUrl ?? '');

	const metricsQuery = createQuery(() => ({
		queryKey: ['post-metrics', $shortUrlStore],
		enabled: Boolean($shortUrlStore),
		queryFn: async () => {
			if (!$shortUrlStore) return null;
			const detail = await getPostDetail(undefined, $shortUrlStore);
			return detail?.metrics ?? null;
		}
	}));

	$effect(() => {
		if (!metricsQuery.data) return;
		updateModelData((prev) => {
			if (!prev) return prev;
			return {
				...prev,
				metrics: metricsQuery.data ?? prev.metrics
			};
		});
	});
</script>

<span>
	{#if metricsQuery.data}
		<span class="text-sm text-gray-500 dark:text-gray-400">
			阅读 {metricsQuery.data.views} · 喜欢 {metricsQuery.data.likes} · 评论 {metricsQuery.data.comments}
		</span>
	{/if}
</span>
