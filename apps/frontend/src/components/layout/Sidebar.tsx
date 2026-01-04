
import { NavLink } from 'react-router-dom';
import { LayoutDashboard, Users, Server, Settings, Activity } from 'lucide-react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { Button } from '@/components/ui/button';

function cn(...inputs: (string | undefined | null | false)[]) {
    return twMerge(clsx(inputs));
}

const NAV_ITEMS = [
    { name: 'Dashboard', path: '/', icon: LayoutDashboard },
    { name: 'Accounts', path: '/accounts', icon: Users },
    { name: 'Monitor', path: '/monitor', icon: Activity },
    { name: 'Proxy', path: '/proxy', icon: Server },
    { name: 'Settings', path: '/settings', icon: Settings },
];

export function Sidebar() {
    return (
        <aside className="w-64 bg-sidebar text-sidebar-foreground flex flex-col h-full border-r border-border">
            <div className="p-6">
                <h1 className="text-xl font-bold text-sidebar-foreground tracking-wider flex items-center gap-2">
                    <span className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white">S</span>
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
                                variant={isActive ? "secondary" : "ghost"}
                                className={cn(
                                    "w-full justify-start gap-3",
                                    isActive ? "bg-sidebar-accent text-sidebar-accent-foreground" : "text-sidebar-foreground/70 hover:text-sidebar-foreground hover:bg-sidebar-accent/50"
                                )}
                            >
                                <item.icon size={20} />
                                {item.name}
                            </Button>
                        )}
                    </NavLink>
                ))}
            </nav>

            <div className="p-4 border-t border-sidebar-border">
                <div className="text-xs text-sidebar-foreground/50 text-center">
                    v1.2.0
                </div>
            </div>
        </aside>
    );
}
