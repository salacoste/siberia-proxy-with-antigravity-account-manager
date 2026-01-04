# Story-42: Migrate to History API Routing (Clean URLs)

**Parent Epic**: Epic-14 (Release Polish)

## Problem
The application currently uses `HashRouter`, resulting in URLs like `http://localhost:7101/#/settings`.
The user requested a cleaner URL structure (`http://localhost:7101/settings`) without the hash symbol.

## Solution
Switch React Router from `HashRouter` to `BrowserRouter`.

## Risks / Considerations
- **Wails Production**: In a compiled Wails app, `BrowserRouter` functionality depends on the asset handler correctly rewrites 404s to `index.html`.
- **Wails Dev**: Should support history API fallback.
- **Vite Standalone**: Supports history API fallback via `--history-api-fallback` (usually default in dev).

## Implementation Steps
1.  Modify `apps/frontend/src/App.tsx` to replace `HashRouter` with `BrowserRouter`.
2.  Verify navigation works in `npm run dev` (Web Mode).
3.  Verify refreshing a sub-page (e.g., `/settings`) does not 404.

## Acceptance Criteria
- [ ] URL is clean: `http://localhost:7101/settings` instead of `.../#/settings`.
- [ ] Navigation between pages works.
- [ ] Refreshing `/settings` loads the page correctly (in Web Mode).
