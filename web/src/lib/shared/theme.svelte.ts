import { browser } from '$app/environment';

export type Theme = 'light' | 'dark' | 'system';
export type ResolvedTheme = 'light' | 'dark';

class ThemeManager {
	current = $state<Theme>('system');

	set(theme: Theme) {
		this.current = theme;
	}
}

export const themeManager = new ThemeManager();

export const resolveTheme = (theme: Theme): ResolvedTheme => {
	if (!browser || theme !== 'system') {
		return theme === 'dark' ? 'dark' : 'light';
	}

	return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
};

export const initTheme = (manager: ThemeManager): void => {
	if (!browser) return;

	const saved = localStorage.getItem('theme') as Theme | null;
	if (saved === 'light' || saved === 'dark' || saved === 'system') {
		manager.set(saved);
	}
};

export const startThemeSync = (manager: ThemeManager): void => {
	$effect(() => {
		if (!browser) return;

		const media = window.matchMedia('(prefers-color-scheme: dark)');
		const apply = () => {
			const resolved = resolveTheme(manager.current);
			document.documentElement.classList.toggle('dark', resolved === 'dark');
			localStorage.setItem('theme', manager.current);
		};

		apply();

		if (manager.current !== 'system') return;

		const onChange = () => apply();
		media.addEventListener('change', onChange);
		return () => media.removeEventListener('change', onChange);
	});
};
