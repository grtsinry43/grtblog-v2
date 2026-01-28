import { createModelDataContext } from 'svatoms';
import type { NavMenuItem } from '$lib/features/navigation/types';
import type { PostSummary } from '$lib/features/post/types';

export type PostListPageData = {
	navMenus: NavMenuItem[];
	posts: PostSummary[];
	pagination: {
		total: number;
		page: number;
		size: number;
	};
};

export const postContext = createModelDataContext<PostListPageData>({
	name: 'postContext',
	initial: null,
})
