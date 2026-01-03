import { useEffect, useState } from 'react';
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Input } from "@/components/ui/input";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ResumeRequest } from "../../../wailsjs/go/main/App";
import { proxy } from "../../../wailsjs/go/models";

// Need Textarea component, assume standard HTML for now or check UI lib
// Assuming we don't have a configured Textarea yet, I'll use standard <textarea> with tailwind classes

interface PendingRequestDialogProps {
    pendingReq: proxy.PendingRequest | null;
    onClose: () => void;
}

export function PendingRequestDialog({ pendingReq, onClose }: PendingRequestDialogProps) {
    const [body, setBody] = useState("");
    const [headers, setHeaders] = useState<Record<string, string>>({});
    const [method, setMethod] = useState("");
    const [url, setUrl] = useState("");

    useEffect(() => {
        if (pendingReq) {
            setBody(pendingReq.body || "");
            setHeaders({ ...pendingReq.headers });
            setMethod(pendingReq.method);
            setUrl(pendingReq.url);
        }
    }, [pendingReq]);

    const handleResume = async (drop: boolean) => {
        if (!pendingReq) return;

        const mod = new proxy.ModifiedRequest({
            drop: drop,
            method: method,
            url: url,
            headers: headers,
            body: body
        });

        await ResumeRequest(pendingReq.id, mod);
        onClose();
    };

    if (!pendingReq) return null;

    return (
        <Dialog open={!!pendingReq} onOpenChange={() => { }}>
            <DialogContent className="max-w-4xl max-h-[90vh] flex flex-col pointer-events-auto">
                <DialogHeader>
                    <DialogTitle className="text-red-500 flex items-center gap-2">
                        🛑 Paused Request
                    </DialogTitle>
                </DialogHeader>

                <div className="flex gap-2 mb-2">
                    <Input
                        value={method}
                        onChange={e => setMethod(e.target.value)}
                        className="w-24 font-bold"
                    />
                    <Input
                        value={url}
                        onChange={e => setUrl(e.target.value)}
                        className="flex-1 font-mono text-sm"
                    />
                </div>

                <Tabs defaultValue="body" className="flex-1 flex flex-col overflow-hidden">
                    <TabsList>
                        <TabsTrigger value="body">Body</TabsTrigger>
                        <TabsTrigger value="headers">Headers</TabsTrigger>
                    </TabsList>

                    <TabsContent value="body" className="flex-1 mt-2">
                        <textarea
                            className="w-full h-full min-h-[300px] p-2 font-mono text-xs bg-muted/50 rounded border resize-none focus:outline-none focus:ring-1 focus:ring-ring"
                            value={body}
                            onChange={(e) => setBody(e.target.value)}
                        />
                    </TabsContent>

                    <TabsContent value="headers" className="flex-1 mt-2">
                        <ScrollArea className="h-[300px] border rounded p-2">
                            {Object.entries(headers).map(([k, v]) => (
                                <div key={k} className="flex gap-2 mb-2">
                                    <Input
                                        value={k}
                                        disabled
                                        className="w-1/3 text-xs bg-muted"
                                    />
                                    <Input
                                        value={v}
                                        onChange={(e) => setHeaders({ ...headers, [k]: e.target.value })}
                                        className="flex-1 text-xs"
                                    />
                                </div>
                            ))}
                        </ScrollArea>
                    </TabsContent>
                </Tabs>

                <DialogFooter className="gap-2 sm:justify-between">
                    <Button variant="destructive" onClick={() => handleResume(true)}>
                        Drop Request
                    </Button>
                    <div className="flex gap-2">
                        <Button variant="outline" onClick={() => handleResume(false)}>
                            Resume (No Edits)
                        </Button>
                        <Button onClick={() => handleResume(false)}>
                            Resume with Changes
                        </Button>
                    </div>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
}
