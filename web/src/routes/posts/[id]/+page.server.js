import { getPostDetail } from '$lib/features/post/api';

export const load = async ({ fetch, params }) => {
	const post = await getPostDetail(fetch, params.id);
	return { post };
};
