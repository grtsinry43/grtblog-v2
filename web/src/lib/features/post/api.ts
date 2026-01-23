import { getApi } from '$lib/shared/clients/api';
import type { PostDetail, PostLatestCheckResponse, PostListResponse } from '$lib/features/post/types';

type PostListOptions = {
	page?: number;
	pageSize?: number;
};

export const getPostList = async (
	fetcher?: typeof fetch,
	{ page = 1, pageSize = 10 }: PostListOptions = {}
): Promise<PostListResponse> => {
	const api = getApi(fetcher);
	const query = new URLSearchParams({
		page: String(page),
		pageSize: String(pageSize)
	});
	const result = await api<PostListResponse>(`/articles?${query.toString()}`);
	return result ?? { items: [], total: 0, page, size: pageSize };
};

export const getPostDetail = async (
	fetcher: typeof fetch | undefined,
	shortUrl: string
): Promise<PostDetail | null> => {
	const api = getApi(fetcher);
	const result = await api<PostDetail>(`/articles/short/${shortUrl}`);
	return result ?? null;
};

export const checkPostLatest = async (
	fetcher: typeof fetch | undefined,
	id: number,
	hash: string
): Promise<PostLatestCheckResponse | null> => {
	const api = getApi(fetcher);
	const result = await api<PostLatestCheckResponse>(`/articles/${id}/latest`, {
		method: 'POST',
		body: { hash }
	});
	return result ?? null;
};
