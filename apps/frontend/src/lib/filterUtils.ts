import { ProxyEvent } from '../contexts/TrafficContext';

export interface FilterOptions {
    query: string;
    isRegex: boolean;
    methods: string[]; // e.g., ["GET", "POST"] or empty for all
    statusCategory: string; // "All", "2xx", "3xx", "4xx", "5xx", "Error"
}

export const filterEvent = (evt: ProxyEvent, options: FilterOptions): boolean => {
    // 1. Method Filter
    if (options.methods.length > 0) {
        if (!options.methods.includes(evt.method.toUpperCase())) {
            return false;
        }
    }

    // 2. Status Filter
    if (options.statusCategory !== "All") {
        const status = evt.status;
        switch (options.statusCategory) {
            case "2xx": if (status < 200 || status >= 300) return false; break;
            case "3xx": if (status < 300 || status >= 400) return false; break;
            case "4xx": if (status < 400 || status >= 500) return false; break;
            case "5xx": if (status < 500 || status >= 600) return false; break;
            case "Error": if (status < 400) return false; break; // 4xx and 5xx
        }
    }

    // 3. Query Filter (Regex or Smart Search)
    const query = options.query.trim();
    if (!query) return true;

    if (options.isRegex) {
        try {
            // Case insensitive regex
            const regex = new RegExp(query, 'i');
            return regex.test(evt.url) || regex.test(evt.method) || regex.test(evt.status.toString());
        } catch {
            return false; // Invalid regex matches nothing (or should we match everything? usually nothing gives feedback)
        }
    } else {
        // Smart Search (existing logic: allows !negative, method:GET, etc.)
        // Refactored from MonitorPage.tsx
        const terms = query.split(/\s+/);
        return terms.every(term => {
            let checkTerm = term;
            let isNegative = false;

            if (checkTerm.startsWith('!') || (checkTerm.startsWith('-') && checkTerm.length > 1)) {
                isNegative = true;
                checkTerm = checkTerm.substring(1);
            }

            let match = false;

            // Field filter?
            if (checkTerm.includes(':')) {
                const parts = checkTerm.split(':');
                if (parts.length === 2) {
                    const [key, val] = parts;
                    const v = val.toLowerCase();
                    if (key.toLowerCase() === 'method') match = evt.method.toLowerCase().includes(v);
                    else if (key.toLowerCase() === 'status') match = evt.status.toString().startsWith(v);
                    else if (key.toLowerCase() === 'url' || key === 'host') match = evt.url.toLowerCase().includes(v);
                }
            } else {
                // Standard text match against common fields
                const lowerTerm = checkTerm.toLowerCase();
                match = evt.url.toLowerCase().includes(lowerTerm) ||
                    evt.method.toLowerCase().includes(lowerTerm) ||
                    evt.status.toString().includes(lowerTerm);
            }

            return isNegative ? !match : match;
        });
    }
};
