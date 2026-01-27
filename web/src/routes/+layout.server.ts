import { fetchNavMenuTree } from '$lib/features/navigation/api';
import type { NavMenuItem } from '$lib/features/navigation/types';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ fetch }) => {
	let navMenus: NavMenuItem[] = [];
	try {
		navMenus = await fetchNavMenuTree(fetch);
	} catch (error) {
		console.error('Failed to load nav menus:', error);
	}

	return {
		navMenus,
	};
};
