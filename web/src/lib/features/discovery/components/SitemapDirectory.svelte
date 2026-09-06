<script lang="ts">
	import { resolveHref } from '$lib/shared/utils/resolve-path';
	import type { DiscoveryCatalog } from '../types';
	let { catalog }: { catalog: DiscoveryCatalog } = $props();
	let query = $state('');
	const sections = [
		{ kind: 'navigation', title: '从这里出发', description: '找到适合此刻的阅读方式。' },
		{ kind: 'posts', title: '文章', description: '完整的写作、经验与探索。' },
		{ kind: 'moments', title: '手记', description: '生活片段与阶段性的记录。' },
		{ kind: 'pages', title: '独立页面', description: '关于作者，也关于这个博客。' },
		{ kind: 'categories', title: '分类', description: '沿着主题继续阅读。' },
		{ kind: 'columns', title: '专栏', description: '把相关的手记串在一起。' },
		{ kind: 'albums', title: '相册', description: '用影像留下记忆。' },
		{ kind: 'photos', title: '照片', description: '相册中的每一个瞬间。' }
	];
	const groups = $derived(
		sections
			.map((section) => ({
				...section,
				entries: catalog.entries.filter(
					(entry) =>
						entry.kind === section.kind &&
						`${entry.title} ${entry.summary ?? ''}`
							.toLocaleLowerCase()
							.includes(query.trim().toLocaleLowerCase())
				)
			}))
			.filter((section) => section.entries.length > 0)
	);
	const visibleCount = $derived(groups.reduce((count, group) => count + group.entries.length, 0));
</script>

<header class="mb-12 border-b border-ink-200 pb-10 dark:border-ink-800">
	<p class="mb-4 font-mono text-xs tracking-widest text-jade-600 dark:text-jade-400">
		{catalog.siteName}
	</p>
	<h1 class="font-serif text-4xl text-ink-900 sm:text-5xl dark:text-ink-100">站点地图</h1>
	<p class="mt-5 max-w-xl leading-relaxed text-ink-500">
		一份持续更新的阅读目录。你可以从一篇文章出发，也可以沿着分类，寻找感兴趣的内容。
	</p>
	<nav
		aria-label="其他阅读格式"
		class="mt-6 flex flex-wrap gap-5 text-sm text-jade-700 dark:text-jade-400"
	>
		<a href="/sitemap.xml" data-sveltekit-reload class="underline-offset-4 hover:underline"
			>XML 地图 ↗</a
		>
		<a href="/llms.txt" data-sveltekit-reload class="underline-offset-4 hover:underline"
			>AI 阅读指南 ↗</a
		>
		<a href="/feed" data-sveltekit-reload class="underline-offset-4 hover:underline">RSS 订阅 ↗</a>
	</nav>
</header>

<div class="mb-10 flex flex-wrap items-end justify-between gap-4">
	<label class="block w-full max-w-md text-sm text-ink-600 dark:text-ink-300">
		<span class="mb-2 block">查找标题或摘要</span>
		<input
			type="search"
			bind:value={query}
			placeholder="输入一个词，寻找一段文字…"
			class="w-full rounded-lg border border-ink-200 bg-transparent px-4 py-3 outline-none focus:border-jade-500 dark:border-ink-700"
		/>
	</label>
	<p aria-live="polite" class="font-mono text-xs text-ink-400">{visibleCount} 个阅读入口</p>
</div>

{#each groups as group (group.kind)}
	<section aria-labelledby={`sitemap-${group.kind}`} class="mb-12">
		<div class="mb-5 flex flex-wrap items-baseline gap-x-4 gap-y-2">
			<h2 id={`sitemap-${group.kind}`} class="font-serif text-2xl text-ink-900 dark:text-ink-100">
				{group.title}
			</h2>
			<p class="text-sm text-ink-400">{group.description}</p>
		</div>
		<ul class="grid gap-x-10 sm:grid-cols-2">
			{#each group.entries as entry (entry.path)}
				<li class="border-b border-ink-100 py-4 dark:border-ink-800">
					<a
						href={resolveHref(entry.path)}
						class="font-medium text-ink-700 transition-colors hover:text-jade-600 dark:text-ink-200 dark:hover:text-jade-400"
						>{entry.title}</a
					>
					{#if entry.modified}<time
							datetime={entry.modified}
							class="mt-2 block font-mono text-xs text-ink-400"
							>更新于 {entry.modified.slice(0, 10)}</time
						>{/if}
				</li>
			{/each}
		</ul>
	</section>
{/each}
{#if groups.length === 0}<p class="py-12 text-center text-ink-500">
		没有找到匹配的内容，试试另一个关键词。
	</p>{/if}
