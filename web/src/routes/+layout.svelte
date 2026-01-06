<script lang="ts">
    import './layout.css';
    import favicon from '$lib/assets/favicon.svg';
    import Sidebar from "../components/Sidebar.svelte";
    import { SectionId } from "../components/types";

    let {children} = $props();
</script>

<Sidebar activeSection={SectionId.Home} onNavigate={() => {}} />

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
