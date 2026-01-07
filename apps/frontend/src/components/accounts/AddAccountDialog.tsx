import { useState } from 'react';
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Plus, Loader2 } from "lucide-react";

interface AddAccountDialogProps {
    onAccountAdded: () => void;
}

export function AddAccountDialog({ onAccountAdded }: AddAccountDialogProps) {
    const [open, setOpen] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState('');

    const [formData, setFormData] = useState({
        email: '',
        password: '',
        recovery: '',
        proxyGroup: 'default'
    });

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError('');

        try {
            // @ts-ignore
            if (window.go?.main?.App?.CreateAccount) {
                // @ts-ignore
                await window.go.main.App.CreateAccount(
                    formData.email,
                    formData.password,
                    formData.recovery,
                    formData.proxyGroup
                );

                setOpen(false);
                setFormData({ email: '', password: '', recovery: '', proxyGroup: 'default' });
                onAccountAdded();
            } else {
                console.warn("Backend not available");
                setError("Backend not connected");
            }
        } catch (err: any) {
            console.error("Failed to create account:", err);
            setError(err.toString());
        } finally {
            setIsLoading(false);
        }
    };

    const handleOAuthLogin = async () => {
        setIsLoading(true);
        setError('');
        try {
            // @ts-ignore
            if (window.go?.main?.App?.LoginWithOAuth) {
                // @ts-ignore
                await window.go.main.App.LoginWithOAuth("google");
                setOpen(false);
                onAccountAdded();
            } else {
                console.warn("Backend LoginWithOAuth not available");
                setError("Backend API mismatch");
            }
        } catch (err: any) {
            console.error("OAuth failed:", err);
            setError(err.toString());
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Dialog open={open} onOpenChange={setOpen}>
            <DialogTrigger asChild>
                <Button>
                    <Plus className="mr-2 h-4 w-4" /> Add Account
                </Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">
                <div className="grid gap-4 py-4">
                    <Button variant="outline" onClick={handleOAuthLogin} disabled={isLoading} className="w-full">
                        <img src="https://www.google.com/favicon.ico" className="mr-2 h-4 w-4" alt="Google" />
                        Add with Google
                    </Button>

                    <div className="relative">
                        <div className="absolute inset-0 flex items-center">
                            <span className="w-full border-t" />
                        </div>
                        <div className="relative flex justify-center text-xs uppercase">
                            <span className="bg-background px-2 text-muted-foreground">
                                Or manually
                            </span>
                        </div>
                    </div>

                    <form onSubmit={handleSubmit}>
                        <div className="grid gap-4 py-2">
                            {error && (
                                <div className="text-sm text-red-500 font-medium">
                                    {error}
                                </div>
                            )}
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="email" className="text-right">
                                    Email
                                </Label>
                                <Input
                                    id="email"
                                    type="email"
                                    value={formData.email}
                                    onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                                    className="col-span-3"
                                    required
                                />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="password" className="text-right">
                                    Password
                                </Label>
                                <Input
                                    id="password"
                                    type="password"
                                    value={formData.password}
                                    onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                                    className="col-span-3"
                                    required
                                />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="recovery" className="text-right">
                                    Recovery
                                </Label>
                                <Input
                                    id="recovery"
                                    type="email"
                                    placeholder="Optional"
                                    value={formData.recovery}
                                    onChange={(e) => setFormData({ ...formData, recovery: e.target.value })}
                                    className="col-span-3"
                                />
                            </div>
                            <div className="grid grid-cols-4 items-center gap-4">
                                <Label htmlFor="group" className="text-right">
                                    Group
                                </Label>
                                <Input
                                    id="group"
                                    value={formData.proxyGroup}
                                    onChange={(e) => setFormData({ ...formData, proxyGroup: e.target.value })}
                                    className="col-span-3"
                                />
                            </div>
                        </div>
                        <DialogFooter>
                            <Button type="submit" disabled={isLoading}>
                                {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                Add Account
                            </Button>
                        </DialogFooter>
                    </form>
            </DialogContent>
        </Dialog>
    );
}
