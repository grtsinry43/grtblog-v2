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
	contentHash: string;
	leadIn?: string | null;
	toc?: TOCNode[];
	shortUrl: string;
	cover?: string | null;
	isPublished: boolean;
	metrics: {
		views: number;
		likes: number;
		comments: number;
	};
	isTop: boolean;
	isHot: boolean;
	isOriginal: boolean;
	createdAt: string;
	updatedAt: string;
};

export type TOCNode = {
	name: string;
	anchor: string;
	children?: TOCNode[];
};

export type PostLatestCheckResponse = {
	latest: boolean;
	contentHash: string;
	title?: string;
	leadIn?: string | null;
	toc?: TOCNode[];
	content?: string;
};

export type PostContentPayload = {
	contentHash: string;
	title?: string;
	leadIn?: string | null;
	toc?: TOCNode[];
	content?: string;
};

export type PostListResponse = {
	items: PostSummary[];
	total: number;
	page: number;
	size: number;
};
