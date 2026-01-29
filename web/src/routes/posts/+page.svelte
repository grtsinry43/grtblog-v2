<script lang="ts">
	import { PostList } from '$lib/features/post';
    import {postContext} from "$routes/posts/post-context";

    let { data } = $props();

    const postStore = postContext.mountModelData(data);

    $effect(() => {
        postContext.syncModelData(postStore, data);
    });
</script>

<PostList />
