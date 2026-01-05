import React, { createContext, useContext, useEffect, useState, useRef } from 'react';
import { proxy } from "../../wailsjs/go/models";

// Re-using interface from MonitorPage (should probably be shared)
export interface ProxyEvent {
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
    connection_id: string;
}

interface TrafficContextType {
    events: ProxyEvent[];
    wsFrames: proxy.WebSocketFrame[];
    pendingReq: proxy.PendingRequest | null;
    paused: boolean;
    setPaused: (paused: boolean) => void;
    clearEvents: () => void;
    setPendingReq: (req: proxy.PendingRequest | null) => void;
    addRule: (rule: proxy.BreakpointRule) => void; // Placeholder if needed globally
    addEvents: (newEvents: ProxyEvent[]) => void;
}

const TrafficContext = createContext<TrafficContextType | undefined>(undefined);

const MAX_LOGS = 5000;

export function TrafficProvider({ children }: { children: React.ReactNode }) {
    const [events, setEvents] = useState<ProxyEvent[]>([]);
    const [wsFrames, setWsFrames] = useState<proxy.WebSocketFrame[]>([]);
    const [pendingReq, setPendingReq] = useState<proxy.PendingRequest | null>(null);
    const [paused, setPaused] = useState(false);

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
                    // Prepend new event
                    const next = [event, ...prev];
                    // Slice from end if too large suitable for virtual scroll (index 0 is newest)
                    if (next.length > MAX_LOGS) {
                        return next.slice(0, MAX_LOGS);
                    }
                    return next;
                });
            });

            // @ts-ignore
            window.runtime.EventsOn("breakpoint:hit", (req: proxy.PendingRequest) => {
                setPendingReq(req);
            });

            // @ts-ignore
            window.runtime.EventsOn("proxy:ws:frame", (frame: proxy.WebSocketFrame) => {
                setWsFrames(prev => {
                    const next = [...prev, frame];
                    if (next.length > MAX_LOGS) return next.slice(-MAX_LOGS); // Keep last N for WS
                    return next;
                });
            });

            console.log("[TrafficProvider] Listeners attached.");
        }

        return () => {
            // Cleanup listeners if possible, but Wails listeners are global.
            // Usually we don't need to unregister for a global singleton provider.
        }
    }, []);

    const clearEvents = () => {
        setEvents([]);
        setWsFrames([]);
    };

    // Placeholder
    const addRule = (_rule: proxy.BreakpointRule) => { };

    const addEvents = (newEvents: ProxyEvent[]) => {
        setEvents(prev => {
            const next = [...newEvents, ...prev];
            if (next.length > MAX_LOGS) return next.slice(0, MAX_LOGS);
            return next;
        });
    };

    return (
        <TrafficContext.Provider value={{
            events,
            wsFrames,
            pendingReq,
            paused,
            setPaused,
            clearEvents,
            setPendingReq,
            addRule,
            addEvents
        }}>
            {children}
        </TrafficContext.Provider>
    );
}

export function useTraffic() {
    const context = useContext(TrafficContext);
    if (context === undefined) {
        throw new Error('useTraffic must be used within a TrafficProvider');
    }
    return context;
}
