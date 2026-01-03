import { useEffect, useState, useRef } from 'react';
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Trash2, Pause, Play, Eye } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ScrollArea } from "@/components/ui/scroll-area";

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
        }
    }, []);

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

    return (
        <div className="p-8 h-full flex flex-col space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold">Traffic Monitor</h1>
                    <p className="text-sm text-muted-foreground">Real-time HTTP request inspection</p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" size="sm" onClick={() => setPaused(!paused)}>
                        {paused ? <Play className="mr-2 h-4 w-4" /> : <Pause className="mr-2 h-4 w-4" />}
                        {paused ? "Resume" : "Pause"}
                    </Button>
                    <Button variant="outline" size="sm" onClick={() => setEvents([])}>
                        <Trash2 className="mr-2 h-4 w-4" /> Clear
                    </Button>
                </div>
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
                        {events.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={7} className="text-center h-24 text-muted-foreground">
                                    No requests captured yet.
                                </TableCell>
                            </TableRow>
                        ) : (
                            events.map((evt, i) => (
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

            <Dialog open={selectedEvent !== null} onOpenChange={(open) => !open && setSelectedEvent(null)}>
                <DialogContent className="max-w-4xl max-h-[80vh] flex flex-col">
                    <DialogHeader>
                        <DialogTitle className="font-mono flex items-center gap-2">
                            {selectedEvent?.method} {selectedEvent?.status}
                        </DialogTitle>
                        <DialogDescription className="truncate font-mono text-xs">
                            {selectedEvent?.url}
                        </DialogDescription>
                    </DialogHeader>

                    {selectedEvent && (
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
                                                {selectedEvent.req_headers && Object.entries(selectedEvent.req_headers).map(([k, v]) => `${k}: ${v}\n`).join('')}
                                            </pre>
                                        </ScrollArea>
                                    </div>
                                    <div className="border rounded bg-muted/50 p-2 flex flex-col">
                                        <h4 className="text-xs font-semibold mb-2">Body</h4>
                                        <ScrollArea className="flex-1">
                                            <pre className="text-[10px] whitespace-pre-wrap">{selectedEvent.req_body || "<Empty>"}</pre>
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
                                                {selectedEvent.resp_headers && Object.entries(selectedEvent.resp_headers).map(([k, v]) => `${k}: ${v}\n`).join('')}
                                            </pre>
                                        </ScrollArea>
                                    </div>
                                    <div className="border rounded bg-muted/50 p-2 flex flex-col">
                                        <h4 className="text-xs font-semibold mb-2">Body</h4>
                                        <ScrollArea className="flex-1">
                                            <pre className="text-[10px] whitespace-pre-wrap">{selectedEvent.resp_body || "<Empty>"}</pre>
                                        </ScrollArea>
                                    </div>
                                </div>
                            </TabsContent>
                        </Tabs>
                    )}
                </DialogContent>
            </Dialog>
        </div>
    );
}
