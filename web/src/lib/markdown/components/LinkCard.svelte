<script lang="ts">
	let { href = '', title = '', desc = '', newtab = 'true', contentHtml = '' } = $props<{
		href?: string;
		title?: string;
		desc?: string;
		newtab?: string | boolean;
		contentHtml?: string;
	}>();

	const openInNewTab = typeof newtab === 'string' ? newtab !== 'false' : Boolean(newtab);
	const target = openInNewTab ? '_blank' : '_self';
	const rel = openInNewTab ? 'noreferrer' : undefined;
</script>

<a
	class="group block rounded-2xl border border-ink-200/70 bg-white/80 p-6 shadow-subtle transition hover:-translate-y-0.5 hover:shadow-float"
	href={href || '#'}
	target={target}
	rel={rel}
>
	<div class="flex items-start gap-4">
		<span class="inline-flex h-10 w-10 items-center justify-center rounded-2xl bg-ink-100 text-ink-600">
			<svg class="h-5 w-5" viewBox="0 0 24 24" fill="none">
				<path
					d="M14 5h5v5m-9 9h-5v-5m14-4L10 19m9-9L5 19"
					stroke="currentColor"
					stroke-width="1.8"
					stroke-linecap="round"
					stroke-linejoin="round"
				/>
			</svg>
		</span>
		<div class="space-y-2">
			<div class="flex flex-wrap items-center gap-2">
				<h4 class="text-lg font-semibold text-ink-900">{title}</h4>
				<span class="rounded-full border border-ink-200/80 bg-white/70 px-2.5 py-0.5 text-[11px] uppercase tracking-[0.18em] text-ink-500">
					Link
				</span>
			</div>
			{#if contentHtml}
				<div class="text-sm leading-relaxed text-ink-600">
					{@html contentHtml}
				</div>
			{:else if desc}
				<p class="text-sm leading-relaxed text-ink-600">{desc}</p>
			{/if}
		</div>
	</div>
</a>
