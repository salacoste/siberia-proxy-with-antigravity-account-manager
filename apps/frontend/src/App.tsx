import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import RootLayout from './components/layout/RootLayout';
import DashboardPage from './pages/DashboardPage';
import AccountsPage from './pages/AccountsPage';
import MonitorPage from './pages/MonitorPage';
import ProxyPage from './pages/ProxyPage';
import CloudPage from './pages/CloudPage';
import SettingsPage from './pages/SettingsPage';
import { useConfigStore } from './stores/useConfigStore';
import { useTheme } from './stores/useTheme';
import { TrafficProvider } from './contexts/TrafficContext';
import { SaveWindowSize } from '../wailsjs/go/main/App';

import { ErrorBoundary } from './components/ErrorBoundary';

function App() {
    // Initialize config and theme
    useTheme();
    React.useEffect(() => {
        useConfigStore.getState().fetchConfig();

        // Window Persistence
        let timeout: NodeJS.Timeout;
        const handleResize = () => {
            clearTimeout(timeout);
            timeout = setTimeout(async () => {
                // @ts-ignore
                if (!window.go) return;
                try {
                    await SaveWindowSize(window.outerWidth, window.outerHeight);
                } catch (e) {
                    console.error("Failed to save window size", e);
                }
            }, 1000);
        };
        window.addEventListener('resize', handleResize);
        return () => {
            window.removeEventListener('resize', handleResize);
            clearTimeout(timeout);
        };
    }, []);

    return (
        <ErrorBoundary>
            <TrafficProvider>
                <BrowserRouter>
                    <Routes>
                        <Route element={<RootLayout />}>
                            <Route path="/" element={<Navigate to="/dashboard" replace />} />
                            <Route path="/dashboard" element={<DashboardPage />} />
                            <Route path="/accounts" element={<AccountsPage />} />
                            <Route path="accounts/add" element={<AccountsPage />} />
                            <Route path="monitor" element={<MonitorPage />} />
                            <Route path="proxy" element={<ProxyPage />} />
                            <Route path="/cloud" element={<CloudPage />} />
                            <Route path="/settings" element={<SettingsPage />} />
                            {/* Fallback */}
                            <Route path="*" element={<Navigate to="/" replace />} />
                        </Route>
                    </Routes>
                </BrowserRouter>
            </TrafficProvider>
        </ErrorBoundary>
    );
}

export default App;
