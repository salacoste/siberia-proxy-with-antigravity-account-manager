import { Outlet } from 'react-router-dom';
import { Sidebar } from './Sidebar';
import { Toaster } from "@/components/ui/sonner"

export default function RootLayout() {
    return (
        <div className="flex h-screen w-screen bg-background text-foreground overflow-hidden">
            <Sidebar />
            <main className="flex-1 overflow-auto bg-muted/30">
                <Outlet />
            </main>
            <Toaster />
        </div>
    );
}
