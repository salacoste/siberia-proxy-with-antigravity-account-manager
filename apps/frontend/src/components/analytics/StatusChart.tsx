import React, { useMemo } from 'react';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { useAnalytics } from './AnalyticsProvider';

const COLORS = {
    '2xx': '#22c55e', // green-500
    '3xx': '#3b82f6', // blue-500
    '4xx': '#eab308', // yellow-500
    '5xx': '#ef4444', // red-500
    'other': '#94a3b8' // slate-400
};

export const StatusChart: React.FC = () => {
    const { stats } = useAnalytics();

    const data = useMemo(() => {
        if (!stats || !stats.response_codes) return [];
        return Object.entries(stats.response_codes)
            .map(([key, value]) => ({
                name: key,
                value: value
            }))
            .filter(item => item.value > 0)
            .sort((a, b) => b.value - a.value); // Sort by value desc
    }, [stats]);

    return (
        <Card className="h-full flex flex-col border-none shadow-none bg-transparent">
            <CardHeader className="pb-2 px-0 pt-0">
            </CardHeader>
            <CardContent className="flex-1 min-h-0 px-0 pb-0 flex flex-col">
                {data.length === 0 ? (
                    <div className="flex items-center justify-center flex-1 h-full text-muted-foreground text-sm">
                        No Data
                    </div>
                ) : (
                    <>
                        {/* Chart Side - Top */}
                        <div className="flex-1 w-full min-h-[140px] relative">
                            <ResponsiveContainer width="100%" height="100%">
                                <PieChart>
                                    <Pie
                                        data={data}
                                        cx="50%"
                                        cy="50%"
                                        innerRadius={50}
                                        outerRadius={70}
                                        paddingAngle={4}
                                        dataKey="value"
                                        stroke="none"
                                    >
                                        {data.map((entry, index) => {
                                            const color = COLORS[entry.name as keyof typeof COLORS] || COLORS['other'];
                                            return <Cell key={`cell-${index}`} fill={color} />;
                                        })}
                                    </Pie>
                                    <Tooltip
                                        contentStyle={{
                                            backgroundColor: 'hsl(var(--popover))',
                                            borderColor: 'hsl(var(--border))',
                                            borderRadius: 'var(--radius)',
                                            color: 'hsl(var(--popover-foreground))',
                                            boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)'
                                        }}
                                        itemStyle={{ color: 'hsl(var(--foreground))' }}
                                    />
                                </PieChart>
                            </ResponsiveContainer>
                        </div>

                        {/* Legend Side - Bottom */}
                        <div className="w-full mt-2 grid grid-cols-2 gap-x-4 gap-y-1 text-sm">
                            {data.map((entry) => {
                                const color = COLORS[entry.name as keyof typeof COLORS] || COLORS['other'];
                                return (
                                    <div key={entry.name} className="flex items-center justify-between">
                                        <div className="flex items-center gap-2 truncate">
                                            <div className="w-2.5 h-2.5 rounded-full flex-shrink-0" style={{ backgroundColor: color }} />
                                            <span className="font-medium truncate">{entry.name}</span>
                                        </div>
                                        <span className="text-muted-foreground ml-1">{entry.value}</span>
                                    </div>
                                );
                            })}
                        </div>
                    </>
                )}
            </CardContent>
        </Card>
    );
};
