import { useEffect, useState } from 'react';
import { Badge } from "@/components/ui/badge";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Trash2, Pause, Play } from "lucide-react";
import { Button } from "@/components/ui/button";

interface ProxyEvent {
    method: string;
    url: string;
    status: number;
    duration_ms: number;
    time: string;
}

export default function MonitorPage() {
    const [events, setEvents] = useState<ProxyEvent[]>([]);
    const [paused, setPaused] = useState(false);
    // Auto-scroll ref not actually used with Table approach easily without scroll area
    // Removing it to satisfy linter

    useEffect(() => {
        // @ts-ignore
        if (window.runtime && window.runtime.EventsOn) {
            // @ts-ignore
            window.runtime.EventsOn("proxy:request", (event: ProxyEvent) => {
                if (paused) return;
                setEvents(prev => {
                    const newEvents = [event, ...prev];
                    if (newEvents.length > 100) return newEvents.slice(0, 100);
                    return newEvents;
                });
            });
        }

        return () => {
            // Cleanup if needed, though Wails events usually persist
            // @ts-ignore
            if (window.runtime && window.runtime.EventsOff) {
                // @ts-ignore
                window.runtime.EventsOff("proxy:request");
            }
        };
    }, [paused]);

    const getStatusColor = (status: number) => {
        if (status >= 200 && status < 300) return "bg-green-500/10 text-green-500 hover:bg-green-500/20";
        if (status >= 300 && status < 400) return "bg-blue-500/10 text-blue-500 hover:bg-blue-500/20";
        if (status >= 400 && status < 500) return "bg-yellow-500/10 text-yellow-500 hover:bg-yellow-500/20";
        if (status >= 500) return "bg-red-500/10 text-red-500 hover:bg-red-500/20";
        return "bg-gray-500/10 text-gray-500";
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

            <div className="border rounded-md flex-1 overflow-hidden bg-background/50 backdrop-blur">
                <Table>
                    <TableHeader className="sticky top-0 bg-background z-10">
                        <TableRow>
                            <TableHead className="w-[100px]">Time</TableHead>
                            <TableHead className="w-[80px]">Method</TableHead>
                            <TableHead className="w-[80px]">Status</TableHead>
                            <TableHead>URL</TableHead>
                            <TableHead className="w-[100px] text-right">Duration</TableHead>
                        </TableRow>
                    </TableHeader>
                    <TableBody>
                        {events.length === 0 ? (
                            <TableRow>
                                <TableCell colSpan={5} className="text-center h-24 text-muted-foreground">
                                    No requests captured yet.
                                </TableCell>
                            </TableRow>
                        ) : (
                            events.map((evt, i) => (
                                <TableRow key={i} className="font-mono text-xs">
                                    <TableCell className="text-muted-foreground">{evt.time}</TableCell>
                                    <TableCell>
                                        <Badge variant={getMethodColor(evt.method) as any} className="text-[10px] px-1 py-0 h-5">
                                            {evt.method}
                                        </Badge>
                                    </TableCell>
                                    <TableCell>
                                        <Badge variant="outline" className={`border-0 ${getStatusColor(evt.status)}`}>
                                            {evt.status > 0 ? evt.status : '???'}
                                        </Badge>
                                    </TableCell>
                                    <TableCell className="max-w-[400px] truncate" title={evt.url}>
                                        {evt.url}
                                    </TableCell>
                                    <TableCell className="text-right text-muted-foreground">
                                        {evt.duration_ms}ms
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
