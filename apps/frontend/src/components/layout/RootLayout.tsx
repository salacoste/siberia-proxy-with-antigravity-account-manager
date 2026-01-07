import { Outlet } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { Toaster } from "@/components/ui/sonner"

export default function RootLayout() {
    return (
        <div className="flex h-screen w-screen bg-background text-foreground overflow-hidden font-sans">
            <Sidebar />
            <main className="flex-1 overflow-auto bg-[#F3F4F6] dark:bg-zinc-950/50">
                <Outlet />
            </main>
            <Toaster />
        </div>
    );
}
