import { getPostList } from '$lib/features/post/api';

export const load = async ({ fetch, params, url }) => {
	const rawPage = Number(params.page ?? '1');
	const page = Number.isFinite(rawPage) && rawPage > 0 ? rawPage : 1;
	const rawPageSize = Number(url.searchParams.get('pageSize') ?? '10');
	const pageSize = Number.isFinite(rawPageSize) && rawPageSize > 0 ? rawPageSize : 10;
	const data = await getPostList(fetch, { page, pageSize });
	return { posts: data.items, pagination: { total: data.total, page: data.page, size: data.size } };
};
