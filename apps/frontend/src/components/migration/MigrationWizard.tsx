import React, { useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Check, Loader2, Play } from "lucide-react";
// @ts-ignore
import { ScanIDEAccounts, ImportDiscoveredAccount } from "../../../wailsjs/go/main/App";

// Mock type if bindings aren't generated yet
interface DiscoveredAccount {
    name: string;
    path: string;
    raw_token: string;
    masked_token: string;
    source: string;
}

interface MigrationWizardProps {
    open: boolean;
    onOpenChange: (open: boolean) => void;
    onImportComplete?: () => void;
}

export function MigrationWizard({ open, onOpenChange, onImportComplete }: MigrationWizardProps) {
    const [step, setStep] = useState<'scan' | 'list' | 'importing' | 'success'>('scan');
    const [accounts, setAccounts] = useState<DiscoveredAccount[]>([]);
    const [selectedAccounts, setSelectedAccounts] = useState<string[]>([]); // raw_tokens


    // Auto-scan on open
    React.useEffect(() => {
        if (open && step === 'scan') {
            performScan();
        }
    }, [open]);

    const performScan = async () => {
        try {
            const results = await ScanIDEAccounts();
            if (results && results.length > 0) {
                setAccounts(results);
                // Select all by default
                setSelectedAccounts(results.map((r: DiscoveredAccount) => r.raw_token));
                setStep('list');
            } else {
                // No accounts found
                setStep('list'); // List empty
            }
        } catch (e) {
            console.error("Scan failed", e);
            setStep('list');
        }
    };

    const handleImport = async () => {
        setStep('importing');
        for (const token of selectedAccounts) {
            try {
                await ImportDiscoveredAccount(token);
            } catch (e) {
                console.error("Import failed for token", e);
            }
        }
        setStep('success');
        if (onImportComplete) onImportComplete();
        // Close after delay
        setTimeout(() => {
            onOpenChange(false);
        }, 2000);
    };

    return (
        <Dialog open={open} onOpenChange={onOpenChange}>
            <DialogContent className="sm:max-w-[500px]">
                <DialogHeader>
                    <DialogTitle>Welcome to Siberia</DialogTitle>
                    <DialogDescription>
                        Streamline your setup by importing existing sessions from your IDEs.
                    </DialogDescription>
                </DialogHeader>

                <div className="py-4">
                    {step === 'scan' && (
                        <div className="flex flex-col items-center justify-center space-y-4 py-8">
                            <Loader2 className="h-8 w-8 animate-spin text-primary" />
                            <p className="text-sm text-muted-foreground">Scanning for existing accounts...</p>
                        </div>
                    )}

                    {step === 'list' && (
                        <div className="space-y-4">
                            {accounts.length === 0 ? (
                                <div className="text-center py-8 text-muted-foreground">
                                    No accounts found in standard locations.
                                </div>
                            ) : (
                                <div className="space-y-2 max-h-[200px] overflow-y-auto border rounded-md p-2">
                                    {accounts.map((acc, idx) => (
                                        <div key={idx} className="flex items-center space-x-3 p-2 hover:bg-muted/50 rounded cursor-pointer"
                                            onClick={() => {
                                                if (selectedAccounts.includes(acc.raw_token)) {
                                                    setSelectedAccounts(selectedAccounts.filter(t => t !== acc.raw_token));
                                                } else {
                                                    setSelectedAccounts([...selectedAccounts, acc.raw_token]);
                                                }
                                            }}
                                        >
                                            <div className={`w-4 h-4 border rounded flex items-center justify-center ${selectedAccounts.includes(acc.raw_token) ? "bg-primary border-primary text-primary-foreground" : "border-muted-foreground"}`}>
                                                {selectedAccounts.includes(acc.raw_token) && <Check className="h-3 w-3" />}
                                            </div>
                                            <div className="flex-1 overflow-hidden">
                                                <p className="text-sm font-medium truncate">{acc.name}</p>
                                                <p className="text-xs text-muted-foreground truncate">{acc.path}</p>
                                            </div>
                                            <div className="text-xs bg-muted px-2 py-1 rounded">
                                                {acc.source}
                                            </div>
                                        </div>
                                    ))}
                                </div>
                            )}
                            <div className="text-xs text-muted-foreground">
                                {selectedAccounts.length} account(s) selected for import.
                            </div>
                        </div>
                    )}

                    {step === 'importing' && (
                        <div className="flex flex-col items-center justify-center space-y-4 py-8">
                            <Loader2 className="h-8 w-8 animate-spin text-primary" />
                            <p className="text-sm text-muted-foreground">Importing secure tokens...</p>
                        </div>
                    )}

                    {step === 'success' && (
                        <div className="flex flex-col items-center justify-center space-y-4 py-8 text-green-500">
                            <div className="h-12 w-12 rounded-full bg-green-100 flex items-center justify-center dark:bg-green-900/30">
                                <Check className="h-6 w-6" />
                            </div>
                            <p className="text-sm font-medium text-foreground">Import Successful!</p>
                        </div>
                    )}
                </div>

                <DialogFooter>
                    {step === 'list' && (
                        <>
                            <Button variant="ghost" onClick={() => onOpenChange(false)}>Skip</Button>
                            <Button onClick={handleImport} disabled={selectedAccounts.length === 0}>
                                <Play className="mr-2 h-4 w-4" /> Import Selected
                            </Button>
                        </>
                    )}
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
