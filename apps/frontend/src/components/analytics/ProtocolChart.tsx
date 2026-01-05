import React, { useMemo } from 'react';
import { useAnalytics } from './AnalyticsProvider';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

export const ProtocolChart: React.FC = () => {
    const { stats } = useAnalytics();

    const data = useMemo(() => {
        if (!stats || !stats.protocol_breakdown) return [];
        const total = Object.values(stats.protocol_breakdown).reduce((a, b) => a + b, 0);

        return Object.entries(stats.protocol_breakdown)
            .map(([name, value]) => ({
                name,
                value,
                percentage: total > 0 ? (value / total) * 100 : 0
            }))
            .sort((a, b) => b.value - a.value); // Descending
    }, [stats]);

    return (
        <Card className="h-full flex flex-col border-none shadow-none bg-transparent">
            <CardHeader className="pb-2 px-0 pt-0">
                <CardTitle className="text-sm font-medium text-muted-foreground hidden">Protocol Distribution</CardTitle>
            </CardHeader>
            <CardContent className="flex-1 min-h-0 px-0 pb-0">
                {data.length === 0 ? (
                    <div className="flex items-center justify-center h-full text-muted-foreground text-sm">
                        No Data
                    </div>
                ) : (
                    <ScrollArea className="h-full w-full pr-2">
                        <div className="space-y-4 pt-1">
                            {data.map((item) => (
                                <div key={item.name} className="space-y-1">
                                    <div className="flex items-center justify-between text-sm">
                                        <span className="font-medium text-foreground/90">{item.name}</span>
                                        <span className="text-muted-foreground text-xs">{item.value.toLocaleString()} ({item.percentage.toFixed(1)}%)</span>
                                    </div>
                                    <div className="h-2 w-full bg-secondary/30 rounded-full overflow-hidden">
                                        <div
                                            className="h-full bg-blue-500 rounded-full transition-all duration-500 ease-out"
                                            style={{ width: `${item.percentage}%` }}
                                        />
                                    </div>
                                </div>
                            ))}
                        </div>
                    </ScrollArea>
                )}
            </CardContent>
        </Card>
    );
};
