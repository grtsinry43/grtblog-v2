import { getPostDetail } from '$lib/infrastructure/post';

export const load = async ({ fetch, params }) => {
	const post = await getPostDetail(fetch, params.id);
	return { post };
};
