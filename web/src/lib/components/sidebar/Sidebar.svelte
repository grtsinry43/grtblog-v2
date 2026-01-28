<script lang="ts">
	import type { NavMenuItem } from '$lib/features/navigation/types';
	import DynamicLucideIcon from '../icons/DynamicLucideIcon.svelte';
	import ThemeIcon from './ThemeIcon.svelte';
    import { page } from '$app/state';


	let { menuTree = [] } = $props<{ menuTree: NavMenuItem[] }>();
    const isActive = (href: string) => page.url.pathname.startsWith(href);
</script>

<div
	class="nav-container w-[7.5em] h-full absolute left-0 top-0 bg-[#FBF9F5] dark:bg-ink-800 flex flex-col items-center pt-6 pb-4 border-r border-ink-200 dark:border-ink-700"
>
	<div class="nav-author-avatar m-8 relative z-0">
		<img
			src="https://dogeoss.grtsinry43.com/img/author.jpeg"
			alt="Author Avatar"
			class="avatar-image w-[53px] h-[53px] object-cover rounded-default z-10"
		/>
	</div>
	<nav class="sidebar-nav">
		<ul>
			{#each menuTree as item}
				<li class="p-2 cursor-pointer hover:bg-ink-600 rounded-md mb-2 {isActive(item.url) ? 'bg-ink-200' : ''}">
					<a href={item.url}>
						{#if item.icon}
							<DynamicLucideIcon name={item.icon} className="nav-icon w-8 h-8 {isActive(item.url) ? 'text-ink-950' : ''}" />
						{/if}
						<span>{item.title}</span>
					</a>
					{#if item.children && item.children.length > 0}
						<ul>
							{#each item.children as child}
								<li>
									<a href={child.url}>
										{#if child.icon}
											<DynamicLucideIcon name={child.icon} className="nav-icon" />
										{/if}
										<span>{child.title}</span>
									</a>
								</li>
							{/each}
						</ul>
					{/if}
				</li>
			{/each}
		</ul>
	</nav>
	<div class="mt-auto pb-2">
		<ThemeIcon />
	</div>
</div>

<style lang="postcss">
	@reference "$routes/layout.css";

	.nav-author-avatar::before {
		content: '';
        @apply w-[53px] h-[53px];
		@apply absolute inset-0 rounded-default border-2 border-ink-200 dark:border-ink-900;
		@apply translate-x-1 translate-y-1 -z-10;
	}
</style>
