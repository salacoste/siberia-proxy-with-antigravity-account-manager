import { useRef, useEffect, useState } from 'react';
import { useVirtualizer } from '@tanstack/react-virtual';
import { Table, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Eye, ArrowDown } from "lucide-react";
import { ProxyEvent } from '../../contexts/TrafficContext';

interface TrafficTableProps {
    events: ProxyEvent[];
    onSelectEvent: (evt: ProxyEvent) => void;
    onInjectLoad?: () => void; // For testing
}

export function TrafficTable({ events, onSelectEvent, onInjectLoad }: TrafficTableProps) {
    const parentRef = useRef<HTMLDivElement>(null);
    const [autoScroll, setAutoScroll] = useState(true);

    const rowVirtualizer = useVirtualizer({
        count: events.length,
        getScrollElement: () => parentRef.current,
        estimateSize: () => 35, // Fixed height for compact density
        overscan: 10,
    });

    // Auto-scroll logic
    useEffect(() => {
        if (autoScroll && events.length > 0) {
            rowVirtualizer.scrollToIndex(0, { align: 'start' });
        }
    }, [events.length, autoScroll, rowVirtualizer]);

    // Detect user scroll to disable auto-scroll
    const handleScroll = (e: React.UIEvent<HTMLDivElement>) => {
        const { scrollTop } = e.currentTarget;
        if (scrollTop > 50) { // If user scrolled down (since index 0 is top)
            // Wait, usually index 0 is top.
            // If we prepend new events, they are at index 0.
            // So "Auto Scroll" means staying at Top?
            // Yes, traffic logs usually show newest at top.
            // So if scrollTop > 0, we are looking at older logs.
            setAutoScroll(false);
        } else {
            setAutoScroll(true);
        }
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

    return (
        <div className="flex-1 relative border rounded-md bg-background/50 backdrop-blur overflow-hidden flex flex-col">
            {/* Header is Sticky OUTSIDE the virtual scroller */}
            <div className="border-b bg-muted/50 z-10">
                <Table>
                    <TableHeader className="bg-transparent">
                        <TableRow className="hover:bg-transparent border-none">
                            <TableHead className="w-[100px]">Time</TableHead>
                            <TableHead className="w-[80px]">Method</TableHead>
                            <TableHead className="w-[100px]">Status</TableHead>
                            <TableHead>URL</TableHead>
                            <TableHead className="w-[100px] text-right">Size</TableHead>
                            <TableHead className="w-[100px] text-right">Duration</TableHead>
                            <TableHead className="w-[60px]"></TableHead>
                        </TableRow>
                    </TableHeader>
                </Table>
            </div>

            {/* Auto-Scroll Indicator */}
            {!autoScroll && (
                <div className="absolute bottom-4 right-4 z-20">
                    <Button size="sm" onClick={() => {
                        setAutoScroll(true);
                        rowVirtualizer.scrollToIndex(0);
                    }}>
                        <ArrowDown className="mr-2 h-4 w-4" /> Latest
                    </Button>
                </div>
            )}

            {/* Load Test Button (Dev Only) */}
            {onInjectLoad && (
                <div className="absolute bottom-4 left-4 z-20">
                    <Button size="sm" variant="outline" onClick={onInjectLoad}>
                        Simulate 5k
                    </Button>
                </div>
            )}

            <div
                ref={parentRef}
                className="flex-1 overflow-auto"
                onScroll={handleScroll}
            >
                <div
                    style={{
                        height: `${rowVirtualizer.getTotalSize()}px`,
                        width: '100%',
                        position: 'relative',
                    }}
                >
                    {rowVirtualizer.getVirtualItems().map((virtualRow) => {
                        const evt = events[virtualRow.index];
                        return (
                            <div
                                key={virtualRow.key}
                                style={{
                                    position: 'absolute',
                                    top: 0,
                                    left: 0,
                                    width: '100%',
                                    height: `${virtualRow.size}px`,
                                    transform: `translateY(${virtualRow.start}px)`,
                                }}
                                className="flex items-center border-b hover:bg-muted/50 transition-colors cursor-pointer text-xs font-mono"
                                onClick={() => onSelectEvent(evt)}
                            >
                                {/* We mimic Table Cells with fixed widths */}
                                <div className="w-[100px] px-4 text-muted-foreground whitespace-nowrap overflow-hidden">{evt.time}</div>
                                <div className="w-[80px] px-4">
                                    <Badge variant={getMethodColor(evt.method) as any} className="text-[10px] px-2 py-0.5 h-5 font-bold">
                                        {evt.method}
                                    </Badge>
                                </div>
                                <div className="w-[100px] px-4">
                                    <Badge variant="outline" className={`border ${getStatusColor(evt.status)}`}>
                                        {evt.status > 0 ? evt.status : '---'}
                                    </Badge>
                                </div>
                                <div className="flex-1 px-4 truncate" title={evt.url}>{evt.url}</div>
                                <div className="w-[100px] px-4 text-right text-muted-foreground">{evt.size} B</div>
                                <div className="w-[100px] px-4 text-right font-medium">{evt.duration_ms}ms</div>
                                <div className="w-[60px] px-4 text-center">
                                    <Button variant="ghost" size="icon" className="h-6 w-6">
                                        <Eye className="h-4 w-4" />
                                    </Button>
                                </div>
                            </div>
                        );
                    })}

                    {events.length === 0 && (
                        <div className="flex items-center justify-center p-8 text-muted-foreground">
                            No requests captured yet.
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
