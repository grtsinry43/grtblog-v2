<script lang="ts">
    import './layout.css';
    import favicon from '$lib/assets/favicon.svg';

    let {children} = $props();
</script>

<svelte:head>
    <link rel="icon" href={favicon}/>
</svelte:head>

<main class="main">
    {@render children()}
</main>

<style>
    main {
        padding: 20px;
    }
</style>
