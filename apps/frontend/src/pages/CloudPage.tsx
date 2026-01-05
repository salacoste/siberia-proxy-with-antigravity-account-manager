import { useState, useEffect } from "react";
// @ts-ignore
import { CloudLogin, CloudSignUp, CloudGetStatus, CloudLogout, CloudSync } from "../../wailsjs/go/main/App";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Loader2, Cloud, LogOut, RefreshCw } from "lucide-react";
import { toast } from "sonner";

export default function CloudPage() {
    const [status, setStatus] = useState({ enabled: false, email: "", last_sync: "" });
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        refreshStatus();
    }, []);

    const refreshStatus = async () => {
        // @ts-ignore
        if (!window.go) {
            console.warn("Web Mode: Cloud Sync unavailable");
            return;
        }
        try {
            const s = await CloudGetStatus();
            setStatus(s);
        } catch (e) {
            console.error(e);
        }
    };

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");

    const handleLogin = async () => {
        // @ts-ignore
        if (!window.go) { toast.error("Unavailable in Web Mode"); return; }
        setLoading(true);
        try {
            await CloudLogin(email, password);
            toast.success("Logged in successfully");
            await refreshStatus();
        } catch (e: any) {
            toast.error("Login failed: " + e);
        } finally {
            setLoading(false);
        }
    };

    const handleSignUp = async () => {
        // @ts-ignore
        if (!window.go) { toast.error("Unavailable in Web Mode"); return; }
        setLoading(true);
        try {
            await CloudSignUp(email, password);
            toast.success("Account created! You are now logged in.");
            await refreshStatus();
        } catch (e: any) {
            toast.error("Sign up failed: " + e);
        } finally {
            setLoading(false);
        }
    };

    const handleLogout = async () => {
        // @ts-ignore
        if (!window.go) { toast.error("Unavailable in Web Mode"); return; }
        setLoading(true);
        try {
            await CloudLogout();
            toast.info("Logged out");
            await refreshStatus();
        } catch (e: any) {
            toast.error("Logout failed: " + e);
        } finally {
            setLoading(false);
        }
    };

    const handleSync = async () => {
        // @ts-ignore
        if (!window.go) { toast.error("Unavailable in Web Mode"); return; }
        setLoading(true);
        try {
            await CloudSync();
            toast.success("Sync completed!");
            await refreshStatus();
        } catch (e: any) {
            toast.error("Sync failed: " + e);
        } finally {
            setLoading(false);
        }
    };

    if (status.enabled) {
        return (
            <div className="flex flex-col items-center justify-center p-10 h-full space-y-6">
                <Card className="w-[400px] border-emerald-500/20 bg-emerald-950/10 backdrop-blur">
                    <CardHeader className="text-center">
                        <div className="mx-auto bg-emerald-500/20 p-3 rounded-full mb-2 w-fit">
                            <Cloud className="w-8 h-8 text-emerald-400" />
                        </div>
                        <CardTitle className="text-xl">Cloud Sync Active</CardTitle>
                        <CardDescription>Your traffic profiles are synchronized.</CardDescription>
                    </CardHeader>
                    <CardContent className="space-y-4">
                        <div className="flex items-center justify-between p-3 bg-secondary/50 rounded-lg">
                            <span className="text-sm text-muted-foreground">Account</span>
                            <span className="font-medium text-sm">{status.email}</span>
                        </div>
                        <div className="flex items-center justify-between p-3 bg-secondary/50 rounded-lg">
                            <span className="text-sm text-muted-foreground">Last Sync</span>
                            <span className="font-medium text-sm">
                                {status.last_sync ? new Date(status.last_sync).toLocaleString() : "Never"}
                            </span>
                        </div>
                    </CardContent>
                    <CardFooter className="flex gap-2">
                        <Button variant="outline" className="flex-1" onClick={handleLogout} disabled={loading}>
                            <LogOut className="w-4 h-4 mr-2" />
                            Logout
                        </Button>
                        <Button className="flex-1 bg-emerald-600 hover:bg-emerald-700" onClick={handleSync} disabled={loading}>
                            {loading ? <Loader2 className="w-4 h-4 animate-spin mr-2" /> : <RefreshCw className="w-4 h-4 mr-2" />}
                            Sync Now
                        </Button>
                    </CardFooter>
                </Card>
            </div>
        );
    }

    return (
        <div className="flex flex-col items-center justify-center p-10 h-full animate-in fade-in zoom-in-95 duration-500">
            <div className="mb-8 text-center space-y-2">
                <h1 className="text-3xl font-bold tracking-tight bg-gradient-to-r from-blue-400 to-emerald-400 bg-clip-text text-transparent">
                    Siberia Cloud
                </h1>
                <p className="text-muted-foreground">Sync your settings, profiles, and interceptors across devices.</p>
                {/* @ts-ignore */}
                {!window.go && (
                    <div className="mx-auto mt-4 px-3 py-1 text-xs font-mono bg-yellow-500/10 text-yellow-500 rounded border border-yellow-500/20 w-fit">
                        WEB MODE: Unavailable
                    </div>
                )}
            </div>

            <Tabs defaultValue="login" className="w-[400px]">
                <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="login">Login</TabsTrigger>
                    <TabsTrigger value="signup">Register</TabsTrigger>
                </TabsList>

                <TabsContent value="login">
                    <Card>
                        <CardHeader>
                            <CardTitle>Welcome back</CardTitle>
                            <CardDescription>Enter your credentials to access your cloud profile.</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="space-y-2">
                                <label className="text-sm font-medium">Email</label>
                                <Input placeholder="user@example.com" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
                            </div>
                            <div className="space-y-2">
                                <label className="text-sm font-medium">Password</label>
                                <Input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
                            </div>
                        </CardContent>
                        <CardFooter>
                            <Button className="w-full" onClick={handleLogin} disabled={loading}>
                                {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                Log In
                            </Button>
                        </CardFooter>
                    </Card>
                </TabsContent>

                <TabsContent value="signup">
                    <Card>
                        <CardHeader>
                            <CardTitle>Create Account</CardTitle>
                            <CardDescription>Start syncing your proxy configuration securely.</CardDescription>
                        </CardHeader>
                        <CardContent className="space-y-4">
                            <div className="space-y-2">
                                <label className="text-sm font-medium">Email</label>
                                <Input placeholder="user@example.com" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
                            </div>
                            <div className="space-y-2">
                                <label className="text-sm font-medium">Password</label>
                                <Input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
                            </div>
                        </CardContent>
                        <CardFooter>
                            <Button className="w-full" onClick={handleSignUp} disabled={loading}>
                                {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                Create Account
                            </Button>
                        </CardFooter>
                    </Card>
                </TabsContent>
            </Tabs>
        </div>
    );
}
