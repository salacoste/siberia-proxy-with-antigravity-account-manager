import React, { useEffect, useState } from 'react';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useAnalytics } from './AnalyticsProvider';

interface DataPoint {
    time: string;
    rps: number;
}

const MAX_POINTS = 60; // 60 seconds history

export const TrafficChart: React.FC = () => {
    const { stats } = useAnalytics();
    const [data, setData] = useState<DataPoint[]>([]);

    useEffect(() => {
        if (!stats) return;

        setData(prev => {
            const now = new Date();
            const timeStr = `${now.getHours()}:${now.getMinutes()}:${now.getSeconds()}`;
            const newPoint = { time: timeStr, rps: stats.rps };
            const newData = [...prev, newPoint];
            if (newData.length > MAX_POINTS) {
                return newData.slice(newData.length - MAX_POINTS);
            }
            return newData;
        });
    }, [stats]);

    return (
        <Card className="h-full flex flex-col">
            <CardHeader className="pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">Requests per Second</CardTitle>
            </CardHeader>
            <CardContent className="flex-1 min-h-0">
                <ResponsiveContainer width="100%" height="100%">
                    <AreaChart data={data}>
                        <defs>
                            <linearGradient id="colorRps" x1="0" y1="0" x2="0" y2="1">
                                <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.8} />
                                <stop offset="95%" stopColor="#3b82f6" stopOpacity={0} />
                            </linearGradient>
                        </defs>
                        <XAxis
                            dataKey="time"
                            hide={true}
                        />
                        <YAxis
                            allowDecimals={false}
                            width={30}
                            tick={{ fontSize: 12 }}
                            axisLine={false}
                            tickLine={false}
                        />
                        <Tooltip
                            contentStyle={{ backgroundColor: 'var(--background)', borderColor: 'var(--border)' }}
                            itemStyle={{ color: 'var(--foreground)' }}
                        />
                        <Area
                            type="monotone"
                            dataKey="rps"
                            stroke="#3b82f6"
                            fillOpacity={1}
                            fill="url(#colorRps)"
                            isAnimationActive={false}
                        />
                    </AreaChart>
                </ResponsiveContainer>
            </CardContent>
        </Card>
    );
};
