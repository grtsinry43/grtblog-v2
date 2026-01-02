<script>
    let { data } = $props();
</script>

{#if data.article}
    <h1>{data.article.title}</h1>
    <p>{data.article.content}</p>
{:else}
    <p>Article not found.</p>
{/if}