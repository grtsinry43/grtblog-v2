<script lang="ts">
	import { resolvePath } from '$lib/shared/utils/resolve-path';
	import { resolveFooterThemeConfig } from '$lib/features/footer/theme';
	import { websiteInfoCtx } from '$lib/features/website-info/context';
	import { onMount } from 'svelte';

	type Props = {
		onlineCount?: number;
		presenceConnected?: boolean;
		onOpenPresence?: () => void;
	};

	let { onlineCount = 0, presenceConnected = false, onOpenPresence = () => {} }: Props = $props();

	const currentYear = new Date().getFullYear();
	const footerThemeStore = websiteInfoCtx.selectModelData((data) => resolveFooterThemeConfig(data));
	const desktopColumnCount = $derived(Math.min($footerThemeStore.sections.length, 6));
	const desktopGridStyle = $derived.by(() => {
		const sectionCols = Math.max(desktopColumnCount, 1);
		return `repeat(${sectionCols}, minmax(0, 1fr)) minmax(220px, 1.2fr)`;
	});
	const brandFullRow = $derived($footerThemeStore.sections.length + 1 > 6);
	let nowMs = $state(0);

	const formatPresenceText = (template: string, count: number): string =>
		template.replaceAll('{count}', String(count));

	const formatUptimeText = (
		template: string,
		parts: { days: number; hours: number; minutes: number; seconds: number; totalSeconds: number }
	): string =>
		template
			.replaceAll('{days}', String(parts.days))
			.replaceAll('{hours}', String(parts.hours))
			.replaceAll('{minutes}', String(parts.minutes))
			.replaceAll('{seconds}', String(parts.seconds))
			.replaceAll('{totalSeconds}', String(parts.totalSeconds));

	const preloadDataAttr = (href: string): 'off' | undefined => {
		if (!href || /^(https?:|mailto:)/i.test(href)) return undefined;
		const path = href.split(/[?#]/, 1)[0];
		return path === '/feed' || path === '/rss.xml' ? 'off' : undefined;
	};

	const uptimeText = $derived.by(() => {
		if (nowMs <= 0) return '';
		const startAt = Date.parse($footerThemeStore.siteStartTime);
		if (!Number.isFinite(startAt)) return '';
		const totalSeconds = Math.max(0, Math.floor((nowMs - startAt) / 1000));
		const days = Math.floor(totalSeconds / 86400);
		const hours = Math.floor((totalSeconds % 86400) / 3600);
		const minutes = Math.floor((totalSeconds % 3600) / 60);
		const seconds = totalSeconds % 60;
		return formatUptimeText($footerThemeStore.uptimeTextTemplate, {
			days,
			hours,
			minutes,
			seconds,
			totalSeconds
		});
	});

	onMount(() => {
		nowMs = Date.now();
		const timer = window.setInterval(() => {
			nowMs = Date.now();
		}, 1000);
		return () => window.clearInterval(timer);
	});
</script>

<footer
	class="mt-32 border-t border-jade-100/80 dark:border-ink-800 bg-jade-50/30 dark:bg-ink-950/30 backdrop-blur-sm"
>
	<div class="max-w-[1200px] mx-auto px-6 py-12 md:py-16">
		<!-- Mobile Compact Layout (Hidden on Desktop) -->
		<div class="flex flex-col gap-4 mb-12 md:hidden">
			{#each $footerThemeStore.sections as section (section.title)}
				<div class="flex flex-col gap-2">
					<div
						class="text-sm font-serif font-bold text-ink-900 dark:text-ink-100 flex items-center justify-between"
					>
						{section.title}
						<span class="text-ink-300 dark:text-ink-700 font-mono font-normal">></span>
					</div>
					<div class="flex flex-wrap gap-x-4 gap-y-2">
						{#each section.links as link (link.name)}
							<a
								href={/^(https?:|mailto:)/i.test(link.href) ? link.href : resolvePath(link.href)}
								data-sveltekit-preload-data={preloadDataAttr(link.href)}
								class="text-sm text-ink-500 hover:text-jade-600 dark:hover:text-jade-400 transition-colors"
							>
								{link.name}
							</a>
						{/each}
					</div>
				</div>
			{/each}

			<!-- Brand Info below mobile footer links -->
			<div class="flex flex-col gap-4 pt-4 border-t border-ink-100 dark:border-ink-800/50">
				<div class="flex flex-col">
					<div class="text-lg font-mono font-bold text-ink-900 dark:text-ink-100">
						{$footerThemeStore.brandName}
					</div>
					<p class="text-[11px] font-mono text-ink-400 mt-1 uppercase tracking-wider">
						{$footerThemeStore.brandTagline}
					</p>
				</div>
				<button
					onclick={onOpenPresence}
					class="flex items-center gap-2 w-fit transition-colors underline-offset-2 hover:underline focus-visible:underline focus-visible:outline-none {presenceConnected
						? 'text-jade-700/80 dark:text-jade-400/80'
						: 'text-red-600 dark:text-red-400'}"
				>
					<span class="relative flex h-1.5 w-1.5">
						<span
							class="absolute inline-flex h-full w-full rounded-full opacity-75 {presenceConnected
								? 'bg-jade-400'
								: 'bg-red-400'}"
						></span>
						<span
							class="relative inline-flex rounded-full h-1.5 w-1.5 {presenceConnected
								? 'bg-jade-500'
								: 'bg-red-500'}"
						></span>
					</span>
					<span class="text-[10px] font-mono">
						{#if presenceConnected}
							{formatPresenceText($footerThemeStore.presenceConnectedText, onlineCount)}
						{:else}
							{$footerThemeStore.presenceLoadingText}
						{/if}
					</span>
				</button>
			</div>
		</div>

		<!-- Desktop Multi-column Layout (Hidden on Mobile) -->
		<div
			class="hidden md:grid gap-8 lg:gap-12 mb-16"
			style:grid-template-columns={desktopGridStyle}
		>
			{#each $footerThemeStore.sections as section (section.title)}
				<div class="flex flex-col gap-6">
					<h3
						class="text-sm font-serif font-bold text-ink-900 dark:text-ink-100 flex items-center gap-2"
					>
						<span class="w-1 h-3 bg-jade-500 rounded-full"></span>
						{section.title}
					</h3>
					<ul class="flex flex-col gap-3">
						{#each section.links as link (link.name)}
							<li>
								<a
									href={/^(https?:|mailto:)/i.test(link.href) ? link.href : resolvePath(link.href)}
									data-sveltekit-preload-data={preloadDataAttr(link.href)}
									class="text-sm text-ink-500 hover:text-jade-600 dark:hover:text-jade-400 transition-colors"
								>
									{link.name}
								</a>
							</li>
						{/each}
					</ul>
				</div>
			{/each}

			<!-- Brand Info inside Desktop Grid -->
			<div
				class="flex flex-col gap-6 items-end text-right"
				style:grid-column={brandFullRow ? '1 / -1' : undefined}
			>
				<div class="flex flex-col items-end">
					<div class="text-xl font-mono font-bold text-ink-900 dark:text-ink-100">
						{$footerThemeStore.brandName}
					</div>
					<p class="text-[11px] font-mono text-ink-400 mt-1 uppercase tracking-wider">
						{$footerThemeStore.brandTagline}
					</p>
				</div>
				<button
					onclick={onOpenPresence}
					class="flex items-center gap-2 w-fit transition-colors underline-offset-2 hover:underline focus-visible:underline focus-visible:outline-none {presenceConnected
						? 'text-jade-700/80 dark:text-jade-400/80'
						: 'text-red-600 dark:text-red-400'}"
				>
					<span class="relative flex h-1.5 w-1.5">
						<span
							class="absolute inline-flex h-full w-full rounded-full opacity-75 {presenceConnected
								? 'bg-jade-400'
								: 'bg-red-400'}"
						></span>
						<span
							class="relative inline-flex rounded-full h-1.5 w-1.5 {presenceConnected
								? 'bg-jade-500'
								: 'bg-red-500'}"
						></span>
					</span>
					<span class="text-[10px] font-mono">
						{#if presenceConnected}
							{formatPresenceText($footerThemeStore.presenceConnectedText, onlineCount)}
						{:else}
							{$footerThemeStore.presenceLoadingText}
						{/if}
					</span>
				</button>
			</div>
		</div>

		<!-- Bottom Copyright (Universal) -->
		<div
			class="flex flex-col md:flex-row justify-between items-center pt-8 border-t border-ink-100 dark:border-ink-800/50 gap-4"
		>
			<div class="text-[10px] md:text-[11px] font-mono text-ink-400 text-center md:text-left">
				<p>
					Copyright © {$footerThemeStore.copyrightStartYear} - {currentYear}
					{$footerThemeStore.copyrightOwner}. All rights reserved.
				</p>
				<div class="flex flex-wrap justify-center md:justify-start gap-x-3 mt-1">
					<span class="hidden md:inline"
						>Powered by <a
							href="https://grtblog.js.org/"
							class="text-jade-500 hover:text-jade-600 transition-colors">Grtblog-v2</a
						></span
					>
					{#if uptimeText}
						<span class="hidden md:inline text-ink-200 dark:text-ink-800">|</span>
						<span class="hidden md:inline">{uptimeText}</span>
					{/if}
					{#if $footerThemeStore.beianText && $footerThemeStore.beianUrl}
						<span class="hidden md:inline text-ink-200 dark:text-ink-800">|</span>
						<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
						<a
							href={$footerThemeStore.beianUrl}
							target="_blank"
							rel="noreferrer"
							class="hover:text-jade-600 transition-colors">{$footerThemeStore.beianText}</a
						>
					{/if}
					{#if $footerThemeStore.beianGongAnText}
						<span class="hidden md:inline text-ink-200 dark:text-ink-800">|</span>
						<span>{$footerThemeStore.beianGongAnText}</span>
					{/if}
				</div>
			</div>

			<div class="hidden md:flex items-center gap-4 text-[11px] font-mono text-ink-300">
				<span>{$footerThemeStore.designedWithText}</span>
			</div>
		</div>
	</div>
</footer>

<style lang="postcss">
	@reference "$routes/layout.css";
</style>
