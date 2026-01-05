import { describe, it, expect } from 'vitest';
import { filterEvent, FilterOptions } from './filterUtils';
import { ProxyEvent } from '../contexts/TrafficContext';

// Mock Event Factory
const createEvent = (overrides: Partial<ProxyEvent> = {}): ProxyEvent => ({
    time: "10:00:00",
    method: "GET",
    url: "http://example.com/api/test",
    status: 200,
    duration_ms: 100,
    size: 500,
    req_headers: {},
    resp_headers: {},
    req_body: "",
    resp_body: "",
    ...overrides
});

// Default Options
const defaultOpts: FilterOptions = {
    query: "",
    isRegex: false,
    methods: [],
    statusCategory: "All"
};

describe('filterUtils', () => {
    it('should match everything with empty options', () => {
        expect(filterEvent(createEvent(), defaultOpts)).toBe(true);
    });

    it('should filter by Method', () => {
        const getEvt = createEvent({ method: "GET" });
        const postEvt = createEvent({ method: "POST" });

        const opts: FilterOptions = { ...defaultOpts, methods: ["POST"] };

        expect(filterEvent(getEvt, opts)).toBe(false);
        expect(filterEvent(postEvt, opts)).toBe(true);
    });

    it('should filter by Status Category', () => {
        const success = createEvent({ status: 200 });
        const error = createEvent({ status: 500 });

        const opts2xx: FilterOptions = { ...defaultOpts, statusCategory: "2xx" };
        const optsErr: FilterOptions = { ...defaultOpts, statusCategory: "Error" };

        expect(filterEvent(success, opts2xx)).toBe(true);
        expect(filterEvent(error, opts2xx)).toBe(false);

        expect(filterEvent(success, optsErr)).toBe(false);
        expect(filterEvent(error, optsErr)).toBe(true);
    });

    it('should filter by basic query (partial match)', () => {
        const evt = createEvent({ url: "http://example.com/api/users" });
        const opts: FilterOptions = { ...defaultOpts, query: "users" };

        expect(filterEvent(evt, opts)).toBe(true);
        expect(filterEvent(evt, { ...defaultOpts, query: "products" })).toBe(false);
    });

    it('should filter by negative query', () => {
        const evt = createEvent({ url: "http://example.com/api/css" });
        const opts: FilterOptions = { ...defaultOpts, query: "!css" };

        expect(filterEvent(evt, opts)).toBe(false);
        expect(filterEvent(createEvent({ url: "api/js" }), opts)).toBe(true);
    });

    it('should filter by field syntax (method:POST)', () => {
        const evt = createEvent({ method: "POST" });
        const opts: FilterOptions = { ...defaultOpts, query: "method:POST" };

        expect(filterEvent(evt, opts)).toBe(true);
        expect(filterEvent(createEvent({ method: "GET" }), opts)).toBe(false);
    });

    it('should filter by Regex when enabled', () => {
        const evt = createEvent({ url: "http://example.com/api/v1/users/123" });

        // Match numbers at end
        const opts: FilterOptions = { ...defaultOpts, query: "users\\/\\d+$", isRegex: true };
        expect(filterEvent(evt, opts)).toBe(true);

        const optsFail: FilterOptions = { ...defaultOpts, query: "users\\/abc$", isRegex: true };
        expect(filterEvent(evt, optsFail)).toBe(false);
    });

    it('should handle invalid regex gracefully', () => {
        const evt = createEvent();
        const opts: FilterOptions = { ...defaultOpts, query: "[", isRegex: true }; // Invalid regex

        expect(filterEvent(evt, opts)).toBe(false); // Should not throw
    });
});
