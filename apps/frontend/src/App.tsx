import React from 'react';
import { HashRouter, Routes, Route, Navigate } from 'react-router-dom';
import RootLayout from './components/layout/RootLayout';
import DashboardPage from './pages/DashboardPage';
import AccountsPage from './pages/AccountsPage';
import ProxyPage from './pages/ProxyPage';
import SettingsPage from './pages/SettingsPage';
import { useConfigStore } from './stores/useConfigStore';

function App() {
    // Initialize config on load
    React.useEffect(() => {
        useConfigStore.getState().fetchConfig();
    }, []);

    return (
        <HashRouter>
            <Routes>
                <Route element={<RootLayout />}>
                    <Route path="/" element={<DashboardPage />} />
                    <Route path="/accounts" element={<AccountsPage />} />
                    <Route path="/proxy" element={<ProxyPage />} />
                    <Route path="/settings" element={<SettingsPage />} />
                    {/* Fallback */}
                    <Route path="*" element={<Navigate to="/" replace />} />
                </Route>
            </Routes>
        </HashRouter>
    );
}

export default App;
