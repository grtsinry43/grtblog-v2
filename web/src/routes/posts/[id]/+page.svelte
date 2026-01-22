<script lang="ts">
	import { get, writable } from 'svelte/store';
	import { PostDetail } from '$lib/ui/post';
	import { checkPostLatest } from '$lib/infrastructure/post';
	import type { PostContentPayload, PostDetail as PostDetailModel } from '$lib/domain/post';
	import { browser } from '$app/environment';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
	const postStore = writable<PostDetailModel | null>(null);
	let socket: WebSocket | null = null;
	let showUpdateHint = $state(false);
	let updateHintTimer: ReturnType<typeof setTimeout> | null = null;

	$effect(() => {
		postStore.set(data.post ?? null);
	});

	const refreshPostIfNeeded = async () => {
		const current = get(postStore);
		if (!current?.id || !current.contentHash) {
			return;
		}

		try {
			const latest = await checkPostLatest(undefined, current.id, current.contentHash);
			if (!latest || latest.latest) {
				return;
			}

			postStore.update((prev) => {
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
		showUpdateHint = true;
		if (updateHintTimer) {
			clearTimeout(updateHintTimer);
		}
		updateHintTimer = setTimeout(() => {
			showUpdateHint = false;
			updateHintTimer = null;
		}, 2400);
	};

	const connectToPostUpdates = (postId: number) => {
		if (!browser) {
			return;
		}

		const wsUrl = new URL('/api/v2/ws', window.location.origin);
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
				postStore.update((prev) => {
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
		const current = data.post;
		if (!current?.id) {
			return;
		}
		connectToPostUpdates(current.id);
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
		const postId = data.post?.id;
		if (!postId) {
			return;
		}
		void refreshPostIfNeeded();
	});
</script>

<PostDetail post={$postStore} updated={showUpdateHint} />
