import { browser } from '$app/environment';

export type Theme = 'light' | 'dark' | 'system';

class ThemeManager {
    current = $state<Theme>('system');

    constructor() {
        if (browser) {
            const saved = localStorage.getItem('theme') as Theme;
            if (saved) {
                this.current = saved;
            }
        }

        $effect.root(() => {
            $effect(() => {
                if (!browser) return;

                const isDark =
                    this.current === 'dark' ||
                    (this.current === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches);

                document.documentElement.classList.toggle('dark', isDark);
                localStorage.setItem('theme', this.current);
            });
        });
    }

    set(theme: Theme) {
        this.current = theme;
    }
}

export const themeManager = new ThemeManager();
