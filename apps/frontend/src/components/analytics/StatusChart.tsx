import React, { useMemo } from 'react';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
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

        return Object.entries(stats.response_codes).map(([key, value]) => ({
            name: key,
            value: value
        })).filter(item => item.value > 0);
    }, [stats]);

    return (
        <Card className="h-full flex flex-col">
            <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">Response Codes</CardTitle>
            </CardHeader>
            <CardContent className="flex-1 min-h-0">
                {data.length === 0 ? (
                    <div className="flex items-center justify-center h-full text-muted-foreground text-sm">
                        No Data
                    </div>
                ) : (
                    <ResponsiveContainer width="100%" height="100%">
                        <PieChart>
                            <Pie
                                data={data}
                                cx="50%"
                                cy="50%"
                                innerRadius={60}
                                outerRadius={80}
                                paddingAngle={5}
                                dataKey="value"
                            >
                                {data.map((entry, index) => {
                                    const color = COLORS[entry.name as keyof typeof COLORS] || COLORS['other'];
                                    return <Cell key={`cell-${index}`} fill={color} stroke="none" />;
                                })}
                            </Pie>
                            <Tooltip
                                contentStyle={{ backgroundColor: 'var(--background)', borderColor: 'var(--border)' }}
                                itemStyle={{ color: 'var(--foreground)' }}
                            />
                            <Legend verticalAlign="bottom" height={36} iconType="circle" />
                        </PieChart>
                    </ResponsiveContainer>
                )}
            </CardContent>
        </Card>
    );
};
