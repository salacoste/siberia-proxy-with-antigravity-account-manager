import React from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import RootLayout from './components/layout/RootLayout';
import DashboardPage from './pages/DashboardPage';
import AccountsPage from './pages/AccountsPage';
import MonitorPage from './pages/MonitorPage';
import ProxyPage from './pages/ProxyPage';
import CloudPage from './pages/CloudPage';
import SettingsPage from './pages/SettingsPage';
import MapLocalPage from './pages/MapLocalPage';
import ScriptingPage from './pages/ScriptingPage';
import { useConfigStore } from './stores/useConfigStore';

import { useTheme } from './stores/useTheme';

import { TrafficProvider } from './contexts/TrafficContext';
import { SaveWindowSize, ListAccounts } from '../wailsjs/go/main/App';

import { ErrorBoundary } from './components/ErrorBoundary';
import { MigrationWizard } from './components/migration/MigrationWizard';

function App() {
    // Initialize config and theme
    useTheme();
    const [showWizard, setShowWizard] = React.useState(false);

    React.useEffect(() => {
        useConfigStore.getState().fetchConfig();

        // Check for accounts to trigger Wizard
        const checkAccounts = async () => {
            // @ts-ignore
            if (!window.go) return;
            try {
                const accounts = await ListAccounts();
                if (accounts && accounts.length === 0) {
                    setShowWizard(true);
                }
            } catch (e) {
                console.error("Failed to list accounts", e);
            }
        };
        checkAccounts();

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
                            <Route path="/map-local" element={<MapLocalPage />} />
                            <Route path="/scripting" element={<ScriptingPage />} />
                            <Route path="/settings" element={<SettingsPage />} />



                            {/* Fallback */}
                            <Route path="*" element={<Navigate to="/" replace />} />
                        </Route>
                    </Routes>
                </BrowserRouter>
                <MigrationWizard open={showWizard} onOpenChange={setShowWizard} />
            </TrafficProvider>
        </ErrorBoundary>
    );
}

export default App;
