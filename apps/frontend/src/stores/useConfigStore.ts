import { create } from 'zustand';
import { GetAppConfig, UpdateAppConfig } from '../../wailsjs/go/main/App';
import { config } from '../../wailsjs/go/models';

interface ConfigState {
    config: config.AppConfig | null;
    isLoading: boolean;
    error: string | null;

    // Actions
    fetchConfig: () => Promise<void>;
    updateConfig: (newConfig: config.AppConfig) => Promise<void>;
    setTheme: (theme: string) => Promise<void>;
}

export const useConfigStore = create<ConfigState>((set, get) => ({
    config: null,
    isLoading: true,
    error: null,

    fetchConfig: async () => {
        try {
            set({ isLoading: true });
            const data = await GetAppConfig();
            set({ config: data, isLoading: false });
        } catch (err: any) {
            set({ error: err.toString(), isLoading: false });
        }
    },

    updateConfig: async (newConfig: config.AppConfig) => {
        try {
            // Optimistic update
            set({ config: newConfig });
            await UpdateAppConfig(newConfig);
        } catch (err: any) {
            // Revert on failure (reload from backend)
            get().fetchConfig();
            set({ error: err.toString() });
        }
    },

    setTheme: async (theme: string) => {
        const current = get().config;
        if (!current) return;

        const newConfig = { ...current, theme };
        await get().updateConfig(newConfig);
    },
}));
