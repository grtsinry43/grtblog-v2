import { getDiscoveryCatalog } from '$lib/features/discovery/api';
import type { PageServerLoad } from './$types';

export const prerender = false;

export const load: PageServerLoad = async ({ fetch, setHeaders }) => {
	setHeaders({ 'cache-control': 'no-store' });
	return { catalog: await getDiscoveryCatalog(fetch) };
};
