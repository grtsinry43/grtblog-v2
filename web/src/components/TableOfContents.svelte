<script lang="ts">
	import { X, Hash, ChevronRight } from 'lucide-svelte';
	import { fade, fly } from 'svelte/transition';
	import { cubicOut } from 'svelte/easing';

	let { isOpen, onClose } = $props<{
		isOpen: boolean;
		onClose: () => void;
	}>();

	// Mock Data for Visual Design
	const tocItems = [
		{ level: 1, text: 'Introduction', id: 'intro' },
		{ level: 2, text: 'why Svelte?', id: 'why' },
		{ level: 2, text: 'The Compiler Approach', id: 'compiler' },
		{ level: 1, text: 'Core Concepts', id: 'core' },
		{ level: 2, text: 'Reactivity', id: 'reactivity' },
		{ level: 3, text: 'Runes', id: 'runes' },
		{ level: 2, text: 'Templating', id: 'templating' },
		{ level: 1, text: 'Conclusion', id: 'conclusion' }
	];
</script>

{#if isOpen}
	<!-- Backdrop -->
	<div
		class="fixed inset-0 z-[60] bg-ink-900/20 backdrop-blur-sm transition-opacity"
		transition:fade={{ duration: 300 }}
		onclick={onClose}
	></div>

	<!-- Sidebar Drawer -->
	<div
		class="fixed right-0 top-0 bottom-0 z-[61] w-72 bg-white/90 dark:bg-ink-900/90 backdrop-blur-2xl border-l border-white/50 dark:border-ink-700 shadow-2xl flex flex-col"
		transition:fly={{ x: 300, duration: 400, easing: cubicOut, opacity: 1 }}
	>
		<!-- Header -->
		<div
			class="flex items-center justify-between px-5 h-16 border-b border-ink-100 dark:border-ink-800/50 shrink-0"
		>
			<div class="flex items-center gap-2 text-ink-800 dark:text-jade-100 font-serif font-bold">
				<Hash size={18} class="text-jade-600 dark:text-jade-400" />
				<span>Contents</span>
			</div>
			<button
				onclick={onClose}
				class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-ink-100 dark:hover:bg-white/10 text-ink-500 transition-colors"
			>
				<X size={18} />
			</button>
		</div>

		<!-- Scrollable Content -->
		<div class="flex-1 overflow-y-auto p-4 flex flex-col gap-1">
			{#each tocItems as item}
				<a
					href="#{item.id}"
					onclick={(e) => {
						e.preventDefault();
						onClose();
					}}
					class="group flex items-center gap-2 py-2 px-3 rounded-lg text-sm transition-all hover:bg-jade-50 dark:hover:bg-white/5"
					style:padding-left="{item.level * 12 + 4}px"
				>
					{#if item.level === 1}
						<ChevronRight
							size={14}
							class="text-jade-500 opacity-0 group-hover:opacity-100 transition-opacity -ml-4"
						/>
					{/if}
					<span
						class="truncate {item.level === 1
							? 'font-medium text-ink-900 dark:text-ink-200'
							: 'text-ink-500 dark:text-ink-400'}"
					>
						{item.text}
					</span>
				</a>
			{/each}
		</div>

		<!-- Footer Decoration -->
		<div class="p-4 border-t border-ink-100 dark:border-ink-800/50 text-center">
			<div class="text-[10px] text-ink-300 uppercase tracking-widest">On this page</div>
		</div>
	</div>
{/if}
