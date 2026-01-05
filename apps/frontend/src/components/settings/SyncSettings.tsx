import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { RefreshCw, UploadCloud, DownloadCloud, Lock, LogIn, UserPlus, ShieldCheck } from 'lucide-react';
import { toast } from 'sonner';

// Type definition for the window object to include the Go backend
declare global {
    interface Window {
        go: {
            main: {
                App: {
                    SyncPush: (password: string) => Promise<void>;
                    SyncPull: (password: string) => Promise<string>;
                    SyncSignUp: (email: string, password: string) => Promise<void>;
                    SyncSignIn: (email: string, password: string) => Promise<void>;
                    SyncGetUser: () => Promise<string>;
                }
            }
        }
    }
}

export function SyncSettings() {
    const [password, setPassword] = useState("");
    const [isSyncing, setIsSyncing] = useState(false);
    const [lastSync, setLastSync] = useState<string | null>(null);
    const [userId, setUserId] = useState<string>("");

    // Auth State
    const [email, setEmail] = useState("");
    const [authPassword, setAuthPassword] = useState("");

    useEffect(() => {
        checkAuth();
    }, []);

    const checkAuth = async () => {
        // @ts-ignore
        if (!window.go) return;
        try {
            const id = await window.go.main.App.SyncGetUser();
            setUserId(id);
        } catch (e) {
            console.log("Not logged in");
        }
    };

    const handleLogin = async () => {
        // @ts-ignore
        if (!window.go) {
            toast.error("Not available in Web Mode");
            return;
        }
        try {
            await window.go.main.App.SyncSignIn(email, authPassword);
            toast.success("Logged in successfully");
            checkAuth();
        } catch (e: any) {
            toast.error("Login failed: " + e);
        }
    };

    const handleSignup = async () => {
        // @ts-ignore
        if (!window.go) {
            toast.error("Not available in Web Mode");
            return;
        }
        try {
            await window.go.main.App.SyncSignUp(email, authPassword);
            toast.success("Account created! You are now logged in.");
            checkAuth();
        } catch (e: any) {
            toast.error("Signup failed: " + e);
        }
    };

    const handlePush = async () => {
        // @ts-ignore
        if (!window.go) {
            toast.error("Not available in Web Mode");
            return;
        }
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
        // @ts-ignore
        if (!window.go) {
            toast.error("Not available in Web Mode");
            return;
        }
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

    if (!userId) {
        return (
            <Card className="w-full max-w-md mx-auto mt-6">
                <CardHeader>
                    <CardTitle>Cloud Sync Login</CardTitle>
                    <CardDescription>Sign in to sync your profiles across devices.</CardDescription>
                </CardHeader>
                <CardContent>
                    <Tabs defaultValue="login" className="w-full">
                        <TabsList className="grid w-full grid-cols-2">
                            <TabsTrigger value="login">Login</TabsTrigger>
                            <TabsTrigger value="signup">Sign Up</TabsTrigger>
                        </TabsList>

                        <TabsContent value="login" className="space-y-4 mt-4">
                            <div className="space-y-2">
                                <Label>Email</Label>
                                <Input value={email} onChange={(e) => setEmail(e.target.value)} type="email" />
                            </div>
                            <div className="space-y-2">
                                <Label>Password</Label>
                                <Input value={authPassword} onChange={(e) => setAuthPassword(e.target.value)} type="password" />
                            </div>
                            <Button className="w-full" onClick={handleLogin}>
                                <LogIn className="mr-2 h-4 w-4" /> Sign In
                            </Button>
                        </TabsContent>

                        <TabsContent value="signup" className="space-y-4 mt-4">
                            <div className="space-y-2">
                                <Label>Email</Label>
                                <Input value={email} onChange={(e) => setEmail(e.target.value)} type="email" />
                            </div>
                            <div className="space-y-2">
                                <Label>Password</Label>
                                <Input value={authPassword} onChange={(e) => setAuthPassword(e.target.value)} type="password" />
                            </div>
                            <Button className="w-full" onClick={handleSignup}>
                                <UserPlus className="mr-2 h-4 w-4" /> Create Account
                            </Button>
                        </TabsContent>
                    </Tabs>
                </CardContent>
            </Card>
        );
    }

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

                <div className="flex items-center justify-between p-4 border rounded-lg bg-green-500/10 border-green-500/20">
                    <div className="flex items-center gap-3">
                        <ShieldCheck className="h-5 w-5 text-green-600" />
                        <div className="flex flex-col">
                            <span className="font-medium text-sm text-green-700">Authenticated</span>
                            <span className="text-xs text-green-600/80">User ID: {userId.slice(0, 8)}...</span>
                        </div>
                    </div>
                    <Button variant="ghost" size="sm" onClick={() => setUserId("")} className="text-xs">
                        Logout
                    </Button>
                </div>

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
                <ShieldCheck className="h-3 w-3 text-green-500" />
                <span>End-to-End Encrypted. The server cannot see your data.</span>
            </CardFooter>
        </Card>
    );
}
