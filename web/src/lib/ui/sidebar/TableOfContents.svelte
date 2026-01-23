<script lang="ts">
	import { X, Hash, ChevronRight } from 'lucide-svelte';
	import { fade, fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';
	import type { TOCNode } from '$lib/features/post/types';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
		toc?: TOCNode[];
		activeAnchor?: string | null;
		onAnchorClick?: (anchor: string, e: MouseEvent) => void;
	}

	let { isOpen, onClose, toc = [], activeAnchor, onAnchorClick }: Props = $props();

	// Flatten TOC for mobile simpler display if needed, or keep nested
	const flattenedToc = $derived.by(() => {
		const result: { level: number; name: string; anchor: string }[] = [];
		const walk = (nodes: TOCNode[], level: number) => {
			for (const node of nodes) {
				result.push({ level, name: node.name, anchor: node.anchor });
				if (node.children) walk(node.children, level + 1);
			}
		};
		walk(toc, 1);
		return result;
	});
</script>

{#if isOpen}
	<!-- Backdrop -->
	<button
		type="button"
		class="backdrop"
		transition:fade={{ duration: 300 }}
		onclick={onClose}
		aria-label="关闭目录"
	></button>

	<!-- Sidebar Drawer -->
	<div class="drawer" transition:fly={{ x: 280, duration: 400, easing: cubicOut, opacity: 1 }}>
		<!-- Header -->
		<header class="drawer-header">
			<div class="header-content">
				<Hash size={16} class="text-jade-600 dark:text-jade-400" />
				<span>本页目录</span>
			</div>
			<button onclick={onClose} class="close-btn">
				<X size={16} />
			</button>
		</header>

		<!-- Scrollable Content -->
		<nav class="drawer-content no-scrollbar">
			{#if flattenedToc.length > 0}
				{#each flattenedToc as item}
					<a
						href="#{item.anchor}"
						class="toc-link {activeAnchor === item.anchor ? 'active' : ''}"
						style:padding-left="{item.level * 12 + 8}px"
						onclick={(e) => {
							if (onAnchorClick) onAnchorClick(item.anchor, e);
							onClose();
						}}
					>
						{#if item.level === 1 && activeAnchor === item.anchor}
							<div class="active-indicator" transition:fade></div>
						{/if}
						<span class="truncate">{item.name}</span>
					</a>
				{/each}
			{:else}
				<div class="empty-state">未发现目录内容</div>
			{/if}
		</nav>

		<!-- Footer Decoration -->
		<footer class="drawer-footer">
			<span class="footer-label">匠心营造 · 秩序与呼吸</span>
		</footer>
	</div>
{/if}

<style>
	@reference "../../../routes/layout.css";

	.backdrop {
		@apply fixed inset-0 z-[100] bg-ink-950/20 backdrop-blur-[2px];
	}

	.drawer {
		@apply fixed top-0 right-0 bottom-0 z-[101] w-64 bg-white/95 dark:bg-ink-900/95;
		@apply flex flex-col border-l border-ink-100 shadow-2xl backdrop-blur-xl dark:border-ink-800/50;
	}

	.drawer-header {
		@apply flex h-14 shrink-0 items-center justify-between border-b border-ink-50 px-5 dark:border-ink-800/30;
	}

	.header-content {
		@apply flex items-center gap-2 font-serif text-sm font-bold text-ink-900 dark:text-jade-100;
	}

	.close-btn {
		@apply flex h-8 w-8 items-center justify-center rounded-full text-ink-400 transition-colors hover:bg-ink-50 dark:hover:bg-white/5;
	}

	.drawer-content {
		@apply flex flex-1 flex-col gap-0.5 overflow-y-auto py-4;
	}

	.toc-link {
		@apply relative flex items-center px-4 py-2 text-[13px] text-ink-500 transition-all hover:bg-jade-50/50 hover:text-jade-600 dark:text-ink-400 dark:hover:bg-white/5 dark:hover:text-jade-400;
	}

	.toc-link.active {
		@apply bg-jade-50/80 font-bold text-jade-700 dark:bg-jade-900/20 dark:text-jade-400;
	}

	.active-indicator {
		@apply absolute top-1/2 left-0 h-4 w-1 -translate-y-1/2 rounded-r-full bg-jade-600 dark:bg-jade-500;
	}

	.drawer-footer {
		@apply border-t border-ink-50 p-4 text-center dark:border-ink-800/30;
	}

	.footer-label {
		@apply font-mono text-[9px] tracking-[0.2em] text-ink-300 uppercase dark:text-ink-500;
	}

	.empty-state {
		@apply p-8 text-center text-xs font-light text-ink-300 italic;
	}
</style>
