import { useState, useEffect } from 'react';
import { ThemeToggle } from '@/components/ThemeToggle';
import { SyncSettings } from "@/components/settings/SyncSettings";
import { useConfigStore } from '@/stores/useConfigStore';
import { toast } from "sonner";
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Switch } from '@/components/ui/switch';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Shield } from 'lucide-react';

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

    const handleSaveUpstream = () => {
        if (config) {
            // @ts-ignore
            updateConfig({ ...config, upstream_proxy: upstream });
            toast.success("Settings saved");
        }
    };

    const handleInstallCert = async () => {
        try {
            // @ts-ignore
            await window.runtime.Invoke("InstallCert");
            toast.success("Certificate installed successfully.");
        } catch (error: any) {
            toast.error("Failed to install certificate: " + error);
        }
    };

    // Helper for config updates
    const handleConfigChange = (key: string, value: any) => {
        if (config) {
            // @ts-ignore
            updateConfig({ ...config, [key]: value });
        }
    };

    const handleCheckUpdates = async () => {
        setCheckingUpdate(true);
        try {
            const info = await CheckForUpdates();
            if (info.available) {
                toast.promise(
                    new Promise((resolve) => {
                        resolve(true);
                    }),
                    {
                        loading: 'Update found!',
                        success: (
                            <div className="flex flex-col gap-2">
                                <span className="font-bold">Update {info.latest_version} available!</span>
                                <Button size="sm" onClick={() => BrowserOpenURL(info.download_url)}>
                                    Download Update
                                </Button>
                            </div>
                        ),
                        error: 'Error'
                    }
                );
            } else {
                if (info.error) {
                    toast.error(`Error checking updates: ${info.error} `);
                } else {
                    toast.success(`You are on the latest version(${info.current_version})`);
                }
            }
        } catch (error) {
            console.error(error);
            toast.error("Failed to check for updates");
        } finally {
            setCheckingUpdate(false);
        }
    };

    if (!config) return <div className="p-8">Loading configuration...</div>;

    return (
        <div className="p-8 space-y-6">
            <h1 className="text-3xl font-bold">Settings</h1>

            {/* Theme */}
            <Card>
                <CardHeader>
                    <CardTitle>Appearance</CardTitle>
                    <CardDescription>Customize application theme</CardDescription>
                </CardHeader>
                <CardContent className="flex justify-between items-center">
                    <span className="text-sm">Select Theme</span>
                    <ThemeToggle />
                </CardContent>
            </Card>

            {/* Target IDE */}
            <Card>
                <CardHeader>
                    <CardTitle>Target IDE</CardTitle>
                    <CardDescription>Select which IDE you are using (VS Code, Cursor, Windsurf)</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="flex items-center justify-between">
                        <Label>Current IDE</Label>
                        <select
                            className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 w-[200px]"
                            // @ts-ignore
                            value={config.target_ide || "vscode"}
                            onChange={(e) => handleConfigChange("target_ide", e.target.value)}
                        >
                            <option value="vscode">VS Code</option>
                            <option value="cursor">Cursor</option>
                            <option value="windsurf">Windsurf</option>
                        </select>
                    </div>
                </CardContent>
            </Card>

            {/* HTTPS Inspection */}
            <Card>
                <CardHeader>
                    <CardTitle className="flex items-center gap-2">
                        <Shield className="h-5 w-5" />
                        HTTPS Inspection
                    </CardTitle>
                    <CardDescription>
                        Enable decryption of HTTPS traffic to see full request details.
                    </CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="flex items-center justify-between">
                        <div className="space-y-0.5">
                            <Label>Certificate Trust</Label>
                            <div className="text-sm text-muted-foreground">
                                Root CA must be trusted by the OS.
                            </div>
                        </div>
                        <Button variant="outline" onClick={handleInstallCert}>
                            Install Certificate
                        </Button>
                    </div>
                    <div className="flex items-center justify-between">
                        <div className="space-y-0.5">
                            <Label>Enable Decryption</Label>
                            <div className="text-sm text-muted-foreground">
                                Decrypts traffic for inspection.
                            </div>
                        </div>
                        <Switch
                            checked={config.mitm_enabled}
                            onCheckedChange={(checked) => handleConfigChange("mitm_enabled", checked)}
                        />
                    </div>
                </CardContent>
            </Card>

            {/* Upstream Proxy */}
            <Card>
                <CardHeader>
                    <CardTitle>Upstream Proxy</CardTitle>
                    <CardDescription>Chain all traffic through another proxy</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="flex gap-2">
                        <Input
                            placeholder="http://127.0.0.1:8888"
                            value={upstream}
                            onChange={(e) => setUpstream(e.target.value)}
                        />
                        <Button onClick={handleSaveUpstream}>Save</Button>
                    </div>
                </CardContent>
            </Card>

            {/* z.ai Provider */}
            <Card>
                <CardHeader>
                    <CardTitle>z.ai Provider</CardTitle>
                    <CardDescription>Use z.ai as the LLM backend</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="flex items-center justify-between">
                        <Label>Enable Provider</Label>
                        <Switch
                            checked={config.zai_enabled}
                            onCheckedChange={(c) => handleConfigChange("zai_enabled", c)}
                        />
                    </div>

                    {config.zai_enabled && (
                        <div className="space-y-4 pt-4 border-t">
                            <div className="space-y-2">
                                <Label>Base URL</Label>
                                <Input
                                    placeholder="https://api.z.ai/v1"
                                    value={config.zai_base_url}
                                    onChange={(e) => handleConfigChange("zai_base_url", e.target.value)}
                                />
                            </div>
                            <div className="space-y-2">
                                <Label>API Key</Label>
                                <Input
                                    type="password"
                                    placeholder="sk-..."
                                    value={config.zai_api_key}
                                    onChange={(e) => handleConfigChange("zai_api_key", e.target.value)}
                                />
                            </div>
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Proxy Security */}
            <Card>
                <CardHeader>
                    <CardTitle>Proxy Security</CardTitle>
                    <CardDescription>Protect access to the proxy with a token</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                    <div className="flex items-center justify-between">
                        <Label>Enable Authentication</Label>
                        <Switch
                            checked={config.auth_enabled}
                            onCheckedChange={(c) => handleConfigChange("auth_enabled", c)}
                        />
                    </div>

                    {config.auth_enabled && (
                        <div className="space-y-4 pt-4 border-t">
                            <div className="space-y-2">
                                <Label>Access Token</Label>
                                <div className="flex gap-2">
                                    <Input
                                        className="font-mono"
                                        placeholder="Enter or generate a token"
                                        value={config.auth_token}
                                        onChange={(e) => handleConfigChange("auth_token", e.target.value)}
                                    />
                                    <Button variant="secondary" onClick={() => handleConfigChange("auth_token", 'sk-sib-' + Math.random().toString(36).substring(2, 15))}>
                                        Generate
                                    </Button>
                                </div>
                            </div>
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Cloud Sync */}
            <SyncSettings />

            {/* Application Info */}
            <Card>
                <CardHeader>
                    <CardTitle>Application Info</CardTitle>
                </CardHeader>
                <CardContent className="flex items-center justify-between">
                    <div className="text-sm">
                        <span className="text-muted-foreground mr-2">Current Version:</span>
                        <span className="font-mono font-bold">{version}</span>
                    </div>
                    <Button onClick={handleCheckUpdates} disabled={checkingUpdate} variant="outline">
                        {checkingUpdate ? "Checking..." : "Check for Updates"}
                    </Button>
                </CardContent>
            </Card>
        </div>
    );
}
