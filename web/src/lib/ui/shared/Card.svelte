<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		children: Snippet;
		className?: string;
		padding?: string;
	}

	let { children, className = '', padding = 'p-6' }: Props = $props();
</script>

<div
	class="glass-panel rounded-3xl {padding} {className} hover:shadow-float transition-shadow duration-500"
>
	{@render children()}
</div>
