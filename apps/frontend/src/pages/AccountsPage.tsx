import { useEffect, useState } from 'react';
import { Play, Loader2 } from "lucide-react";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { MoreHorizontal } from "lucide-react";
import { AddAccountDialog } from "@/components/accounts/AddAccountDialog";

// Define the interface manually since bindings aren't auto-generated in this env
// In a real dev flow, we'd import { ListAccounts, DeleteAccount } from '../../wailsjs/go/main/App';
interface Account {
    id: number;
    email: string;
    proxy_group: string;
    is_active: boolean;
    created_at: string;
}

export default function AccountsPage() {
    const [accounts, setAccounts] = useState<Account[]>([]);
    const [activatingId, setActivatingId] = useState<number | null>(null);

    const handleActivate = async (id: number) => {
        if (activatingId) return;
        setActivatingId(id);
        try {
            // @ts-ignore
            if (window.go?.main?.App?.ActivateAccount) {
                // @ts-ignore
                await window.go.main.App.ActivateAccount(id);
                console.log("Activation successful");
            } else {
                console.warn("Backend not available");
            }
        } catch (error) {
            console.error("Failed to activate account:", error);
        } finally {
            setActivatingId(null);
        }
    };

    // Mock/Real data fetch
    const fetchAccounts = async () => {
        try {
            // @ts-ignore: Wails bindings might not be generated yet in this env
            if (window.go && window.go.main && window.go.main.App && window.go.main.App.ListAccounts) {
                // @ts-ignore
                const data = await window.go.main.App.ListAccounts();
                setAccounts(data || []);
            } else {
                console.warn("Wails backend not available, usage mock data?");
            }
        } catch (error) {
            console.error("Failed to fetch accounts:", error);
        }
    };

    const handleDelete = async (id: number) => {
        try {
            // @ts-ignore
            if (window.go?.main?.App?.DeleteAccount) {
                // @ts-ignore
                await window.go.main.App.DeleteAccount(id);
                fetchAccounts();
            }
        } catch (error) {
            console.error("Failed to delete account:", error);
        }
    };

    useEffect(() => {
        fetchAccounts();
    }, []);

    return (
        <div className="p-8 space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold">Accounts</h1>
                    <p className="text-sm text-muted-foreground">Manage your Google accounts</p>
                </div>
                <AddAccountDialog onAccountAdded={fetchAccounts} />
            </div>

            <div className="border rounded-md">
                <Table>
                    <TableHeader>
                        <TableRow>
                            <TableHead className="w-[50px]"></TableHead>
                            <TableHead>Email</TableHead>
                            <TableHead>Group</TableHead>
                            <TableHead>Status</TableHead>
                            <TableHead className="text-right">Actions</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {accounts.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={5} className="text-center h-24 text-muted-foreground">
                                    No accounts found. Add one to get started.
                                </TableCell>
                            </TableRow>
                        ) : (
                            accounts.map((acc) => (
                                <TableRow key={acc.id}>
                                    <TableCell>
                                        <Avatar className="h-8 w-8">
                                            <AvatarFallback>{acc.email.substring(0, 2).toUpperCase()}</AvatarFallback>
                                        </Avatar>
                                    </TableCell>
                                    <TableCell className="font-medium">{acc.email}</TableCell>
                                    <TableCell>{acc.proxy_group}</TableCell>
                                    <TableCell>
                                        <Badge variant={acc.is_active ? "default" : "secondary"}>
                                            {acc.is_active ? "Active" : "Inactive"}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="text-right">
                                        <Button
                                            variant="ghost"
                                            size="sm"
                                            onClick={() => handleActivate(acc.id)}
                                            disabled={!!activatingId}
                                            className="mr-2"
                                        >
                                            {activatingId === acc.id ? (
                                                <Loader2 className="h-4 w-4 animate-spin text-primary" />
                                            ) : (
                                                <Play className="h-4 w-4 text-green-600 hover:text-green-700" />
                                            )}
                                            <span className="sr-only">Activate</span>
                                        </Button>
                                        <DropdownMenu>
                                            <DropdownMenuTrigger asChild>
                                                <Button variant="ghost" className="h-8 w-8 p-0">
                                                    <span className="sr-only">Open menu</span>
                                                    <MoreHorizontal className="h-4 w-4" />
                                                </Button>
                                            </DropdownMenuTrigger>
                                            <DropdownMenuContent align="end">
                                                <DropdownMenuLabel>Actions</DropdownMenuLabel>
                                                <DropdownMenuItem onClick={() => navigator.clipboard.writeText(acc.email)}>
                                                    Copy Email
                                                </DropdownMenuItem>
                                                <DropdownMenuItem className="text-destructive" onClick={() => handleDelete(acc.id)}>
                                                    Delete
                                                </DropdownMenuItem>
                                            </DropdownMenuContent>
                                        </DropdownMenu>
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>
        </div>
    );
}
