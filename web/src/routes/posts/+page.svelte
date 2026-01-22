<script lang="ts">
	import { PostList } from '$lib/ui/post';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<PostList posts={data.posts ?? []} pagination={data.pagination} />
