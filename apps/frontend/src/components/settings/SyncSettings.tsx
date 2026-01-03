import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RefreshCw, UploadCloud, DownloadCloud, Lock, AlertTriangle } from 'lucide-react';
import { toast } from 'sonner';

// Type definition for the window object to include the Go backend
declare global {
    interface Window {
        go: {
            main: {
                App: {
                    SyncPush: (password: string) => Promise<void>;
                    SyncPull: (password: string) => Promise<string>;
                }
            }
        }
    }
}

export function SyncSettings() {
    const [password, setPassword] = useState("");
    const [isSyncing, setIsSyncing] = useState(false);
    const [lastSync, setLastSync] = useState<string | null>(null);

    const handlePush = async () => {
        if (!password) {
            toast.error("Master Password required for encryption");
            return;
        }

        setIsSyncing(true);
        try {
            await window.go.main.App.SyncPush(password);
            toast.success("Profile pushed to cloud successfully");
            setLastSync(new Date().toLocaleTimeString());
        } catch (err: any) {
            console.error(err);
            toast.error("Sync Push Failed: " + err);
        } finally {
            setIsSyncing(false);
        }
    };

    const handlePull = async () => {
        if (!password) {
            toast.error("Master Password required for decryption");
            return;
        }

        setIsSyncing(true);
        try {
            const data = await window.go.main.App.SyncPull(password);
            toast.success("Profile pulled from cloud");
            console.log("Decrypted Data:", data); // For MVP debug
            setLastSync(new Date().toLocaleTimeString());
        } catch (err: any) {
            console.error(err);
            toast.error("Sync Pull Failed: " + err);
        } finally {
            setIsSyncing(false);
        }
    };

    return (
        <Card className="w-full max-w-2xl mx-auto mt-6">
            <CardHeader>
                <div className="flex items-center gap-2">
                    <RefreshCw className={`w-6 h-6 ${isSyncing ? 'animate-spin' : ''} text-primary`} />
                    <CardTitle>Cloud Profile Sync (Zero-Knowledge)</CardTitle>
                </div>
                <CardDescription>
                    Sync your accounts and settings securely across devices.
                    Your data is encrypted with your Master Password before leaving this device.
                </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">

                <div className="space-y-2">
                    <Label htmlFor="master-password">Master Password (Encryption Key)</Label>
                    <div className="flex gap-2">
                        <div className="relative flex-1">
                            <Lock className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                            <Input
                                id="master-password"
                                type="password"
                                placeholder="Enter secure password..."
                                className="pl-9"
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                            />
                        </div>
                    </div>
                    <p className="text-xs text-muted-foreground">
                        Your password is never sent to the server. It is only used to derive encryption keys locally.
                    </p>
                </div>

                <div className="flex items-center justify-between p-4 border rounded-lg bg-muted/30">
                    <div className="flex items-center gap-3">
                        <div className="h-2 w-2 rounded-full bg-green-500"></div>
                        <div className="flex flex-col">
                            <span className="font-medium text-sm">Sync Status: Idle</span>
                            <span className="text-xs text-muted-foreground">
                                {lastSync ? `Last synced at ${lastSync}` : 'Not synced yet'}
                            </span>
                        </div>
                    </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                    <Button
                        variant="outline"
                        className="h-24 flex flex-col gap-2 hover:border-primary/50 hover:bg-primary/5"
                        onClick={handlePush}
                        disabled={isSyncing}
                    >
                        <UploadCloud className="h-8 w-8 text-blue-500" />
                        <span className="font-bold">Push to Cloud</span>
                        <span className="text-xs text-muted-foreground font-normal">Encrypt & Upload</span>
                    </Button>

                    <Button
                        variant="outline"
                        className="h-24 flex flex-col gap-2 hover:border-primary/50 hover:bg-primary/5"
                        onClick={handlePull}
                        disabled={isSyncing}
                    >
                        <DownloadCloud className="h-8 w-8 text-green-500" />
                        <span className="font-bold">Pull from Cloud</span>
                        <span className="text-xs text-muted-foreground font-normal">Download & Decrypt</span>
                    </Button>
                </div>

            </CardContent>
            <CardFooter className="bg-muted/10 p-4 rounded-b-lg border-t text-xs text-muted-foreground flex items-center gap-2">
                <AlertTriangle className="h-3 w-3 text-yellow-500" />
                <span>MVP Warning: "Mock Login" is active. Conflict resolution is strictly "Server Wins".</span>
            </CardFooter>
        </Card>
    );
}
