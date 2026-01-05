import { useState, useEffect } from 'react';

import { Trash2, Pause, Play, Bug } from "lucide-react";
import { Button } from "@/components/ui/button";

import { WebSocketViewer } from "../components/monitor/WebSocketViewer";
import { BreakpointPanel } from "../components/monitor/BreakpointPanel";
import { PendingRequestDialog } from "../components/monitor/PendingRequestDialog";
import { RequestDetails } from "../components/monitor/RequestDetails";
import { proxy } from "../../wailsjs/go/models";
import { useTraffic, ProxyEvent } from '../contexts/TrafficContext';
import { TrafficTable } from '../components/monitor/TrafficTable';
import { GetBreakpointRules } from "../../wailsjs/go/main/App";

import { FilterBar } from "@/components/monitor/FilterBar";
import { filterEvent, FilterOptions } from "@/lib/filterUtils";

export default function MonitorPage() {
    const { events, wsFrames, pendingReq, paused, setPaused, clearEvents, setPendingReq, addRule, addEvents } = useTraffic();
    const [selectedEvent, setSelectedEvent] = useState<ProxyEvent | null>(null);

    const [filterOptions, setFilterOptions] = useState<FilterOptions>({
        query: "",
        isRegex: false,
        methods: [],
        statusCategory: "All"
    });

    // Breakpoint State (Still local for UI toggle, but rules could be global if needed)
    const [showBreakpoints, setShowBreakpoints] = useState(false);
    const [rules, setRules] = useState<proxy.BreakpointRule[]>([]);

    useEffect(() => {
        // @ts-ignore
        if (window.go) {
            GetBreakpointRules().then((fetched: proxy.BreakpointRule[]) => {
                if (fetched) setRules(fetched);
            });
        }
    }, []);

    // WebSocket State
    const [showWs, setShowWs] = useState(false);

    // Filter First, Then Virtualize
    const filteredEvents = events.filter(evt => filterEvent(evt, filterOptions));

    // Dev: Simulate Load
    const simulateLoad = () => {
        const newEvents: ProxyEvent[] = [];
        const now = new Date();
        for (let i = 0; i < 5000; i++) {
            newEvents.push({
                time: now.toLocaleTimeString(),
                method: ['GET', 'POST', 'PUT', 'DELETE'][Math.floor(Math.random() * 4)],
                url: `http://example.com/api/resource/${i}`,
                status: [200, 201, 400, 404, 500][Math.floor(Math.random() * 5)],
                duration_ms: Math.floor(Math.random() * 1000),
                size: Math.floor(Math.random() * 5000),
                req_headers: {},
                resp_headers: {},
                req_body: '',
                resp_body: '',
                connection_id: `mock-conn-${i}`
            });
        }
        addEvents(newEvents);
    };

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
                    <Button variant="outline" size="sm" onClick={() => clearEvents()}>
                        <Trash2 className="mr-2 h-4 w-4" /> Clear
                    </Button>
                </div>
            </div>

            {showWs && (
                <div className="h-64 shrink-0">
                    <WebSocketViewer frames={wsFrames} onClear={() => { }} />
                </div>
            )}

            {showBreakpoints && (
                <BreakpointPanel
                    rules={rules}
                    onAdd={(pattern) => {
                        const r = new proxy.BreakpointRule({ id: Date.now().toString(), pattern, enabled: true, method: "*" });
                        setRules([...rules, r]);
                        addRule(r);
                    }}
                    onDelete={(id) => setRules(rules.filter(r => r.id !== id))}
                />
            )}

            <div className="flex gap-2">
                <FilterBar options={filterOptions} onChange={setFilterOptions} />
            </div>

            {/* Virtualized Table */}
            <TrafficTable events={filteredEvents} onSelectEvent={setSelectedEvent} onInjectLoad={simulateLoad} />

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
