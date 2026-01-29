<script lang="ts">
	import { get, writable } from 'svelte/store';
	import { PostDetail } from '$lib/features/post';
	import { checkPostLatest } from '$lib/features/post/api';
	import type { PostContentPayload, PostDetail as PostDetailModel } from '$lib/features/post/types';
	import { browser } from '$app/environment';
    import {SvelteURL} from "svelte/reactivity";
	import { postDetailCtx } from './post-detail-context.js';

	let { data } = $props();
	let socket: WebSocket | null = null;
	let updateHintTimer: ReturnType<typeof setTimeout> | null = null;

	const postDetailStore = postDetailCtx.mountModelData(data.post ?? null);
	const { updateModelData } = postDetailCtx.useModelActions();
	const postIdStore = postDetailCtx.selectModelData((data) => data?.id ?? null);
	const contentHashStore = postDetailCtx.selectModelData((data) => data?.contentHash ?? null);

	const showUpdateHint = writable(false);

	$effect(() => {
		postDetailCtx.syncModelData(postDetailStore, data.post ?? null);
	});

	const refreshPostIfNeeded = async () => {
		const id = get(postIdStore);
		const contentHash = get(contentHashStore);
		if (!id || !contentHash) {
			return;
		}

		try {
			const latest = await checkPostLatest(undefined, id, contentHash);
			if (!latest || latest.latest) {
				return;
			}

			updateModelData((prev) => {
				if (!prev) {
					return prev;
				}

				return {
					...prev,
					title: latest.title ?? prev.title,
					leadIn: latest.leadIn ?? prev.leadIn,
					toc: latest.toc ?? prev.toc,
					content: latest.content ?? prev.content,
					contentHash: latest.contentHash || prev.contentHash
				};
			});
			triggerUpdateHint();
		} catch {
			// Ignore check failures and keep the initial SSR content.
		}
	};

	const triggerUpdateHint = () => {
		showUpdateHint.set(true);
		if (updateHintTimer) {
			clearTimeout(updateHintTimer);
		}
		updateHintTimer = setTimeout(() => {
			showUpdateHint.set(false);
			updateHintTimer = null;
		}, 2400);
	};

	const connectToPostUpdates = (postId: number | null) => {
		if (!browser || postId === null) {
			return;
		}

		const wsUrl = new SvelteURL('/api/v2/ws', window.location.origin);
		wsUrl.protocol = wsUrl.protocol === 'https:' ? 'wss:' : 'ws:';
		wsUrl.searchParams.set('type', 'article');
		wsUrl.searchParams.set('id', String(postId));

		socket?.close();
		socket = new WebSocket(wsUrl.toString());
		socket.onmessage = (event) => {
			try {
				const payload = JSON.parse(event.data) as PostContentPayload;
				if (!payload?.contentHash) {
					return;
				}
				updateModelData((prev) => {
					if (!prev) {
						return prev;
					}
					return {
						...prev,
						title: payload.title ?? prev.title,
						leadIn: payload.leadIn ?? prev.leadIn,
						toc: payload.toc ?? prev.toc,
						content: payload.content ?? prev.content,
						contentHash: payload.contentHash || prev.contentHash
					};
				});
				triggerUpdateHint();
			} catch {
				// Ignore malformed payloads.
			}
		};
	};

	$effect(() => {
		if (!browser) {
			return;
		}
		const postId = get(postIdStore);
		if (!postId) {
			return;
		}
		connectToPostUpdates(postId);
		return () => {
			socket?.close();
			socket = null;
			if (updateHintTimer) {
				clearTimeout(updateHintTimer);
				updateHintTimer = null;
			}
		};
	});

	$effect(() => {
		if (!browser) {
			return;
		}
		const postId = get(postIdStore);
		if (!postId) {
			return;
		}
		void refreshPostIfNeeded();
	});
</script>

<PostDetail updated={showUpdateHint} />
