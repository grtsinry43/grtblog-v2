<script lang="ts">
	import { PostList } from '$lib/modules/post';

	let { data } = $props();
</script>

<PostList posts={data.posts ?? []} pagination={data.pagination} />
