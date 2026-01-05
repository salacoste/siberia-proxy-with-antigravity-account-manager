import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from "@/components/ui/sheet";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Button } from "@/components/ui/button";
import { Share2, Copy, Check, Terminal, ExternalLink, Code } from "lucide-react";
import { useState } from "react";
// @ts-ignore
import { UploadSession, OpenProjectInIDE } from "../../../wailsjs/go/main/App";
// @ts-ignore
import { BrowserOpenURL } from "../../../wailsjs/runtime/runtime";
import { Input } from "@/components/ui/input";
import { DataViewer } from "./DataViewer";
import { Separator } from "@/components/ui/separator";


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
    const [curlCopied, setCurlCopied] = useState(false);

    // If no event, don't render anything (Sheet controls visibility via open prop)
    // But Sheet MUST be rendered to animate out? No, simple conditional is fine for now or keep generic.

    const handleShare = async () => {
        if (!event) return;
        setSharing(true);
        try {
            const url = await UploadSession(event);
            setShareUrl(url);
        } catch (err) {
            console.error("Failed to share:", err);
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

    const generateCurl = () => {
        if (!event) return;
        let cmd = `curl -X ${event.method} "${event.url}"`;
        if (event.req_headers) {
            Object.entries(event.req_headers).forEach(([k, v]) => {
                cmd += ` \\\n  -H "${k}: ${v}"`;
            });
        }
        if (event.req_body) {
            // Escape single quotes
            const escapedBody = event.req_body.replace(/'/g, "'\\''");
            cmd += ` \\\n  -d '${escapedBody}'`;
        }
        navigator.clipboard.writeText(cmd);
        setCurlCopied(true);
        setTimeout(() => setCurlCopied(false), 2000);
    };

    const copyBody = (body: string) => {
        navigator.clipboard.writeText(body || "");
    };

    const openLink = () => {
        if (shareUrl) {
            BrowserOpenURL(shareUrl);
        }
    };

    // Helper to format headers
    const renderHeaders = (headers: Record<string, string>) => {
        if (!headers || Object.keys(headers).length === 0) return <div className="text-xs text-muted-foreground italic">No headers</div>;
        return (
            <div className="grid grid-cols-[120px_1fr] gap-2 text-xs font-mono">
                {Object.entries(headers).map(([k, v]) => (
                    <div key={k} className="contents">
                        <span className="font-semibold text-muted-foreground truncate" title={k}>{k}:</span>
                        <span className="break-all">{v}</span>
                    </div>
                ))}
            </div>
        );
    };

    return (
        <Sheet open={!!event} onOpenChange={(open: boolean) => !open && onClose()}>
            <SheetContent className="w-[800px] sm:w-[540px] md:w-[900px] overflow-hidden flex flex-col p-6">
                <SheetHeader className="pb-4 border-b">
                    <div className="flex justify-between items-start gap-4 pr-6">
                        <div className="space-y-1 overflow-hidden">
                            <SheetTitle className="font-mono flex items-center gap-2 text-lg">
                                <span className={`px-2 py-0.5 rounded text-sm ${event?.method === 'GET' ? 'bg-blue-500/10 text-blue-500' : 'bg-orange-500/10 text-orange-500'}`}>
                                    {event?.method}
                                </span>
                                <span className={event && event.status >= 400 ? "text-red-500" : "text-green-500"}>
                                    {event?.status}
                                </span>
                            </SheetTitle>
                            <SheetDescription className="mx-0 truncate font-mono text-xs text-primary/80 breakdown-all" title={event?.url}>
                                {event?.url}
                            </SheetDescription>
                        </div>
                    </div>
                    {/* Action Bar */}
                    <div className="flex items-center gap-2 pt-2">
                        <Button size="sm" variant="outline" onClick={generateCurl}>
                            {curlCopied ? <Check className="mr-2 h-3 w-3" /> : <Code className="mr-2 h-3 w-3" />}
                            {curlCopied ? "Copied" : "Copy cURL"}
                        </Button>
                        <Separator orientation="vertical" className="h-6" />
                        {shareUrl ? (
                            <div className="flex items-center gap-2 bg-muted rounded p-1 flex-1">
                                <Input value={shareUrl} readOnly className="h-7 text-[10px] font-mono" />
                                <Button size="icon" variant="ghost" className="h-7 w-7" onClick={copyToClipboard}>
                                    {copied ? <Check className="h-3 w-3" /> : <Copy className="h-3 w-3" />}
                                </Button>
                                <Button size="sm" variant="secondary" className="h-7 text-xs" onClick={openLink}>
                                    <ExternalLink className="h-3 w-3" />
                                </Button>
                            </div>
                        ) : (
                            <Button size="sm" variant="outline" onClick={handleShare} disabled={sharing}>
                                <Share2 className="mr-2 h-3 w-3" />
                                {sharing ? "Sharing..." : "Share Session"}
                            </Button>
                        )}
                        <Button size="sm" variant="outline" onClick={() => OpenProjectInIDE()}>
                            <Terminal className="mr-2 h-3 w-3" />
                            IDE
                        </Button>
                    </div>
                </SheetHeader>

                {event && (
                    <Tabs defaultValue="request" className="flex-1 flex flex-col overflow-hidden pt-4">
                        <TabsList className="grid w-full grid-cols-2">
                            <TabsTrigger value="request">Request</TabsTrigger>
                            <TabsTrigger value="response">Response</TabsTrigger>
                        </TabsList>

                        {/* REQUEST TAB */}
                        <TabsContent value="request" className="flex-1 flex flex-col overflow-hidden data-[state=active]:flex gap-4">
                            <div className="flex-1 flex flex-col gap-4 min-h-0">
                                <div className="space-y-2">
                                    <h4 className="text-sm font-medium">Headers</h4>
                                    <ScrollArea className="h-[150px] border rounded-md p-3 bg-muted/20">
                                        {renderHeaders(event.req_headers)}
                                    </ScrollArea>
                                </div>
                                <Separator />
                                <div className="flex-1 flex flex-col space-y-2 min-h-0">
                                    <div className="flex justify-between items-center">
                                        <h4 className="text-sm font-medium flex items-center gap-2">
                                            Body
                                            <span className="text-xs text-muted-foreground font-normal">
                                                {event.req_body?.length || 0} bytes
                                            </span>
                                        </h4>
                                        <Button size="sm" variant="ghost" className="h-6 text-xs" onClick={() => copyBody(event.req_body)}>
                                            <Copy className="mr-1 h-3 w-3" /> Copy Body
                                        </Button>
                                    </div>
                                    <div className="flex-1 min-h-0">
                                        <DataViewer data={event.req_body} type="json" />
                                    </div>
                                </div>
                            </div>
                        </TabsContent>

                        {/* RESPONSE TAB */}
                        <TabsContent value="response" className="flex-1 flex flex-col overflow-hidden data-[state=active]:flex gap-4">
                            <div className="flex-1 flex flex-col gap-4 min-h-0">
                                <div className="space-y-2">
                                    <h4 className="text-sm font-medium">Headers</h4>
                                    <ScrollArea className="h-[150px] border rounded-md p-3 bg-muted/20">
                                        {renderHeaders(event.resp_headers)}
                                    </ScrollArea>
                                </div>
                                <Separator />
                                <div className="flex-1 flex flex-col space-y-2 min-h-0">
                                    <div className="flex justify-between items-center">
                                        <h4 className="text-sm font-medium flex items-center gap-2">
                                            Body
                                            <span className="text-xs text-muted-foreground font-normal">
                                                {event.resp_body?.length || 0} bytes
                                            </span>
                                        </h4>
                                        <Button size="sm" variant="ghost" className="h-6 text-xs" onClick={() => copyBody(event.resp_body)}>
                                            <Copy className="mr-1 h-3 w-3" /> Copy Body
                                        </Button>
                                    </div>
                                    <div className="flex-1 min-h-0">
                                        <DataViewer data={event.resp_body} type="json" />
                                    </div>
                                </div>
                            </div>
                        </TabsContent>
                    </Tabs>
                )}
            </SheetContent>
        </Sheet>
    );
}
