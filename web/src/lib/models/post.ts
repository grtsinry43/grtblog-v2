export type PostSummary = {
	id: number;
	title: string;
	shortUrl: string;
	summary: string;
	cover?: string | null;
	views: number;
	likes: number;
	comments: number;
	isTop: boolean;
	isHot: boolean;
	isOriginal: boolean;
	createdAt: string;
	updatedAt: string;
};

export type PostDetail = {
	id: number;
	title: string;
	summary: string;
	content: string;
	shortUrl: string;
	cover?: string | null;
	isPublished: boolean;
	isTop: boolean;
	isHot: boolean;
	isOriginal: boolean;
	createdAt: string;
	updatedAt: string;
};

export type PostListResponse = {
	items: PostSummary[];
	total: number;
	page: number;
	size: number;
};
