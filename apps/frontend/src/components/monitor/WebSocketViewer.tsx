import { useState, useEffect, useRef } from 'react';
import { proxy } from '../../../wailsjs/go/models';

interface WebSocketViewerProps {
    frames: proxy.WebSocketFrame[];
    onClear: () => void;
}

export function WebSocketViewer({ frames, onClear }: WebSocketViewerProps) {
    const bottomRef = useRef<HTMLDivElement>(null);
    const [autoScroll, setAutoScroll] = useState(true);

    useEffect(() => {
        if (autoScroll && bottomRef.current) {
            bottomRef.current.scrollIntoView({ behavior: 'smooth' });
        }
    }, [frames, autoScroll]);

    return (
        <div className="flex flex-col h-full bg-slate-900 rounded-lg border border-slate-700 overflow-hidden">
            <div className="flex items-center justify-between px-4 py-2 bg-slate-800 border-b border-slate-700">
                <h3 className="text-sm font-semibold text-slate-200">Live WebSocket Stream</h3>
                <div className="flex gap-2">
                    <button
                        onClick={() => setAutoScroll(!autoScroll)}
                        className={`text-xs px-2 py-1 rounded ${autoScroll ? 'bg-blue-600 text-white' : 'bg-slate-700 text-slate-400'}`}
                    >
                        Auto-scroll
                    </button>
                    <button
                        onClick={onClear}
                        className="text-xs px-2 py-1 rounded bg-slate-700 text-slate-400 hover:text-white"
                    >
                        Clear
                    </button>
                </div>
            </div>

            <div className="flex-1 overflow-y-auto p-2 space-y-1 font-mono text-xs">
                {frames.length === 0 && (
                    <div className="text-center text-slate-500 mt-10">No frames captured yet...</div>
                )}
                {frames.map((frame) => (
                    <div key={frame.id} className={`flex gap-2 p-2 rounded ${frame.direction === 'outgoing' ? 'bg-slate-800/50 ml-4 border-l-2 border-green-500' : 'bg-slate-800/80 mr-4 border-r-2 border-blue-500'}`}>
                        <div className="text-slate-500 min-w-[60px]">{frame.time}</div>
                        <div className={`font-bold uppercase min-w-[40px] ${frame.direction === 'outgoing' ? 'text-green-400' : 'text-blue-400'}`}>
                            {frame.direction === 'outgoing' ? 'UP' : 'DOWN'}
                        </div>
                        <div className="flex-1 break-all text-slate-300">
                            {frame.payload}
                        </div>
                        <div className="text-slate-600 min-w-[40px] text-right">{frame.length}b</div>
                    </div>
                ))}
                <div ref={bottomRef} />
            </div>
        </div>
    );
}
