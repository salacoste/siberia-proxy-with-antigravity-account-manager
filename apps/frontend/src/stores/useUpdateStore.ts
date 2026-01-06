import { create } from 'zustand';
import { CheckForUpdates } from '../../wailsjs/go/main/App';
import { updater } from '../../wailsjs/go/models';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';

interface UpdateState {
    available: boolean;
    checking: boolean;
    updateInfo: updater.UpdateInfo | null;
    error: string | null;

    checkForUpdates: () => Promise<void>;
    openDownloadPage: () => void;
}

export const useUpdateStore = create<UpdateState>((set, get) => ({
    available: false,
    checking: false,
    updateInfo: null,
    error: null,

    checkForUpdates: async () => {
        set({ checking: true, error: null });
        try {
            const info = await CheckForUpdates();
            set({
                checking: false,
                available: info.available,
                updateInfo: info
            });
        } catch (err: any) {
            set({ checking: false, error: err.toString() });
        }
    },

    openDownloadPage: () => {
        const { updateInfo } = get();
        if (updateInfo?.download_url) {
            BrowserOpenURL(updateInfo.download_url);
        }
    }
}));
