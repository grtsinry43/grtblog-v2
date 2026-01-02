<script>
        let {data} = $props();
</script>

<h1>Articles</h1>

{#if data.articles && data.articles.length > 0}
    <ul>
        {#each data.articles as article}
            <li>
                <a href="/posts/{article.shortUrl}">{article.title}</a>
            </li>
        {/each}
    </ul>
{:else}
    <p>No articles found.</p>
{/if}