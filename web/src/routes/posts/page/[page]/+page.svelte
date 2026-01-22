<script lang="ts">
	import { PostList } from '$lib/ui/post';

	let { data } = $props();
</script>

<PostList posts={data.posts ?? []} pagination={data.pagination} />
