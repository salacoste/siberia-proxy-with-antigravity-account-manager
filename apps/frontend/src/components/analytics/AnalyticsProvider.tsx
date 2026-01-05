import React, { createContext, useContext, useEffect, useState } from 'react';
import { analytics } from '../../../wailsjs/go/models';
import { GetStats } from '../../../wailsjs/go/analytics/AnalyticsService';

interface AnalyticsContextType {
    stats: analytics.AnalyticsSnapshot | null;
    isLoading: boolean;
    error: string | null;
}

const AnalyticsContext = createContext<AnalyticsContextType>({
    stats: null,
    isLoading: true,
    error: null,
});

export const useAnalytics = () => useContext(AnalyticsContext);

export const AnalyticsProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [stats, setStats] = useState<analytics.AnalyticsSnapshot | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        let isMounted = true;
        const fetchStats = async () => {
            try {
                const data = await GetStats();
                if (isMounted) {
                    setStats(data);
                    setError(null);
                }
            } catch (err) {
                if (isMounted) {
                    console.error("Failed to fetch stats:", err);
                    // Don't set error on every poll failure to avoid UI flickering, 
                    // just log it. Maybe set error if it persists.
                    // setError("Failed to fetch analytics");
                }
            } finally {
                if (isMounted) setIsLoading(false);
            }
        };

        // Initial fetch
        fetchStats();

        // Poll every 1 second
        const interval = setInterval(fetchStats, 1000);

        return () => {
            isMounted = false;
            clearInterval(interval);
        };
    }, []);

    return (
        <AnalyticsContext.Provider value={{ stats, isLoading, error }}>
            {children}
        </AnalyticsContext.Provider>
    );
};
