export type DiscoveryEntry = {
	title: string;
	url: string;
	path: string;
	kind: string;
	summary?: string;
	modified?: string;
	markdownUrl?: string;
};

export type DiscoveryCatalog = {
	siteName: string;
	baseUrl: string;
	entries: DiscoveryEntry[];
};
