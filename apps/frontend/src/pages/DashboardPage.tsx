import { Activity, ArrowDown, ArrowUp, Globe, Server } from "lucide-react";
import { WidgetCard } from "@/components/dashboard/WidgetCard";
import { AnalyticsProvider, useAnalytics } from "@/components/analytics/AnalyticsProvider";
import { TrafficChart } from "@/components/analytics/TrafficChart";
import { StatusChart } from "@/components/analytics/StatusChart";
import { ProtocolChart } from "@/components/analytics/ProtocolChart";
import { formatBytes } from "@/lib/utils";

function DashboardContent() {
    const { stats } = useAnalytics();

    return (
        <div className="p-6 space-y-6 max-w-7xl mx-auto">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Dashboard</h1>
                    <p className="text-muted-foreground mt-1">Real-time traffic intelligence.</p>
                </div>
                <div className="flex items-center gap-2 text-sm font-medium text-emerald-500 bg-emerald-500/10 px-3 py-1 rounded-full border border-emerald-500/20 shadow-sm">
                    <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse shadow-[0_0_8px_rgba(16,185,129,0.5)]" />
                    System Online
                </div>
            </div>

            {/* Quick Stats Grid */}
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                <WidgetCard title="Request Rate" className="border-l-4 border-l-blue-500 min-h-[140px]">
                    <div className="flex items-center justify-between mt-2">
                        <div className="text-lg sm:text-xl md:text-2xl lg:text-3xl font-bold" title={`${stats?.rps.toFixed(1)} req/s`}>
                            {stats?.rps.toFixed(1) || "0.0"} <span className="text-xs md:text-sm font-normal text-muted-foreground">req/s</span>
                        </div>
                        <Activity className="h-4 w-4 text-blue-500 flex-shrink-0 ml-2" />
                    </div>
                </WidgetCard>

                <WidgetCard title="Bandwidth In" className="border-l-4 border-l-emerald-500 min-h-[140px]">
                    <div className="flex items-center justify-between mt-2">
                        <div className="text-lg sm:text-xl md:text-2xl lg:text-3xl font-bold" title={formatBytes(stats?.bandwidth_in_speed || 0)}>
                            {formatBytes(stats?.bandwidth_in_speed || 0)}<span className="text-xs md:text-sm font-normal text-muted-foreground">/s</span>
                        </div>
                        <ArrowDown className="h-4 w-4 text-emerald-500 flex-shrink-0 ml-2" />
                    </div>
                </WidgetCard>

                <WidgetCard title="Bandwidth Out" className="border-l-4 border-l-orange-500 min-h-[140px]">
                    <div className="flex items-center justify-between mt-2">
                        <div className="text-lg sm:text-xl md:text-2xl lg:text-3xl font-bold" title={formatBytes(stats?.bandwidth_out_speed || 0)}>
                            {formatBytes(stats?.bandwidth_out_speed || 0)}<span className="text-xs md:text-sm font-normal text-muted-foreground">/s</span>
                        </div>
                        <ArrowUp className="h-4 w-4 text-orange-500 flex-shrink-0 ml-2" />
                    </div>
                </WidgetCard>

                <WidgetCard title="Active Connections" className="border-l-4 border-l-indigo-500 min-h-[140px]">
                    <div className="flex items-center justify-between mt-2">
                        <div className="text-lg sm:text-xl md:text-2xl lg:text-3xl font-bold" title={String(stats?.active_connections || 0)}>
                            {stats?.active_connections || 0}
                        </div>
                        <Server className="h-4 w-4 text-indigo-500 flex-shrink-0 ml-2" />
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
                                <div key={i} className="flex items-center justify-between text-sm gap-4">
                                    <div className="flex items-center gap-2 min-w-0 flex-1">
                                        <Globe className="h-3 w-3 text-muted-foreground flex-shrink-0" />
                                        <span className="font-medium truncate" title={d.domain}>{d.domain}</span>
                                    </div>
                                    <div className="flex items-center gap-3 flex-shrink-0">
                                        <span className="text-muted-foreground text-xs">{d.count} reqs</span>
                                        <div className="w-16 h-1.5 bg-secondary/50 rounded-full overflow-hidden">
                                            <div
                                                style={{ width: `${Math.min((d.count / (stats?.top_domains[0]?.count || 1)) * 100, 100)}%` }}
                                                className="h-full bg-primary/70 rounded-full"
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
            <div className="grid gap-4 md:grid-cols-2">
                <WidgetCard title="Response Codes" className="col-span-1 min-h-[250px]">
                    <div className="h-[200px] w-full">
                        <StatusChart />
                    </div>
                </WidgetCard>
                <WidgetCard title="Protocol Distribution">
                    <div className="h-[200px] w-full">
                        <ProtocolChart />
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
