import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { Share2, Copy, Check } from "lucide-react";
import { useState } from "react";
// @ts-ignore
import { UploadSession } from "../../../wailsjs/go/main/App";
// @ts-ignore
import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime";

interface ProxyEvent {
    method: string;
    url: string;
    status: number;
    duration_ms: number;
    time: string;
    size: number;
    req_headers: Record<string, string>;
    resp_headers: Record<string, string>;
    req_body: string;
    resp_body: string;
}

interface RequestDetailsProps {
    event: ProxyEvent | null;
    onClose: () => void;
}

export function RequestDetails({ event, onClose }: RequestDetailsProps) {
    const [sharing, setSharing] = useState(false);
    const [shareUrl, setShareUrl] = useState<string | null>(null);
    const [copied, setCopied] = useState(false);

    if (!event) return null;

    const handleShare = async () => {
        setSharing(true);
        try {
            const url = await UploadSession(event);
            setShareUrl(url);
        } catch (err) {
            console.error("Failed to share:", err);
            alert("Failed to create share link");
        } finally {
            setSharing(false);
        }
    };

    const copyToClipboard = () => {
        if (shareUrl) {
            navigator.clipboard.writeText(shareUrl);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        }
    };

    const openLink = () => {
        if (shareUrl) {
            BrowserOpenURL(shareUrl);
        }
    };

    return (
        <Dialog open={!!event} onOpenChange={(open) => !open && onClose()}>
            <DialogContent className="max-w-4xl max-h-[80vh] flex flex-col">
                <DialogHeader>
                    <div className="flex justify-between items-start mr-8">
                        <div>
                            <DialogTitle className="font-mono flex items-center gap-2">
                                {event.method} {event.status}
                            </DialogTitle>
                            <DialogDescription className="truncate font-mono text-xs mt-1">
                                {event.url}
                            </DialogDescription>
                        </div>
                        <div className="flex gap-2">
                            {shareUrl ? (
                                <div className="flex items-center gap-2 bg-muted rounded p-1">
                                    <Input value={shareUrl} readOnly className="h-7 w-48 text-[10px]" />
                                    <Button size="icon" variant="ghost" className="h-7 w-7" onClick={copyToClipboard}>
                                        {copied ? <Check className="h-3 w-3" /> : <Copy className="h-3 w-3" />}
                                    </Button>
                                    <Button size="sm" variant="secondary" className="h-7 text-xs" onClick={openLink}>
                                        Open
                                    </Button>
                                    <Button size="sm" variant="ghost" className="h-7 text-xs text-muted-foreground" onClick={() => setShareUrl(null)}>
                                        Reset
                                    </Button>
                                </div>
                            ) : (
                                <Button size="sm" variant="outline" onClick={handleShare} disabled={sharing}>
                                    <Share2 className="mr-2 h-3 w-3" />
                                    {sharing ? "Creating..." : "Share Session"}
                                </Button>
                            )}
                        </div>
                    </div>
                </DialogHeader>

                <Tabs defaultValue="request" className="flex-1 flex flex-col overflow-hidden">
                    <TabsList>
                        <TabsTrigger value="request">Request</TabsTrigger>
                        <TabsTrigger value="response">Response</TabsTrigger>
                    </TabsList>
                    <TabsContent value="request" className="flex-1 flex flex-col gap-4 overflow-hidden mt-4">
                        <div className="grid grid-cols-2 gap-4 h-full">
                            <div className="border rounded bg-muted/50 p-2 flex flex-col">
                                <h4 className="text-xs font-semibold mb-2">Headers</h4>
                                <ScrollArea className="flex-1">
                                    <pre className="text-[10px] whitespace-pre-wrap">
                                        {event.req_headers && Object.entries(event.req_headers).map(([k, v]) => `${k}: ${v}\n`).join('')}
                                    </pre>
                                </ScrollArea>
                            </div>
                            <div className="border rounded bg-muted/50 p-2 flex flex-col">
                                <h4 className="text-xs font-semibold mb-2">Body</h4>
                                <ScrollArea className="flex-1">
                                    <pre className="text-[10px] whitespace-pre-wrap">{event.req_body || "<Empty>"}</pre>
                                </ScrollArea>
                            </div>
                        </div>
                    </TabsContent>
                    <TabsContent value="response" className="flex-1 flex flex-col gap-4 overflow-hidden mt-4">
                        <div className="grid grid-cols-2 gap-4 h-full">
                            <div className="border rounded bg-muted/50 p-2 flex flex-col">
                                <h4 className="text-xs font-semibold mb-2">Headers</h4>
                                <ScrollArea className="flex-1">
                                    <pre className="text-[10px] whitespace-pre-wrap">
                                        {event.resp_headers && Object.entries(event.resp_headers).map(([k, v]) => `${k}: ${v}\n`).join('')}
                                    </pre>
                                </ScrollArea>
                            </div>
                            <div className="border rounded bg-muted/50 p-2 flex flex-col">
                                <h4 className="text-xs font-semibold mb-2">Body</h4>
                                <ScrollArea className="flex-1">
                                    <pre className="text-[10px] whitespace-pre-wrap">{event.resp_body || "<Empty>"}</pre>
                                </ScrollArea>
                            </div>
                        </div>
                    </TabsContent>
                </Tabs>
            </DialogContent>
        </Dialog>
    );
}

// Simple Input shim if using ShadCN one within this file without import, but better to import it.
// I'll import Input at top
import { Input } from "@/components/ui/input";
