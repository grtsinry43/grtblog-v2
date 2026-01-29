import type { ComponentType } from 'svelte';

export type LucideIconComponent = ComponentType<{
	size?: number;
	strokeWidth?: number;
	class?: string;
}>;

// Manual whitelist for tree-shaking in SSR/client bundles.
const lucideLoaders = {
	moon: () => import('lucide-svelte/icons/moon'),
	sun: () => import('lucide-svelte/icons/sun'),
	'book-open': () => import('lucide-svelte/icons/book-open'),
	aperture: () => import('lucide-svelte/icons/aperture'),
	feather: () => import('lucide-svelte/icons/feather'),
	hash: () => import('lucide-svelte/icons/hash'),
	archive: () => import('lucide-svelte/icons/archive'),
	ellipsis: () => import('lucide-svelte/icons/ellipsis')
} as const;

export type LucideIconKey = keyof typeof lucideLoaders;
export type LucideIconLoader = (typeof lucideLoaders)[LucideIconKey];

export default lucideLoaders;
