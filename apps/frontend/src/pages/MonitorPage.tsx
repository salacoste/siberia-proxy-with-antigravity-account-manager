import { useEffect, useState, useRef } from 'react';
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Trash2, Pause, Play, Eye, Bug } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { WebSocketViewer } from "../components/monitor/WebSocketViewer";
import { BreakpointPanel } from "../components/monitor/BreakpointPanel";
import { PendingRequestDialog } from "../components/monitor/PendingRequestDialog";
import { RequestDetails } from "../components/monitor/RequestDetails";
import { proxy } from "../../wailsjs/go/models";

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

export default function MonitorPage() {
    const [events, setEvents] = useState<ProxyEvent[]>([]);
    const [paused, setPaused] = useState(false);
    const [selectedEvent, setSelectedEvent] = useState<ProxyEvent | null>(null);

    const [filterQuery, setFilterQuery] = useState("");

    // Breakpoint State
    const [showBreakpoints, setShowBreakpoints] = useState(false);
    const [rules, setRules] = useState<proxy.BreakpointRule[]>([]);
    const [pendingReq, setPendingReq] = useState<proxy.PendingRequest | null>(null);

    // WebSocket State
    const [showWs, setShowWs] = useState(false);
    const [wsFrames, setWsFrames] = useState<proxy.WebSocketFrame[]>([]);


    const pausedRef = useRef(paused);
    useEffect(() => {
        pausedRef.current = paused;
    }, [paused]);

    useEffect(() => {
        // @ts-ignore
        if (window.runtime && window.runtime.EventsOn) {
            // @ts-ignore
            window.runtime.EventsOn("proxy:log", (event: ProxyEvent) => {
                if (pausedRef.current) return;

                setEvents(prev => {
                    const newEvents = [event, ...prev];
                    if (newEvents.length > 500) return newEvents.slice(0, 500); // Cap size
                    return newEvents;
                });
            });

            // @ts-ignore
            window.runtime.EventsOn("breakpoint:hit", (req: proxy.PendingRequest) => {
                setPendingReq(req);
            });

            // @ts-ignore
            window.runtime.EventsOn("proxy:ws:frame", (frame: proxy.WebSocketFrame) => {
                setWsFrames(prev => [...prev, frame]);
            });
        }
    }, []);

    const filterEvent = (evt: ProxyEvent) => {
        if (!filterQuery.trim()) return true;

        try {
            const terms = filterQuery.trim().split(/\s+/);
            return terms.every(term => {
                let checkTerm = term;
                let isNegative = false;
                if (checkTerm.startsWith('!') || (checkTerm.startsWith('-') && checkTerm.length > 1)) {
                    isNegative = true;
                    checkTerm = checkTerm.substring(1);
                }

                let match = false;
                const lowerTerm = checkTerm.toLowerCase();

                // Field filter?
                if (checkTerm.includes(':')) {
                    const [key, val] = checkTerm.split(':');
                    const v = val.toLowerCase();
                    if (key === 'method') match = evt.method.toLowerCase().includes(v);
                    else if (key === 'status') match = evt.status.toString().startsWith(v);
                    else if (key === 'url' || key === 'host') match = evt.url.toLowerCase().includes(v);
                }
                // Regex?
                else if (checkTerm.startsWith('/') && checkTerm.endsWith('/') && checkTerm.length > 2) {
                    try {
                        const regex = new RegExp(checkTerm.slice(1, -1), 'i');
                        match = regex.test(evt.url) || regex.test(evt.method);
                    } catch { match = false; }
                }
                // Standard text
                else {
                    match = evt.url.toLowerCase().includes(lowerTerm) ||
                        evt.method.toLowerCase().includes(lowerTerm) ||
                        evt.status.toString().includes(lowerTerm);
                }

                return isNegative ? !match : match;
            });
        } catch { return true; }
    };

    const getStatusColor = (status: number) => {
        if (status >= 200 && status < 300) return "bg-green-500/10 text-green-500 border-green-500/50";
        if (status >= 300 && status < 400) return "bg-blue-500/10 text-blue-500 border-blue-500/50";
        if (status >= 400 && status < 500) return "bg-yellow-500/10 text-yellow-500 border-yellow-500/50";
        if (status >= 500) return "bg-red-500/10 text-red-500 border-red-500/50";
        return "bg-gray-500/10 text-gray-500 border-gray-500/50";
    }

    const getMethodColor = (method: string) => {
        switch (method) {
            case 'GET': return "default";
            case 'POST': return "secondary";
            case 'PUT': return "outline";
            case 'DELETE': return "destructive";
            default: return "outline";
        }
    }

    const filteredEvents = events.filter(filterEvent);

    return (
        <div className="p-8 h-full flex flex-col space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold">Traffic Monitor</h1>
                    <p className="text-sm text-muted-foreground">Real-time HTTP request inspection</p>
                </div>
                <div className="flex gap-2">
                    <Button variant={showWs ? "secondary" : "outline"} size="sm" onClick={() => setShowWs(!showWs)}>
                        <Bug className="mr-2 h-4 w-4" /> WebSockets ({wsFrames.length})
                    </Button>
                    <Button variant={showBreakpoints ? "secondary" : "outline"} size="sm" onClick={() => setShowBreakpoints(!showBreakpoints)}>
                        <Bug className="mr-2 h-4 w-4" /> Breakpoints
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => setPaused(!paused)}>
                        {paused ? <Play className="mr-2 h-4 w-4" /> : <Pause className="mr-2 h-4 w-4" />}
                        {paused ? "Resume" : "Pause"}
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => { setEvents([]); setWsFrames([]); }}>
                        <Trash2 className="mr-2 h-4 w-4" /> Clear
                    </Button>
                </div>
            </div>

            {showWs && (
                <div className="h-64 shrink-0">
                    <WebSocketViewer frames={wsFrames} onClear={() => setWsFrames([])} />
                </div>
            )}

            {showBreakpoints && (
                <BreakpointPanel
                    rules={rules}
                    onAdd={(pattern) => setRules([...rules, new proxy.BreakpointRule({ id: Date.now().toString(), pattern, enabled: true, method: "*" })])}
                    onDelete={(id) => setRules(rules.filter(r => r.id !== id))}
                />
            )}

            <div className="flex gap-2">
                <Input
                    placeholder="Filter (e.g. method:POST /api/ status:404 !css)"
                    value={filterQuery}
                    onChange={(e) => setFilterQuery(e.target.value)}
                    className="font-mono text-xs"
                />
            </div>

            <div className="border rounded-md flex-1 overflow-auto bg-background/50 backdrop-blur relative">
                <Table>
                    <TableHeader className="sticky top-0 bg-background z-10 shadow-sm">
                        <TableRow>
                            <TableHead className="w-[100px]">Time</TableHead>
                            <TableHead className="w-[80px]">Method</TableHead>
                            <TableHead className="w-[100px]">Status</TableHead>
                            <TableHead>URL</TableHead>
                            <TableHead className="w-[100px] text-right">Size</TableHead>
                            <TableHead className="w-[100px] text-right">Duration</TableHead>
                            <TableHead className="w-[50px]"></TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {filteredEvents.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={7} className="text-center h-24 text-muted-foreground">
                                    {events.length === 0 ? "No requests captured yet." : "No matching requests found."}
                                </TableCell>
                            </TableRow>
                        ) : (
                            filteredEvents.map((evt, i) => (
                                <TableRow key={i} className="font-mono text-xs hover:bg-muted/50 transition-colors cursor-pointer" onClick={() => setSelectedEvent(evt)}>
                                    <TableCell className="text-muted-foreground whitespace-nowrap">{evt.time}</TableCell>
                                    <TableCell>
                                        <Badge variant={getMethodColor(evt.method) as any} className="text-[10px] px-2 py-0.5 h-5 font-bold">
                                            {evt.method}
                                        </Badge>
                                    </TableCell>
                                    <TableCell>
                                        <Badge variant="outline" className={`border ${getStatusColor(evt.status)}`}>
                                            {evt.status > 0 ? evt.status : '---'}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="max-w-[400px] truncate" title={evt.url}>
                                        {evt.url}
                                    </TableCell>
                                    <TableCell className="text-right text-muted-foreground">
                                        {evt.size} B
                                    </TableCell>
                                    <TableCell className="text-right font-medium">
                                        {evt.duration_ms}ms
                                    </TableCell>
                                    <TableCell>
                                        <Button variant="ghost" size="icon" className="h-6 w-6">
                                            <Eye className="h-4 w-4" />
                                        </Button>
                                    </TableCell>
                                </TableRow>
                            ))
                        )}
                    </TableBody>
                </Table>
            </div>

            {/* Pending Request Interception Dialog */}
            <PendingRequestDialog
                pendingReq={pendingReq}
                onClose={() => setPendingReq(null)}
            />

            {/* Request Details Dialog */}
            <RequestDetails
                event={selectedEvent}
                onClose={() => setSelectedEvent(null)}
            />
        </div>
    );
}
