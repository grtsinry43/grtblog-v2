<script lang="ts">
	import type { PostDetail } from '$lib/models/post';

	let { post } = $props<{ post: PostDetail | null }>();
</script>

{#if post}
	<h1>{post.title}</h1>
	<p>{post.content}</p>
{:else}
	<p>Article not found.</p>
{/if}
