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
        set({ isLoading: true });
        // @ts-ignore
        if (!window.go) {
            console.warn("Wails runtime not found. Using Mock Data.");
            const mockConfig = new config.AppConfig({
                proxy_port: 7100,
                control_port: 7101,
                config_dir: "/mock/path",
                mitm_enabled: false,
                upstream_proxy: "",
                auth_enabled: false,
                auth_token: "mock-token",
                zai_enabled: false,
                zai_base_url: "https://api.strll.ai/v1",
                zai_api_key: "",
                theme: "dark",
                target_ide: "vscode"
            });
            set({ config: mockConfig, isLoading: false });
            return;
        }

        try {
            const data = await GetAppConfig();
            set({ config: data, isLoading: false });
        } catch (err: any) {
            set({ error: err.toString(), isLoading: false });
        }
    },

    updateConfig: async (newConfig: config.AppConfig) => {
        // Optimistic update
        set({ config: newConfig });

        // @ts-ignore
        if (!window.go) {
            console.warn("Wails runtime not found. Skipping backend update.");
            return;
        }

        try {
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
