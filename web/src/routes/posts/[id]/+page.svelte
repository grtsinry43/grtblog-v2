<script lang="ts">
import { PostDetail } from '$lib/modules/post';
let { data } = $props();
</script>

<PostDetail post={data.post ?? null} />
