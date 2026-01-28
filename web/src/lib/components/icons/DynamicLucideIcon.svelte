<script lang="ts">
	import type { ComponentType } from 'svelte';

	type IconComponent = ComponentType<{ size?: number; strokeWidth?: number; class?: string }>;

	const iconLoaders = import.meta.glob('/node_modules/lucide-svelte/dist/icons/*.svelte');
	const iconCache = new Map<string, IconComponent | null>();

	let {
		name,
		size = 16,
		strokeWidth = 2,
		className = ''
	} = $props<{ name?: string; size?: number; strokeWidth?: number; className?: string }>();

	const toKebab = (value: string) =>
		value
			.replace(/([a-z0-9])([A-Z])/g, '$1-$2')
			.replace(/[_\\s]+/g, '-')
			.toLowerCase();

	const resolveIcon = async (iconName?: string): Promise<IconComponent | null> => {
		if (!iconName) return null;
		if (iconCache.has(iconName)) return iconCache.get(iconName) ?? null;
		const fileName = toKebab(iconName);
		const key = `/node_modules/lucide-svelte/dist/icons/${fileName}.svelte`;
		const loader = iconLoaders[key];
		if (!loader) {
			iconCache.set(iconName, null);
			return null;
		}
		const mod = (await loader()) as { default: IconComponent };
		iconCache.set(iconName, mod.default);
		return mod.default;
	};

	const iconPromise = $derived.by(() => resolveIcon(name));
</script>

{#if name}
	{#await iconPromise then Icon}
		{#if Icon}
			<svelte:component this={Icon} {size} {strokeWidth} class={className} />
		{:else}
			<span class={className} aria-hidden="true"></span>
		{/if}
	{:catch}
		<span class={className} aria-hidden="true"></span>
	{/await}
{/if}
