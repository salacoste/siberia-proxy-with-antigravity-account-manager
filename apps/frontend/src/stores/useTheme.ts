import { useEffect, useState } from 'react';
import { useConfigStore } from './useConfigStore';

export function useTheme() {
    const { config } = useConfigStore();
    const [resolvedTheme, setResolvedTheme] = useState<'light' | 'dark'>('dark');

    useEffect(() => {
        const themePreference = config?.theme || 'system';
        const root = window.document.documentElement;

        const applyTheme = (theme: 'light' | 'dark') => {
            root.classList.remove('light', 'dark');
            root.classList.add(theme);
            setResolvedTheme(theme);
        };

        if (themePreference === 'system') {
            const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
            applyTheme(mediaQuery.matches ? 'dark' : 'light');

            const listener = (e: MediaQueryListEvent) => {
                applyTheme(e.matches ? 'dark' : 'light');
            };

            mediaQuery.addEventListener('change', listener);
            return () => mediaQuery.removeEventListener('change', listener);
        } else {
            applyTheme(themePreference as 'light' | 'dark');
        }
    }, [config?.theme]);

    return { resolvedTheme };
}
