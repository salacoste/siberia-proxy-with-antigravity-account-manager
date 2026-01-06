
import { NavLink } from 'react-router-dom';
import { LayoutDashboard, Users, Server, Settings, Activity, Cloud, FileCode, Code } from 'lucide-react';
import { clsx } from 'clsx';


import { twMerge } from 'tailwind-merge';
import { Button } from '@/components/ui/button';

function cn(...inputs: (string | undefined | null | false)[]) {
    return twMerge(clsx(inputs));
}

const NAV_ITEMS = [
    { name: "Dashboard", icon: LayoutDashboard, path: "/dashboard" },
    { name: 'Accounts', path: '/accounts', icon: Users },
    { name: 'Monitor', path: '/monitor', icon: Activity },
    { name: 'Proxy', path: '/proxy', icon: Server },
    { name: 'Map Local', path: '/map-local', icon: FileCode },
    { name: 'Scripting', path: '/scripting', icon: Code },
    { name: 'Cloud Sync', path: '/cloud', icon: Cloud },


    { name: 'Settings', path: '/settings', icon: Settings },
];

export function Sidebar() {
    return (
        <aside className="w-64 bg-sidebar/50 dark:bg-sidebar/20 text-sidebar-foreground flex flex-col h-full border-r border-border backdrop-blur-sm">
            <div className="p-6">
                <h1 className="text-xl font-bold text-foreground tracking-tight flex items-center gap-2 font-heading">
                    <span className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white shadow-md shadow-blue-500/20">S</span>
                    SIBERIA
                </h1>
            </div>

            <nav className="flex-1 px-4 py-4 space-y-1">
                {NAV_ITEMS.map((item) => (
                    <NavLink
                        key={item.path}
                        to={item.path}
                        className={({ isActive: _ }) =>
                            cn(
                                "w-full", // wrapper for block width
                            )
                        }
                    >
                        {({ isActive }) => (
                            <Button
                                variant="ghost"
                                className={cn(
                                    "w-full justify-start gap-3 h-10 font-medium transition-all px-4",
                                    isActive
                                        ? "bg-white dark:bg-accent text-blue-600 dark:text-blue-400 shadow-sm border border-black/[0.04] dark:border-white/[0.04]"
                                        : "text-muted-foreground hover:text-foreground hover:bg-black/[0.02] dark:hover:bg-white/[0.04]"
                                )}
                            >
                                <item.icon size={18} className={isActive ? "text-blue-600 dark:text-blue-400" : "text-muted-foreground"} />
                                <span className={isActive ? "font-semibold" : ""}>{item.name}</span>
                            </Button>
                        )}
                    </NavLink>
                ))}
            </nav>

            <div className="p-4 border-t border-sidebar-border">
                <div className="text-xs text-muted-foreground text-center font-mono">
                    v1.2.0
                </div>
            </div>
        </aside>
    );
}
