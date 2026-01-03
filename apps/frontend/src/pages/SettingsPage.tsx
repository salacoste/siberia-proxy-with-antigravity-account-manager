import { useState, useEffect } from 'react';
import { ThemeToggle } from '@/components/ThemeToggle';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { useConfigStore } from '@/stores/useConfigStore';
import { toast } from "sonner";

import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';
import { CheckForUpdates, GetVersion } from '../../wailsjs/go/main/App';

export default function SettingsPage() {
    const { config, updateConfig } = useConfigStore();
    const [upstream, setUpstream] = useState('');
    const [version, setVersion] = useState('Loading...');
    const [checkingUpdate, setCheckingUpdate] = useState(false);

    useEffect(() => {
        if (config) {
            setUpstream(config.upstream_proxy || '');
        }
        GetVersion().then(setVersion);
    }, [config]);

    // ... imports

    const handleSaveUpstream = () => {
        if (config) {
            // @ts-ignore
            updateConfig({ ...config, upstream_proxy: upstream });
            toast.success("Settings saved");
        }
    };

    const handleCheckUpdates = async () => {
        setCheckingUpdate(true);
        try {
            const info = await CheckForUpdates();
            if (info.available) {
                toast.promise(
                    new Promise((resolve) => {
                        // We just show the toast, the button action is below
                        resolve(true);
                    }),
                    {
                        loading: 'Update found!',
                        success: (
                            <div className="flex flex-col gap-2">
                                <span className="font-bold">Update {info.version} available!</span>
                                <Button size="sm" onClick={() => BrowserOpenURL(info.download_url)}>
                                    Download Update
                                </Button>
                            </div>
                        ),
                        error: 'Error'
                    }
                );
            } else {
                toast.success(`You are on the latest version (${info.version})`);
            }
        } catch (error) {
            console.error(error);
            toast.error("Failed to check for updates");
        } finally {
            setCheckingUpdate(false);
        }
    };

    return (
        <div className="p-8 space-y-6">
            <h1 className="text-3xl font-bold">Settings</h1>

            {/* Theme */}
            <div className="flex items-center justify-between p-4 border rounded-lg bg-card text-card-foreground">
                <div>
                    <h2 className="font-semibold">Appearance</h2>
                    <p className="text-sm text-muted-foreground">Customize application theme</p>
                </div>
                <ThemeToggle />
            </div>

            {/* Upstream Proxy */}
            <div className="p-4 border rounded-lg bg-card text-card-foreground space-y-4">
                <div>
                    <h2 className="font-semibold">Upstream Proxy</h2>
                    <p className="text-sm text-muted-foreground">Chain all traffic through another proxy (e.g., http://user:pass@host:port)</p>
                </div>
                <div className="flex gap-2">
                    <Input
                        placeholder="http://127.0.0.1:8888"
                        value={upstream}
                        onChange={(e) => setUpstream(e.target.value)}
                    />
                    <Button onClick={handleSaveUpstream}>Save</Button>
                </div>
            </div>

            {/* z.ai Provider */}
            <div className="p-4 border rounded-lg bg-card text-card-foreground space-y-4">
                <div className="flex items-center justify-between">
                    <div>
                        <h2 className="font-semibold">z.ai Provider</h2>
                        <p className="text-sm text-muted-foreground">Example: Use z.ai as the LLM backend</p>
                    </div>
                    <Button
                        variant={config?.zai_enabled ? "default" : "outline"}
                        onClick={() => config && updateConfig({ ...config, zai_enabled: !config.zai_enabled })}
                    >
                        {config?.zai_enabled ? "Enabled" : "Disabled"}
                    </Button>
                </div>

                {config?.zai_enabled && (
                    <div className="space-y-4 pt-4 border-t">
                        <div className="space-y-2">
                            <label className="text-sm font-medium">Base URL</label>
                            <Input
                                placeholder="https://api.z.ai/v1"
                                value={config.zai_base_url}
                                onChange={(e) => updateConfig({ ...config, zai_base_url: e.target.value })}
                            />
                        </div>
                        <div className="space-y-2">
                            <label className="text-sm font-medium">API Key</label>
                            <Input
                                type="password"
                                placeholder="sk-..."
                                value={config.zai_api_key}
                                onChange={(e) => updateConfig({ ...config, zai_api_key: e.target.value })}
                            />
                        </div>
                    </div>
                )}
            </div>

            {/* Proxy Security */}
            <div className="p-4 border rounded-lg bg-card text-card-foreground space-y-4">
                <div className="flex items-center justify-between">
                    <div>
                        <h2 className="font-semibold">Proxy Security</h2>
                        <p className="text-sm text-muted-foreground">Protect access to the proxy with a token</p>
                    </div>
                    <Button
                        variant={config?.auth_enabled ? "default" : "outline"}
                        onClick={() => config && updateConfig({ ...config, auth_enabled: !config.auth_enabled })}
                    >
                        {config?.auth_enabled ? "Enabled" : "Disabled"}
                    </Button>
                </div>

                {config?.auth_enabled && (
                    <div className="space-y-4 pt-4 border-t">
                        <div className="space-y-2">
                            <label className="text-sm font-medium">Access Token</label>
                            <div className="flex gap-2">
                                <Input
                                    className="font-mono"
                                    placeholder="Enter or generate a token"
                                    value={config.auth_token}
                                    onChange={(e) => updateConfig({ ...config, auth_token: e.target.value })}
                                />
                                <Button variant="secondary" onClick={() => config && updateConfig({ ...config, auth_token: 'sk-sib-' + Math.random().toString(36).substring(2, 15) })}>
                                    Generate
                                </Button>
                            </div>
                        </div>
                    </div>
                )}
            </div>

            {/* Application Info */}
            <div className="p-4 border rounded-lg bg-card text-card-foreground space-y-4">
                <div>
                    <h2 className="font-semibold">Application Info</h2>
                    <p className="text-sm text-muted-foreground">Version and Updates</p>
                </div>
                <div className="flex items-center justify-between">
                    <div className="text-sm">
                        <span className="text-muted-foreground">Current Version: </span>
                        <span className="font-mono font-bold">{version}</span>
                    </div>
                    <Button onClick={handleCheckUpdates} disabled={checkingUpdate}>
                        {checkingUpdate ? "Checking..." : "Check for Updates"}
                    </Button>
                </div>
            </div>
        </div>
    );
}
