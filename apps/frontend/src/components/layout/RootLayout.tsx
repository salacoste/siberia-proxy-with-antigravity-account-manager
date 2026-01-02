import React from 'react';
import { Outlet } from 'react-router-dom';
import { Sidebar } from './Sidebar';

export default function RootLayout() {
    return (
        <div className="flex h-screen w-screen bg-slate-950 text-slate-200 overflow-hidden">
            <Sidebar />
            <main className="flex-1 overflow-auto bg-slate-950">
                <Outlet />
            </main>
        </div>
    );
}
