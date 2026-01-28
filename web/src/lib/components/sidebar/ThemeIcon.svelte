<script lang="ts">
	import DynamicLucideIcon from '../icons/DynamicLucideIcon.svelte';
	import { resolveTheme, themeManager } from '$lib/shared/theme.svelte';

	const theme = themeManager;
	const resolved = $derived.by(() => resolveTheme(theme.current));
	const iconName = $derived.by(() => (resolved === 'dark' ? 'Moon' : 'Sun'));

	type ViewTransitionLike = { ready: Promise<void> };
	type DocumentWithViewTransition = Document & {
		startViewTransition?: (callback: () => void) => ViewTransitionLike;
	};

	const toggleTheme = async (event: MouseEvent) => {
		const next = resolved === 'dark' ? 'light' : 'dark';
		const doc = document as DocumentWithViewTransition;
		if (!doc.startViewTransition) {
			theme.set(next);
			return;
		}

		const x = event.clientX;
		const y = event.clientY;
		const endRadius = Math.hypot(
			Math.max(x, window.innerWidth - x),
			Math.max(y, window.innerHeight - y)
		);

		const transition = doc.startViewTransition.call(doc, () => {
			theme.set(next);
		});

		await transition.ready;

		document.documentElement.animate(
			{
				clipPath: [`circle(0px at ${x}px ${y}px)`, `circle(${endRadius}px at ${x}px ${y}px)`]
			},
			{
				duration: 500,
				easing: 'ease-in-out',
				pseudoElement: '::view-transition-new(root)'
			}
		);
	};
</script>

<button
	type="button"
	class="theme-btn"
	data-theme={resolved}
	aria-label={`Switch to ${resolved === 'dark' ? 'light' : 'dark'} theme`}
	onclick={toggleTheme}
>
	<DynamicLucideIcon name={iconName} size={18} className="theme-icon" />
</button>

<style lang="postcss">
	@reference "$routes/layout.css";

	.theme-btn {
		@apply relative grid h-10 w-10 place-items-center overflow-hidden rounded-full border
		 border-ink-200 dark:border-ink-700 bg-ink-200/80 dark:bg-ink-700/80 text-ink-200 transition-colors duration-300;
		@apply hover:border-ink-300 hover:bg-ink-300/70 hover:dark:border-ink-600 hover:dark:bg-ink-600/70;
		@apply text-ink-900 dark:text-ink-100;
	}

	.theme-icon {
		@apply relative z-10 ;
	}
</style>
