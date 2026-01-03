import { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Trash2, Plus } from "lucide-react";
import { AddBreakpointRule, DeleteBreakpointRule } from "../../../wailsjs/go/main/App";
import { proxy } from "../../../wailsjs/go/models";

// We'll manage local state for the UI, but we should probably fetch from backend on mount?
// For MVP, valid rules are added via one-way operation.

interface BreakpointPanelProps {
    rules: proxy.BreakpointRule[];
    onAdd: (pattern: string) => void;
    onDelete: (id: string) => void;
}

export function BreakpointPanel({ rules, onAdd, onDelete }: BreakpointPanelProps) {
    const [pattern, setPattern] = useState("");

    const handleAdd = async () => {
        if (!pattern) return;
        await AddBreakpointRule(pattern, "*");
        onAdd(pattern); // Notify parent to refresh or just optimistic update
        setPattern("");
    };

    const handleDelete = async (id: string) => {
        await DeleteBreakpointRule(id);
        onDelete(id);
    };

    return (
        <div className="flex flex-col gap-4 p-4 border rounded-md bg-muted/20">
            <h3 className="font-semibold text-sm">Breakpoint Rules</h3>
            <div className="flex gap-2">
                <Input
                    placeholder="URL contains (e.g. /api/auth)"
                    value={pattern}
                    onChange={(e) => setPattern(e.target.value)}
                    className="h-8 text-xs font-mono"
                />
                <Button size="sm" className="h-8" onClick={handleAdd}>
                    <Plus className="w-4 h-4" />
                </Button>
            </div>
            <div className="flex flex-col gap-2">
                {rules.map(rule => (
                    <div key={rule.id} className="flex items-center justify-between text-xs bg-background p-2 rounded border">
                        <span className="font-mono">{rule.pattern}</span>
                        <Button variant="ghost" size="icon" className="h-6 w-6" onClick={() => handleDelete(rule.id)}>
                            <Trash2 className="w-3 h-3 text-red-500" />
                        </Button>
                    </div>
                ))}
                {rules.length === 0 && <span className="text-muted-foreground text-[10px]">No active breakpoints.</span>}
            </div>
        </div>
    );
}
