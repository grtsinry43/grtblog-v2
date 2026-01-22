<script lang="ts">
	let {
		url = '',
		title = '',
		type = 'page',
		cover = '',
		blur = '7px',
		contentHtml = ''
	} = $props<{
		url?: string;
		title?: string;
		type?: string;
		cover?: string;
		blur?: string;
		contentHtml?: string;
	}>();

	const target = url?.startsWith('http') ? '_blank' : '_self';
	const rel = target === '_blank' ? 'noreferrer' : undefined;
</script>

<article class="group relative overflow-hidden rounded-3xl border border-ink-200/70 bg-white/80 shadow-float">
	{#if cover}
		<div class="absolute inset-0">
			<img class="h-full w-full object-cover" src={cover} alt="" loading="lazy" />
			<div class="absolute inset-0 bg-white/70" style={`backdrop-filter: blur(${blur});`}></div>
		</div>
	{:else}
		<div class="absolute inset-0 bg-[radial-gradient(circle_at_top,_rgba(45,212,191,0.15),_transparent_55%)]"></div>
	{/if}
	<div class="relative z-10 flex flex-col gap-6 px-8 py-8 md:flex-row md:items-center md:justify-between">
		<div class="space-y-3">
			<p class="text-xs font-semibold uppercase tracking-[0.28em] text-ink-400">Annual Summary</p>
			<h3 class="text-3xl font-semibold text-ink-900 md:text-4xl">{title}</h3>
			{#if contentHtml}
				<div class="text-base leading-relaxed text-ink-700">
					{@html contentHtml}
				</div>
			{:else}
				<p class="text-base text-ink-700">年度回顾与总结</p>
			{/if}
		</div>
		<div class="flex flex-col items-start gap-3 sm:flex-row sm:items-center">
			<a
				href={url || '#'}
				target={target}
				rel={rel}
				class="inline-flex items-center justify-center rounded-2xl border border-ink-300/70 bg-white/70 px-5 py-3 text-xs font-semibold uppercase tracking-[0.2em] text-ink-700 shadow-subtle transition hover:-translate-y-0.5 hover:shadow-float"
			>
				Preview
			</a>
			<a
				href={url || '#'}
				target={target}
				rel={rel}
				class="inline-flex h-11 w-11 items-center justify-center rounded-full border border-ink-300/70 bg-white/70 text-ink-700 shadow-subtle transition hover:-translate-y-0.5 hover:shadow-float"
				aria-label="Open link"
			>
				<svg class="h-5 w-5" viewBox="0 0 24 24" fill="none">
					<path
						d="M9 7h8m0 0v8m0-8L7 17"
						stroke="currentColor"
						stroke-width="1.8"
						stroke-linecap="round"
						stroke-linejoin="round"
					/>
				</svg>
			</a>
			<span class="rounded-full border border-ink-200/80 bg-white/70 px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.22em] text-ink-500">
				{type}
			</span>
		</div>
	</div>
</article>
