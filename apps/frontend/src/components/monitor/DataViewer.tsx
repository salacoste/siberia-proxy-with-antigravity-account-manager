import { ScrollArea } from "@/components/ui/scroll-area";

interface DataViewerProps {
    data: string;
    type?: 'json' | 'text';
}

export function DataViewer({ data, type = 'text' }: DataViewerProps) {
    if (!data) {
        return <div className="text-muted-foreground italic text-xs p-4 text-center">Empty</div>;
    }

    let displayContent = data;
    if (type === 'json') {
        try {
            const parsed = JSON.parse(data);
            displayContent = JSON.stringify(parsed, null, 2);
        } catch {
            // fallback to text if invalid json
        }
    }

    return (
        <ScrollArea className="h-full w-full rounded-md border bg-muted/50 p-4">
            <pre className="font-mono text-[10px] whitespace-pre-wrap break-all">
                {displayContent}
            </pre>
        </ScrollArea>
    );
}
