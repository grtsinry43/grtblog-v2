import { getPostDetail } from '$lib/queries/post';

export const load = async ({ fetch, params }) => {
	const post = await getPostDetail(fetch, params.id);
	return { post };
};
