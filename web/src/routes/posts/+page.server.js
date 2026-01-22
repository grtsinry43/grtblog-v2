import { getPostList } from '$lib/queries/post';

export const load = async ({ fetch, url }) => {
	const rawPageSize = Number(url.searchParams.get('pageSize') ?? '10');
	const pageSize = Number.isFinite(rawPageSize) && rawPageSize > 0 ? rawPageSize : 10;
	const data = await getPostList(fetch, { page: 1, pageSize });
	return { posts: data.items, pagination: { total: data.total, page: data.page, size: data.size } };
};
