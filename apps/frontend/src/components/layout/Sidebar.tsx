
import { NavLink } from 'react-router-dom';
import { LayoutDashboard, Users, Server, Settings } from 'lucide-react';
import { clsx } from 'clsx';
import { twMerge } from 'tailwind-merge';
import { Button } from '@/components/ui/button';

function cn(...inputs: (string | undefined | null | false)[]) {
    return twMerge(clsx(inputs));
}

const NAV_ITEMS = [
    { name: 'Dashboard', path: '/', icon: LayoutDashboard },
    { name: 'Accounts', path: '/accounts', icon: Users },
    { name: 'Proxy', path: '/proxy', icon: Server },
    { name: 'Settings', path: '/settings', icon: Settings },
];

export function Sidebar() {
    return (
        <aside className="w-64 bg-slate-900 text-slate-300 flex flex-col h-full border-r border-slate-800">
            <div className="p-6">
                <h1 className="text-xl font-bold text-white tracking-wider flex items-center gap-2">
                    <span className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">S</span>
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
                                    isActive ? "bg-blue-600/10 text-blue-400 hover:bg-blue-600/20 hover:text-blue-300" : "text-slate-400 hover:text-white hover:bg-slate-800"
                                )}
                            >
                                <item.icon size={20} />
                                {item.name}
                            </Button>
                        )}
                    </NavLink>
                ))}
            </nav>

            <div className="p-4 border-t border-slate-800">
                <div className="text-xs text-slate-500 text-center">
                    v0.1.0-alpha
                </div>
            </div>
        </aside>
    );
}
