import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';

import { Save, Play, AlertTriangle } from 'lucide-react';
import { toast } from 'sonner';
import { UpdateScript, GetScript } from '../../wailsjs/go/main/App';

const DEFAULT_SCRIPT = `// Traffic Scripting (Goja JS)
// Available hooks:
// 1. onRequest(req) -> return modified req
// 2. onResponse(res) -> return modified res

function onRequest(req) {
    // Example: Add a custom header
    // req.Headers["X-Siberia-Script"] = ["Enabled"];
    
    // Example: Log URL
    // console.log("Request: " + req.URL);
    
    return req;
}

function onResponse(res) {
    // Example: Modify JSON body
    // if (res.StatusCode === 200) {
    //    var body = JSON.parse(res.Body);
    //    body.scripted = true;
    //    res.Body = JSON.stringify(body);
    // }
    
    return res;
}
`;

export default function ScriptingPage() {
    const [script, setScript] = useState(DEFAULT_SCRIPT);
    const [isActive, setIsActive] = useState(false);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        loadScript();
    }, []);

    const loadScript = async () => {
        try {
            const state = await GetScript();
            if (state) {
                setScript(state.code || "");
                setIsActive(state.active || false);
            }
        } catch (error) {
            console.error("Failed to load script", error);
        }
    };


    const handleSave = async () => {
        setLoading(true);
        try {
            await UpdateScript(script, isActive);
            toast.success("Script updated successfully");
        } catch (error) {
            toast.error("Failed to update script: " + error);
        } finally {
            setLoading(false);
        }
    };

    const toggleActive = async (val: boolean) => {
        setIsActive(val);
        // Auto-save on toggle? Or just local state?
        // Let's auto-save to ensure backend is in sync immediately
        try {
            await UpdateScript(script, val);
            toast.success(val ? "Scripting Enabled" : "Scripting Disabled");
        } catch (error) {
            toast.error("Failed to toggle: " + error);
            setIsActive(!val); // Revert
        }
    };

    return (
        <div className="container mx-auto p-6 space-y-6 h-full flex flex-col">
            <div className="flex justify-between items-center shrink-0">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Traffic Scripting</h1>
                    <p className="text-muted-foreground flex items-center gap-2">
                        Modify HTTP traffic on the fly using JavaScript (ES5/Goja).
                    </p>
                </div>
                <div className="flex items-center gap-4">
                    <div className="flex items-center space-x-2">
                        <Switch id="active-mode" checked={isActive} onCheckedChange={toggleActive} />
                        <Label htmlFor="active-mode" className={isActive ? "font-bold text-green-600" : "text-muted-foreground"}>
                            {isActive ? "Enabled" : "Disabled"}
                        </Label>
                    </div>
                    <Button onClick={handleSave} disabled={loading}>
                        <Save className="mr-2 h-4 w-4" />
                        Save & Apply
                    </Button>
                </div>
            </div>

            <Card className="flex-1 flex flex-col overflow-hidden border-zinc-200 dark:border-zinc-800 shadow-sm">
                <CardHeader className="py-3 px-4 bg-muted/30 border-b flex flex-row justify-between items-center shrink-0">
                    <div className="space-y-1">
                        <CardTitle className="text-base font-mono flex items-center gap-2">
                            <Play className="h-4 w-4 text-blue-500" />
                            script.js
                        </CardTitle>
                    </div>
                    <div className="text-xs text-muted-foreground flex items-center gap-1">
                        <AlertTriangle className="h-3 w-3 text-amber-500" />
                        <span>Changes take effect immediately on Save</span>
                    </div>
                </CardHeader>
                <div className="flex-1 relative bg-zinc-950">
                    <textarea
                        className="w-full h-full resize-none bg-transparent text-zinc-100 font-mono text-sm p-4 outline-none leading-relaxed"
                        value={script}
                        onChange={(e) => setScript(e.target.value)}
                        spellCheck={false}
                    />
                </div>
            </Card>
        </div>
    );
}
