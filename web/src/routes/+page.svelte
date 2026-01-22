<script lang="ts">
	import { onMount } from 'svelte';
	import { fade, fly } from 'svelte/transition';
	import PostList from '$lib/ui/post/PostList.svelte';
	import Card from '$lib/ui/shared/Card.svelte';
	import Button from '$lib/ui/shared/Button.svelte';

	let mounted = $state(false);

	onMount(() => {
		mounted = true;
	});

	import type { PageData } from './$types';
	let { data }: { data: PageData } = $props();
</script>

<div class="max-w-5xl mx-auto px-6 md:px-12 py-12 md:py-24 relative">
	{#if mounted}
		<!-- Hero Section -->
		<section
			class="min-h-[60vh] flex flex-col justify-center relative mb-32"
			in:fade={{ duration: 1000 }}
		>
			<div class="relative z-10">
				<div
					class="inline-block px-3 py-1 mb-6 border border-jade-800/30 dark:border-jade-400/30 rounded-full"
					in:fly={{ y: 20, duration: 800, delay: 200 }}
				>
					<span class="text-xs font-mono text-jade-900 dark:text-jade-400 tracking-widest uppercase"
						>Grtblog V2</span
					>
				</div>

				<h1
					class="text-6xl md:text-8xl font-serif font-medium text-ink-900 dark:text-ink-100 mb-8 leading-[1.1] tracking-tight"
					in:fly={{ y: 30, duration: 800, delay: 400 }}
				>
					序 <span class="text-jade-800 dark:text-jade-400 italic">Xù</span>
					<br />
					<span
						class="text-4xl md:text-6xl text-ink-400 dark:text-ink-500 font-sans font-light tracking-normal"
					>
						秩序与呼吸。
					</span>
				</h1>

				<p
					class="text-lg text-ink-600 dark:text-ink-300 max-w-xl leading-relaxed font-light mb-10"
					in:fly={{ y: 30, duration: 800, delay: 600 }}
				>
					一套探索东方美学与现代界面平衡的设计语言。<br />
					在这里，文字与代码在纸墨的温润中寻找共鸣。
				</p>

				<div class="flex gap-4" in:fly={{ y: 30, duration: 800, delay: 800 }}>
					<Button
						variant="primary"
						onclick={() => window.scrollTo({ top: window.innerHeight * 0.8, behavior: 'smooth' })}
					>
						开始阅读
					</Button>
					<Button variant="secondary" onclick={() => (window.location.href = '/about')}>
						关于作者
					</Button>
				</div>
			</div>

			<!-- Decorative Hero Element -->
			<div
				class="absolute top-0 right-0 md:-right-20 w-[30rem] h-[30rem] bg-gradient-radial from-jade-200/30 to-transparent dark:from-jade-900/10 opacity-60 blur-3xl pointer-events-none -z-10"
			></div>
		</section>

		<!-- Recent Posts Section -->
		<section class="relative" in:fade={{ duration: 1000, delay: 1000 }}>
			<div class="flex items-center justify-between mb-12">
				<h2 class="text-2xl font-serif font-medium text-ink-900 dark:text-ink-100 italic">
					近期发布 <span class="text-jade-600 dark:text-jade-400 not-italic ml-2 text-lg"
						>/ Recent Posts</span
					>
				</h2>
				<div class="h-[1px] flex-1 bg-ink-200 dark:bg-ink-800 mx-8"></div>
				<a
					href="/posts"
					class="text-sm font-mono text-ink-400 hover:text-jade-600 dark:hover:text-jade-400 transition-colors uppercase tracking-widest"
					>View All</a
				>
			</div>

			<PostList
				posts={data?.posts || []}
				pagination={{
					total: (data?.posts || []).length,
					page: 1,
					size: (data?.posts || []).length
				}}
			/>
		</section>

		<!-- Footer -->
		<footer
			class="mt-48 pt-12 border-t border-ink-200 dark:border-ink-800 flex flex-col md:flex-row justify-between items-start md:items-center gap-6 pb-12"
		>
			<div>
				<h4 class="font-serif text-xl text-ink-900 dark:text-ink-100 font-bold mb-2">序 Xù</h4>
				<p class="text-xs text-ink-400 dark:text-ink-600">
					© 2026 Grtblog. Powered by SvelteKit & Oriental Aesthetics.
				</p>
			</div>
			<div class="flex gap-8">
				<div class="w-2 h-2 bg-jade-500 dark:bg-jade-400 rounded-full animate-pulse"></div>
				<span class="text-xs font-mono text-ink-300 dark:text-ink-600 uppercase tracking-widest"
					>Status: Flowing</span
				>
			</div>
		</footer>
	{/if}
</div>

<style>
	:global(body) {
		overflow-x: hidden;
	}
</style>
