import React from 'react';
import { HashRouter, Routes, Route, Navigate } from 'react-router-dom';
import RootLayout from './components/layout/RootLayout';
import DashboardPage from './pages/DashboardPage';
import AccountsPage from './pages/AccountsPage';
import MonitorPage from './pages/MonitorPage';
import ProxyPage from './pages/ProxyPage';
import SettingsPage from './pages/SettingsPage';
import { useConfigStore } from './stores/useConfigStore';
import { useTheme } from './stores/useTheme';
import { TrafficProvider } from './contexts/TrafficContext';

function App() {
    // Initialize config and theme
    useTheme();
    React.useEffect(() => {
        useConfigStore.getState().fetchConfig();
    }, []);

    return (
        <TrafficProvider>
            <HashRouter>
                <Routes>
                    <Route element={<RootLayout />}>
                        <Route path="/" element={<DashboardPage />} />
                        <Route path="accounts" element={<AccountsPage />} />
                        <Route path="accounts/add" element={<AccountsPage />} /> {/* Handle add via dialog on same page or separate if desired. We switched to Dialog so this might be redundant but keeping for safety */}
                        <Route path="monitor" element={<MonitorPage />} />
                        <Route path="proxy" element={<ProxyPage />} />
                        <Route path="/settings" element={<SettingsPage />} />
                        {/* Fallback */}
                        <Route path="*" element={<Navigate to="/" replace />} />
                    </Route>
                </Routes>
            </HashRouter>
        </TrafficProvider>
    );
}

export default App;
