<script lang="ts">
	import { Search, Sparkles, ArrowRight } from 'lucide-svelte';
	import Card from '$lib/components/Card.svelte';
	import Button from '$lib/components/ui/button/Button.svelte';
	import Input from '$lib/components/Input.svelte';
	import Badge from '$lib/components/Badge.svelte';
	import Tag from '$lib/components/Tag.svelte';
	import Divider from '$lib/components/Divider.svelte';

	// Mock Data for Preview
	const featuredPost = {
		title: '极简设计的呼吸感',
		excerpt: '在信息过载的时代，我们探讨如何通过留白与秩序，打造让用户能够自由呼吸的数字交互界面。',
		date: '2024年10月24日',
		category: '设计',
		image:
			'https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?auto=format&fit=crop&q=80&w=1000'
	};

	const posts = [
		{
			title: '构建未来的 Web 架构',
			excerpt: '深入理解组件驱动架构的演变，以及它对现代开发者和用户体验的深远影响。',
			date: '2024年10月22日',
			category: '技术'
		},
		{
			title: '数字花园：一种新的书写方式',
			excerpt: '为什么我们应该将个人网站视为不断生长、演化的花园，而非仅仅是静态的作品展示。',
			date: '2024年10月18日',
			category: '思考'
		}
	];

	let searchQuery = $state('');
</script>

<div class="homepage-container">
	{#snippet heroButtonContent()}
		34324234
	{/snippet}

	{#snippet readMoreContent()}
		阅读全文
		<ArrowRight size={14} class="ml-1.5 transition-transform group-hover:translate-x-1" />
	{/snippet}

	{#snippet viewAllContent()}
		查看全部
	{/snippet}

	<!-- Hero Section -->
	<section class="hero-section">
		<div class="hero-content">
			<Badge variant="soft">
				<Sparkles size={10} />
				匠心营造
			</Badge>

			<Button content={heroButtonContent} />

			<h1 class="hero-title">
				秩序与<br />
				<span class="jade-italic">呼吸</span>
			</h1>

			<p class="hero-desc">
				数字工艺的避风港。探索东方美学与现代界面设计的交汇点，在简洁中寻找深邃，在方寸间感悟平衡。
			</p>

			<div class="search-wrapper">
				{#snippet searchIcon()}
					<Search size={14} />
				{/snippet}
				<Input bind:value={searchQuery} placeholder="搜索观点、项目或灵感..." icon={searchIcon} />
			</div>
		</div>
	</section>

	<!-- Featured Post -->
	<section class="featured-section">
		<Divider label="推荐文章" />

		<div class="featured-card group">
			<div class="featured-grid">
				<div class="featured-image-wrapper">
					<img src={featuredPost.image} alt={featuredPost.title} class="featured-image" />
				</div>
				<div class="featured-info">
					<div class="featured-meta">
						<Tag variant="jade">{featuredPost.category}</Tag>
						<span class="meta-date">{featuredPost.date}</span>
					</div>

					<h2 class="featured-title">
						{featuredPost.title}
					</h2>

					<p class="featured-excerpt">
						{featuredPost.excerpt}
					</p>

					<div class="featured-actions">
						<Button variant="ghost" class="text-link" content={readMoreContent} />
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Recent Posts -->
	<section class="posts-section">
		<div class="section-header">
			<Tag variant="ink">近期文章</Tag>
			<Button variant="ghost" class="view-all-btn" content={viewAllContent} />
		</div>

		<div class="posts-list">
			{#each posts as post}
				<Card variant="seamless" hover={true} class="post-item group">
					<div class="post-content">
						<div class="post-meta">
							<span>{post.date}</span>
							<span class="meta-sep">/</span>
							<Tag variant="outline" class="border-none !p-0">{post.category}</Tag>
						</div>
						<h3 class="post-title">
							{post.title}
						</h3>
						<p class="post-excerpt">
							{post.excerpt}
						</p>
					</div>

					<div class="arrow-icon-box">
						<ArrowRight size={14} />
					</div>
				</Card>
			{/each}
		</div>
	</section>
</div>

<style>
	@reference "./layout.css";

	.homepage-container {
		@apply space-y-16 md:space-y-24;
	}

	.hero-section {
		@apply space-y-6;
	}

	.hero-content {
		@apply space-y-5;
	}

	.hero-title {
		@apply font-serif text-3xl leading-[1.2] font-medium tracking-tight text-ink-950 md:text-4xl lg:text-5xl dark:text-ink-50;
	}

	.jade-italic {
		@apply text-jade-800 italic dark:text-jade-400;
	}

	.hero-desc {
		@apply max-w-lg text-[13px] leading-relaxed font-normal text-ink-600 opacity-90 md:text-sm dark:text-ink-400;
	}

	.search-wrapper {
		@apply w-full max-w-sm pt-2;
	}

	.featured-section {
		@apply space-y-6;
	}

	.featured-card {
		@apply cursor-pointer;
	}

	.featured-grid {
		@apply grid items-center gap-6 lg:grid-cols-[1fr_360px] lg:gap-10;
	}

	.featured-image-wrapper {
		@apply relative aspect-[16/9] h-full min-h-[300px] overflow-hidden rounded-lg border border-ink-100/50 shadow-sm lg:aspect-auto dark:border-ink-800/30;
	}

	.featured-image {
		@apply absolute inset-0 h-full w-full object-cover transition-transform duration-700 ease-out group-hover:scale-105;
	}

	.featured-info {
		@apply space-y-4;
	}

	.featured-meta {
		@apply flex items-center gap-3;
	}

	.meta-date {
		@apply font-mono text-[9px] tracking-widest text-ink-400 uppercase;
	}

	.featured-title {
		@apply font-serif text-xl leading-tight font-medium text-ink-950 transition-colors group-hover:text-jade-800 md:text-2xl dark:text-ink-100 dark:group-hover:text-jade-400;
	}

	.featured-excerpt {
		@apply text-[13px] leading-relaxed font-normal text-ink-600 opacity-90 md:text-sm dark:text-ink-400;
	}

	.featured-actions {
		@apply pt-1;
	}

	.text-link {
		@apply !p-0 text-[11px] font-bold tracking-wider text-jade-700 hover:!bg-transparent dark:text-jade-400;
	}

	.posts-section {
		@apply space-y-6;
	}

	.section-header {
		@apply flex items-center justify-between;
	}

	.view-all-btn {
		@apply !p-0 font-mono text-[9px] tracking-[0.2em] text-ink-400 uppercase hover:!bg-transparent hover:text-ink-900;
	}

	.posts-list {
		@apply grid gap-px border-y border-ink-50 bg-ink-50 dark:border-ink-900/50 dark:bg-ink-900/50;
	}

	.post-item {
		@apply flex flex-col gap-5 bg-white p-5 md:flex-row md:items-center md:p-6 dark:bg-ink-950;
	}

	.post-content {
		@apply flex-1 space-y-2;
	}

	.post-meta {
		@apply flex items-center gap-2 font-mono text-[9px] tracking-widest text-ink-400 uppercase;
	}

	.meta-sep {
		@apply text-ink-100 dark:text-ink-800/50;
	}

	.post-title {
		@apply font-serif text-lg font-medium text-ink-900 transition-colors group-hover:text-jade-700 md:text-xl dark:text-ink-100 dark:group-hover:text-jade-400;
	}

	.post-excerpt {
		@apply max-w-2xl text-[11px] leading-relaxed font-normal text-ink-500 opacity-80 md:text-xs dark:text-ink-400;
	}

	.arrow-icon-box {
		@apply flex h-8 w-8 shrink-0 items-center justify-center rounded-full border border-ink-100 text-ink-400 transition-all duration-500 group-hover:border-jade-500 group-hover:bg-jade-500 group-hover:text-white dark:border-ink-800;
	}
</style>
