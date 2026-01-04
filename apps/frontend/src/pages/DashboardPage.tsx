import { Activity, ArrowDown, ArrowUp, Globe, Server } from "lucide-react";
import { WidgetCard } from "@/components/dashboard/WidgetCard";
import { AnalyticsProvider, useAnalytics } from "@/components/analytics/AnalyticsProvider";
import { TrafficChart } from "@/components/analytics/TrafficChart";
import { StatusChart } from "@/components/analytics/StatusChart";

function DashboardContent() {
    const { stats } = useAnalytics();

    return (
        <div className="p-6 space-y-6 max-w-7xl mx-auto">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
                    <p className="text-muted-foreground mt-1">Real-time traffic intelligence.</p>
                </div>
                <div className="flex items-center gap-2 text-sm text-muted-foreground bg-white/50 px-3 py-1 rounded-full border border-border/50">
                    <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
                    System Online
                </div>
            </div>

            {/* Quick Stats Grid */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <WidgetCard title="Request Rate" className="border-l-4 border-l-blue-500">
                    <div className="flex items-center justify-between mt-2">
                        <div className="text-2xl font-bold">{stats?.rps.toFixed(1) || "0.0"} <span className="text-xs font-normal text-muted-foreground">req/s</span></div>
                        <Activity className="h-4 w-4 text-blue-500" />
                    </div>
                </WidgetCard>

                <WidgetCard title="Bandwidth In" className="border-l-4 border-l-emerald-500">
                    <div className="flex items-center justify-between mt-2">
                        {/* TODO: Format bytes properly in next Story */}
                        <div className="text-2xl font-bold">{(stats?.bandwidth_in_speed || 0).toFixed(1)} <span className="text-xs font-normal text-muted-foreground">B/s</span></div>
                        <ArrowDown className="h-4 w-4 text-emerald-500" />
                    </div>
                </WidgetCard>

                <WidgetCard title="Bandwidth Out" className="border-l-4 border-l-orange-500">
                    <div className="flex items-center justify-between mt-2">
                        {/* TODO: Format bytes properly in next Story */}
                        <div className="text-2xl font-bold">{(stats?.bandwidth_out_speed || 0).toFixed(1)} <span className="text-xs font-normal text-muted-foreground">B/s</span></div>
                        <ArrowUp className="h-4 w-4 text-orange-500" />
                    </div>
                </WidgetCard>

                <WidgetCard title="Active Connections" className="border-l-4 border-l-indigo-500">
                    <div className="flex items-center justify-between mt-2">
                        <div className="text-2xl font-bold">{stats?.active_connections || 0}</div>
                        <Server className="h-4 w-4 text-indigo-500" />
                    </div>
                </WidgetCard>
            </div>

            {/* Main Content Grid */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-7">
                {/* Traffic Volume (Big Chart) */}
                <WidgetCard title="Traffic Volume" className="col-span-4 min-h-[300px]">
                    <div className="h-[250px] w-full">
                        <TrafficChart />
                    </div>
                </WidgetCard>

                {/* Top Domains (List) */}
                <WidgetCard title="Top Domains" description="Most visited destinations" className="col-span-3">
                    <div className="space-y-4">
                        {(stats?.top_domains || []).length === 0 ? (
                            <div className="text-sm text-muted-foreground text-center py-8">No traffic data yet</div>
                        ) : (
                            stats?.top_domains.map((d, i) => (
                                <div key={i} className="flex items-center justify-between text-sm">
                                    <div className="flex items-center gap-2">
                                        <Globe className="h-3 w-3 text-muted-foreground" />
                                        <span className="font-medium truncate max-w-[200px]" title={d.domain}>{d.domain}</span>
                                    </div>
                                    <div className="flex items-center gap-3">
                                        <span className="text-muted-foreground text-xs">{d.count} reqs</span>
                                        <div className="w-16 h-1.5 bg-muted rounded-full overflow-hidden">
                                            {/* Simple relative bar, assuming max is first item? Analytics engine sorts descending. */}
                                            <div
                                                style={{ width: `${Math.min((d.count / (stats?.top_domains[0]?.count || 1)) * 100, 100)}%` }}
                                                className="h-full bg-primary/70"
                                            />
                                        </div>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </WidgetCard>
            </div>

            {/* Protocols Breakdown */}
            <div className="grid gap-4 md:grid-cols-3">
                <WidgetCard title="Response Codes" className="col-span-1 min-h-[250px]">
                    <div className="h-[200px] w-full">
                        <StatusChart />
                    </div>
                </WidgetCard>
                <WidgetCard title="Protocol Distribution" className="col-span-2">
                    <div className="flex items-center justify-center h-[200px] text-muted-foreground/30 bg-muted/10 rounded-lg border border-dashed">
                        {/* Placeholder for Protocol Breakdown if needed, or we can use another chart later */}
                        Coming Soon
                    </div>
                </WidgetCard>
            </div>
        </div>
    );
}

export default function DashboardPage() {
    return (
        <AnalyticsProvider>
            <DashboardContent />
        </AnalyticsProvider>
    );
}
