<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import Sidebar from '$lib/ui/layout/sidebar/Sidebar.svelte';
	import MobileNavBar from '$lib/ui/layout/sidebar/MobileNavBar.svelte';
	import { initTheme, startThemeSync, themeManager } from '$lib/shared/theme/theme.svelte.js';
	import { onMount } from 'svelte';
	import { consoleLogInfo } from '$lib/features/console-info/index';
	import Toaster from '$lib/ui/primitives/toaster/Toaster.svelte';
	import QueryRoot from '$lib/ui/common/QueryRoot.svelte';
	import Loading from '$lib/ui/common/Loading.svelte';
	import { navigating } from '$app/stores';
	import { browser } from '$app/environment';
	import { page } from '$app/state';
	import { beforeNavigate, onNavigate } from '$app/navigation';
	import SearchModal from '$lib/ui/search/SearchModal.svelte';
	import Footer from '$lib/ui/layout/Footer.svelte';
	import FloatingWindow from '$lib/ui/common/FloatingWindow.svelte';
	import { uiState } from '$lib/shared/stores/ui.svelte';
	import { windowStore } from '$lib/shared/stores/windowStore.svelte';
	import { presenceStore } from '$lib/features/presence/store.svelte';
	import { ownerStatusStore } from '$lib/features/owner-status/store.svelte';
	import { detailHeroBgSrc } from '$lib/shared/stores/detailHeroBg';
	import DetailHeroBg from '$lib/ui/detail/DetailHeroBg.svelte';
	import { siteHealthStore } from '$lib/features/site-health/store.svelte';
	import SiteHealthBanner from '$lib/features/site-health/components/SiteHealthBanner.svelte';
	import { resolvePresenceView } from '$lib/features/presence/resolve-view';
	import PresencePagesWindow from '$lib/features/presence/components/PresencePagesWindow.svelte';
	import ThinkingCommentsWindow from '$lib/features/thinking/components/ThinkingCommentsWindow.svelte';
	import { getToken } from '$lib/shared/token';
	import { getProfile } from '$lib/features/auth/api';
	import { userStore } from '$lib/shared/stores/userStore';
	import { get } from 'svelte/store';

	function logClientRuntimeError(
		kind: 'error' | 'unhandledrejection',
		message: string,
		detail: string
	) {
		console.error(
			`[renderer][client-error] side=client code=runtime kind=${kind} path=${page.url.pathname} message=${message}\n${detail}`
		);
	}

	function handleKeydown(event: KeyboardEvent) {
		if ((event.metaKey || event.ctrlKey) && (event.key === 'k' || event.key === 'K')) {
			event.preventDefault();
			uiState.toggleSearch();
		}
	}

	function isAlbumDetailPath(pathname: string | null | undefined) {
		return typeof pathname === 'string' && /^\/albums\/[^/]+$/.test(pathname);
	}

	function isAlbumPhotoPath(pathname: string | null | undefined) {
		return typeof pathname === 'string' && /^\/albums\/[^/]+\/photo\/[^/]+$/.test(pathname);
	}

	function shouldSkipNativeViewTransition(
		fromPath: string | null | undefined,
		toPath: string | null | undefined
	) {
		return (
			(isAlbumDetailPath(fromPath) && isAlbumPhotoPath(toPath)) ||
			(isAlbumPhotoPath(fromPath) && isAlbumDetailPath(toPath))
		);
	}

	/**
	 * Avoid back animation and LCP delay caused by lang time animation.
	 * reference: https://innei.in/posts/design/page-transition-animation-and-lcp
	 */
	beforeNavigate(({ type, willUnload }) => {
		if (willUnload || typeof document === 'undefined') return;
		document.documentElement.dataset.navType = type === 'popstate' ? 'back' : 'forward';
	});

	onNavigate((navigation) => {
		if (typeof document === 'undefined' || !document.startViewTransition) return;
		const fromPath = navigation.from?.url.pathname;
		const toPath = navigation.to?.url.pathname;
		if (shouldSkipNativeViewTransition(fromPath, toPath)) return;
		const startViewTransition = document.startViewTransition.bind(document);

		return new Promise((resolve) => {
			startViewTransition(async () => {
				resolve();
				await navigation.complete;
			});
		});
	});

	import '@fontsource/google-sans/400.css';
	import '@fontsource/noto-serif-sc/400.css';
	import '@fontsource/noto-serif-sc/500.css';
	import '@fontsource/noto-serif-sc/600.css';
	import '@fontsource/noto-serif-sc/700.css';
	import '@fontsource-variable/victor-mono/index.css';
	import { websiteInfoCtx } from '$lib/features/website-info/context.js';
	import { resolveSeoMeta } from '$lib/shared/seo/metadata';
	import { resolveHomeThemeConfig } from '$lib/features/home/theme';
	import {
		createEmptyDetailPanelModel,
		detailPanelCtx,
		type DetailPanelModel,
		type DetailPanelRelatedMoment,
		type DetailPanelRelatedPost
	} from '$lib/shared/detail-panel/context';

	type ThinkingWindowData = {
		areaId?: number | null;
		commentsCount?: number;
		thinkingId?: number;
		activityPubObjectId?: string | null;
	};

	let { children, data } = $props();
	let showRouteLoading = $state(false);

	websiteInfoCtx.mountModelData(() => data.websiteInfo ?? null);

	// Initialize health state from SSR data (one-time).
	$effect(() => {
		if (data.healthData) {
			siteHealthStore.initFromSSR(data.healthData);
		}
	});

	const readDetailPanelFromPageData = (view: unknown): DetailPanelModel => {
		const empty = createEmptyDetailPanelModel();
		if (!view || typeof view !== 'object') return empty;
		const viewData = view as {
			post?: {
				title?: string | null;
				toc?: DetailPanelModel['toc'] | null;
				relatedMoments?: DetailPanelRelatedMoment[] | null;
			};
			moment?: {
				title?: string | null;
				toc?: DetailPanelModel['toc'] | null;
				relatedPosts?: DetailPanelRelatedPost[] | null;
			};
			page?: {
				title?: string | null;
				toc?: DetailPanelModel['toc'] | null;
			};
		};

		if (viewData.post) {
			return {
				...empty,
				kind: 'post',
				title: viewData.post.title ?? '',
				toc: viewData.post.toc ?? [],
				relatedMoments: (viewData.post.relatedMoments ?? []).slice(0, 2)
			};
		}
		if (viewData.moment) {
			return {
				...empty,
				kind: 'moment',
				title: viewData.moment.title ?? '',
				toc: viewData.moment.toc ?? [],
				relatedPosts: (viewData.moment.relatedPosts ?? []).slice(0, 2)
			};
		}
		if (viewData.page) {
			return {
				...empty,
				kind: 'page',
				title: viewData.page.title ?? '',
				toc: viewData.page.toc ?? []
			};
		}

		return empty;
	};

	detailPanelCtx.mountModelData(() => {
		const pathname = page.url.pathname;
		const pageData = page.data;
		const hash = browser ? window.location.hash.replace(/^#/, '') : '';
		return {
			...readDetailPanelFromPageData(pageData),
			contentRoot: null,
			activeAnchor:
				pathname === (browser ? window.location.pathname : pathname) ? hash || null : null
		};
	});

	const websiteInfoStore = websiteInfoCtx.selectModelData((model) => model ?? null);
	const avatarOrigin = $derived.by(() => {
		const url = resolveHomeThemeConfig($websiteInfoStore).hero?.avatarUrl;
		if (!url) return null;
		try {
			return new URL(url).origin;
		} catch {
			return null;
		}
	});
	const normalizeIconUrl = (value: unknown): string =>
		typeof value === 'string' ? value.trim() : '';
	const inferIconMimeType = (iconUrl: string): string | null => {
		const lower = iconUrl.toLowerCase();
		const cleaned = lower.split('#')[0]?.split('?')[0] || lower;
		if (cleaned.endsWith('.svg')) return 'image/svg+xml';
		if (cleaned.endsWith('.png')) return 'image/png';
		if (cleaned.endsWith('.jpg') || cleaned.endsWith('.jpeg')) return 'image/jpeg';
		if (cleaned.endsWith('.webp')) return 'image/webp';
		if (cleaned.endsWith('.ico')) return 'image/x-icon';
		return null;
	};
	const siteFavicon = $derived.by(() => normalizeIconUrl($websiteInfoStore?.favicon) || favicon);
	const siteFaviconType = $derived.by(() => inferIconMimeType(siteFavicon));

	// Clip favicon to circle via Canvas
	let circularFaviconUrl = $state('');
	$effect(() => {
		const src = siteFavicon;
		if (!browser || !src) return;

		let cancelled = false;
		const size = 128;
		const img = new Image();
		img.crossOrigin = 'anonymous';
		img.onload = () => {
			if (cancelled) return;
			try {
				const canvas = document.createElement('canvas');
				canvas.width = size;
				canvas.height = size;
				const ctx = canvas.getContext('2d');
				if (!ctx) return;
				ctx.beginPath();
				ctx.arc(size / 2, size / 2, size / 2, 0, Math.PI * 2);
				ctx.closePath();
				ctx.clip();
				ctx.drawImage(img, 0, 0, size, size);
				circularFaviconUrl = canvas.toDataURL('image/png');
			} catch {
				// Canvas tainted by CORS, keep original
			}
		};
		img.src = src;
		return () => {
			cancelled = true;
		};
	});
	const resolvedFavicon = $derived(circularFaviconUrl || siteFavicon);
	const resolvedFaviconType = $derived(
		circularFaviconUrl ? 'image/png' : siteFaviconType || undefined
	);
	const thinkingWindowData = $derived((windowStore.data ?? {}) as ThinkingWindowData);

	const seoMeta = $derived.by(() =>
		resolveSeoMeta({
			pathname: page.url.pathname,
			search: page.url.search,
			routeData: page.data,
			websiteInfo: $websiteInfoStore,
			origin: page.url.origin,
			fallbackSiteIcon: siteFavicon
		})
	);

	// Initialize theme on mount
	const theme = themeManager;

	function openPresenceWindow() {
		windowStore.open('在线页面', null, 'presence-pages');
	}

	onMount(() => {
		const handleWindowError = (event: ErrorEvent) => {
			const message = event.message || 'Unhandled client error';
			const detail =
				event.error instanceof Error
					? `${event.error.name}: ${event.error.message}${event.error.stack ? `\n${event.error.stack}` : ''}`
					: `${event.filename || 'unknown'}:${event.lineno || 0}:${event.colno || 0}`;
			logClientRuntimeError('error', message, detail);
		};
		const handleUnhandledRejection = (event: PromiseRejectionEvent) => {
			const reason = event.reason;
			const detail =
				reason instanceof Error
					? `${reason.name}: ${reason.message}${reason.stack ? `\n${reason.stack}` : ''}`
					: String(reason);
			logClientRuntimeError('unhandledrejection', 'Unhandled promise rejection', detail);
		};
		window.addEventListener('error', handleWindowError);
		window.addEventListener('unhandledrejection', handleUnhandledRejection);

		initTheme(theme);
		consoleLogInfo();
		presenceStore.start();
		ownerStatusStore.start();
		siteHealthStore.start();

		// Auth bootstrap: imperatively check token and fetch profile.
		// This replaces the old AuthBootstrap component, avoiding the race
		// condition caused by QueryRoot's multiple async imports.
		const token = getToken();
		if (token && !get(userStore).isLogin) {
			getProfile()
				.then((profile) => {
					if (profile && !get(userStore).isLogin) {
						userStore.setUser(profile);
					}
				})
				.catch(() => {
					/* token invalid or expired, ignore */
				});
		}

		return () => {
			window.removeEventListener('error', handleWindowError);
			window.removeEventListener('unhandledrejection', handleUnhandledRejection);
			presenceStore.stop();
			ownerStatusStore.stop();
			siteHealthStore.stop();
		};
	});

	startThemeSync(theme);

	$effect(() => {
		if (!browser) {
			return;
		}

		if ($navigating) {
			showRouteLoading = false;
			const timer = setTimeout(() => {
				if ($navigating) {
					showRouteLoading = true;
				}
			}, 2000);

			return () => {
				clearTimeout(timer);
			};
		}

		showRouteLoading = false;
	});

	$effect(() => {
		if (!browser) return;

		const report = resolvePresenceView(page.url.pathname, page.data);
		if (!report) return;

		presenceStore.reportView(report);
	});
</script>

<svelte:head>
	{#if avatarOrigin}
		<link rel="preconnect" href={avatarOrigin} crossorigin="anonymous" />
		<link rel="dns-prefetch" href={avatarOrigin} />
	{/if}
	<link rel="icon" href={resolvedFavicon} type={resolvedFaviconType} />
	<link rel="shortcut icon" href={resolvedFavicon} />
	<link rel="apple-touch-icon" href={resolvedFavicon} />
	<title>{seoMeta.title}</title>
	<link rel="canonical" href={seoMeta.canonicalUrl} />
	<link rel="describedby" href={seoMeta.llmsUrl} type="text/markdown" />
	{#if seoMeta.markdownUrl}
		<link rel="alternate" href={seoMeta.markdownUrl} type="text/markdown" />
	{/if}
	<meta
		name="viewport"
		content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no"
	/>
	<meta name="description" content={seoMeta.description} />
	<meta name="keywords" content={seoMeta.keywords} />
	<meta name="robots" content={seoMeta.robots} />
	<meta name="author" content="grtinry43" />
	<meta property="og:title" content={seoMeta.ogTitle} />
	<meta property="og:description" content={seoMeta.ogDescription} />
	<meta property="og:type" content={seoMeta.ogType} />
	<meta property="og:url" content={seoMeta.ogUrl} />
	<meta property="og:site_name" content={seoMeta.ogSiteName} />
	<meta property="og:image" content={seoMeta.ogImage} />
	{#if seoMeta.ogImageType}
		<meta property="og:image:type" content={seoMeta.ogImageType} />
	{/if}
	{#if seoMeta.ogImageWidth}
		<meta property="og:image:width" content={String(seoMeta.ogImageWidth)} />
	{/if}
	{#if seoMeta.ogImageHeight}
		<meta property="og:image:height" content={String(seoMeta.ogImageHeight)} />
	{/if}
	<meta name="twitter:card" content={seoMeta.twitterCard} />
	<meta name="twitter:title" content={seoMeta.ogTitle} />
	<meta name="twitter:description" content={seoMeta.ogDescription} />
	<meta name="twitter:image" content={seoMeta.ogImage} />
	<script>
		// Inline script to prevent theme flash (fallback before Svelte hydrates)
		(function () {
			try {
				const theme = localStorage.getItem('theme') || 'system';
				const isDark =
					theme === 'dark' ||
					(theme === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches);
				document.documentElement.classList.toggle('dark', isDark);
			} catch (e) {}
		})();
	</script>
</svelte:head>

<div class="hidden md:block">
	<Sidebar menuTree={data.navMenus ?? []} />
</div>
<MobileNavBar menuTree={data.navMenus ?? []} />
<!-- noise background -->
<div class="bg-noise" aria-hidden="true"></div>

<div class="md:pl-24 transition-[padding] duration-300 relative overflow-x-clip">
	{#if $detailHeroBgSrc}
		<DetailHeroBg src={$detailHeroBgSrc} />
	{/if}
	<SiteHealthBanner />
	<main
		class="page-wrapper mx-auto {page.url.pathname.startsWith('/timeline')
			? 'max-w-none px-0 py-0'
			: 'max-w-300 px-4 sm:px-6 lg:px-8 py-10 md:py-16'}"
	>
		<div class="content-container min-h-[60vh]">
			{@render children()}
		</div>
	</main>
	<Footer
		onlineCount={presenceStore.online}
		presenceConnected={presenceStore.isConnected}
		onOpenPresence={openPresenceWindow}
	/>
</div>

{#if showRouteLoading}
	<div
		class="fixed px-12 py-6 left-1/2 top-1/2 z-99999 -translate-x-1/2 -translate-y-1/2 pointer-events-none rounded-default border border-ink-200/70 bg-ink-50/80 shadow-subtle backdrop-blur-lg dark:border-ink-700/70 dark:bg-ink-900/80"
		aria-live="polite"
		aria-busy="true"
	>
		<Loading size="w-8 h-8" duration={900} class="gap-0" text="正在玩命加载中...莫慌" />
	</div>
{/if}

<SearchModal />
<FloatingWindow>
	<!-- Login branch: always mounted (hidden when inactive) to preserve QueryRoot/AuthClient state -->
	<div hidden={windowStore.kind !== 'login'}>
		<QueryRoot
			loader={() => import('$lib/features/auth/components/AuthClient.svelte')}
			fallback={authFallback}
		/>
	</div>
	{#if windowStore.kind === 'tag-contents'}
		<QueryRoot
			loader={() => import('$lib/features/tag/components/TagContentsWindow.svelte')}
			loaderProps={{ tagId: windowStore.data?.id, tagName: windowStore.data?.name }}
		/>
	{:else if windowStore.title === '申请友链'}
		<QueryRoot
			loader={() => import('$lib/features/friend-link/components/ApplyFriendForm.svelte')}
		/>
	{:else if windowStore.kind === 'presence-pages'}
		<PresencePagesWindow />
	{:else if windowStore.kind === 'thinking-comments'}
		<ThinkingCommentsWindow
			areaId={thinkingWindowData.areaId ?? null}
			commentsCount={thinkingWindowData.commentsCount ?? 0}
			thinkingId={thinkingWindowData.thinkingId ?? 0}
			activityPubObjectId={thinkingWindowData.activityPubObjectId ?? null}
		/>
	{:else if windowStore.kind === 'user-center'}
		<QueryRoot
			loader={() => import('$lib/features/user-center/components/UserCenterWindow.svelte')}
		/>
	{:else if windowStore.kind === 'share-poster'}
		<QueryRoot loader={() => import('$lib/ui/share/SharePosterWindow.svelte')} />
	{:else if windowStore.kind !== 'login'}
		<div class="flex flex-col gap-3"></div>
	{/if}
</FloatingWindow>

<svelte:window onkeydown={handleKeydown} />

<Toaster />
{#snippet globalNotificationFallback()}
	<div></div>
{/snippet}
<QueryRoot
	loader={() =>
		import('$lib/features/global-notification/components/GlobalNotificationClient.svelte')}
	fallback={globalNotificationFallback}
/>
{#snippet authFallback()}
	<div></div>
{/snippet}

<style lang="postcss">
	@reference "./layout.css";

	:global(html) {
		scroll-behavior: smooth;
		scroll-padding-top: 80px;
	}
</style>
