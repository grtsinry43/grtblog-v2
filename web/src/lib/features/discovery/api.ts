import { getApi } from '$lib/shared/clients/api';
import type { DiscoveryCatalog } from './types';

export const getDiscoveryCatalog = (fetch: typeof globalThis.fetch) =>
	getApi(fetch)<DiscoveryCatalog>('/public/discovery/catalog');
