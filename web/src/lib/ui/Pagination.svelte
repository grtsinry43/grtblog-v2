<script lang="ts">
	import Button from './Button.svelte';
	import { ArrowLeft, ArrowRight } from 'lucide-svelte';

	interface Props {
		current: number;
		total: number;
		onPageChange: (page: number) => void;
		class?: string;
	}

	let { current, total, onPageChange, class: className = '' }: Props = $props();

	const pages = $derived(Array.from({ length: total }, (_, i) => i + 1));
</script>

<nav class="pagination {className}">
	<Button
		variant="ghost"
		disabled={current <= 1}
		onclick={() => onPageChange(current - 1)}
		class="nav-btn"
	>
		<ArrowLeft size={12} class="mr-1.5" />
		上一页
	</Button>

	<div class="page-numbers">
		{#each pages as page}
			{#if page === current}
				<span class="page-num active">
					{page}
				</span>
			{:else if page === 1 || page === total || (page >= current - 1 && page <= current + 1)}
				<Button variant="ghost" onclick={() => onPageChange(page)} class="page-num-btn">
					{page}
				</Button>
			{:else if (page === current - 2 && page > 1) || (page === current + 2 && page < total)}
				<span class="ellipsis">...</span>
			{/if}
		{/each}
	</div>

	<Button
		variant="ghost"
		disabled={current >= total}
		onclick={() => onPageChange(current + 1)}
		class="nav-btn"
	>
		下一页
		<ArrowRight size={12} class="ml-1.5" />
	</Button>
</nav>

<style>
	@reference "../../routes/layout.css";

	.pagination {
		@apply flex items-center justify-center gap-3;
	}

	.nav-btn {
		@apply h-8 !bg-transparent px-2 font-mono text-[10px] tracking-widest uppercase hover:!text-jade-600;
	}

	.page-numbers {
		@apply flex items-center gap-1.5;
	}

	.page-num {
		@apply flex h-7 w-7 items-center justify-center font-mono text-[11px] font-bold transition-all;
	}

	.page-num.active {
		@apply rounded-md bg-jade-800 text-white shadow-sm;
	}

	.page-num-btn {
		@apply h-7 w-7 rounded-md !p-0 font-mono text-[11px] text-ink-400 hover:text-ink-900;
	}

	.ellipsis {
		@apply px-0.5 font-mono text-[9px] text-ink-300 select-none;
	}
</style>
