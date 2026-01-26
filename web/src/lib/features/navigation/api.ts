import { getApi } from '$lib/shared/clients/api';
import type { NavMenuItem } from './types';

export async function fetchNavMenuTree(fetcher?: typeof fetch): Promise<NavMenuItem[]> {
	const api = getApi(fetcher);
	return api<NavMenuItem[]>('/public/nav-menus');
}
