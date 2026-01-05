# Migration Candidates & Technical Strategy

Based on the comprehensive analysis of `ag-manager-ref`, the following components are identified for migration to the new architecture.

## 1. Core Proxy Engine (High Priority)
The heart of the application. Must be migrated with high fidelity to ensure compatibility.

*   **Source:** `src-tauri/src/proxy/`
*   **Key Components:**
    *   **Server (`server.rs`):** Axum-based server with specific middlewares (Auth, RateLimit, Attribution).
    *   **OpenAI Handler (`handlers/openai.rs`):** Needs to port the streaming logic, tool call array wrapping, and image generation mapping.
    *   **Claude Handler (`handlers/claude.rs`):** Crucial "Thinking" block sanitization and background task downgrading logic.
    *   **Upstream Client (`upstream/client.rs`):** Connection pooling and multi-layer fallback (Prod -> Daily -> Sandbox) logic.
*   **Complexity:** High. Requires careful porting of async logic and error handling.

## 2. Account & Process Management (High Priority)
Required for the proxy to function (it needs tokens).

*   **Source:** `src-tauri/src/modules/`
*   **Key Components:**
    *   **Process Manager (`process.rs`):** Logic to find, kill, and restart the target IDE/Process. Platform-specific (macOS vs Windows).
    *   **DB Injector (`db.rs`):** The Protobuf manipulation logic is critical. Must use the exact same field IDs to avoid breaking the target app.
    *   **Token Manager (`account.rs`):** Logic to refresh tokens and handle rotation.
*   **Complexity:** Medium-High. Platform-specific code and binary data handling.

## 3. OAuth System (Medium Priority)
Can be simplified or replaced, but the reference implementation is robust.

*   **Source:** `src-tauri/src/modules/oauth*.rs`
*   **Key Components:**
    *   **Local Server (`oauth_server.rs`):** Handling the callback.
    *   **Token Exchange (`oauth.rs`):** Standard OAuth code exchange.
*   **Complexity:** Medium. Standard OAuth flow, but needs careful port handling.

## 4. Frontend features (Medium Priority)
The UI needs to providing the switching and monitoring capabilities.

*   **Source:** `src/pages/`, `src/services/`
*   **Key Components:**
    *   **Account List:** Grid/List view with status.
    *   **Proxy Dashboard:** Toggle switch, configuration inputs.
    *   **Code Generator:** The Python snippet generation is a nice-to-have developer feature.
*   **Complexity:** Medium. Standard React/State management.

## 5. Misc / Polish (Low Priority)
*   **Tray Menu (`tray.rs`):** Dynamic menu updates are nice but not strictly "Core" for the MVP.
*   **Quota Checker (`quota.rs`):** Useful for "Paid Tier" detection but can be deferred.
*   **Z.ai Integration:** Only needed if supporting that specific upstream.

## Migration Strategy Recommendation
1.  **Phase 1: The Engine.** Port `proxy` module first. It can run standalone and be tested with `curl`.
2.  **Phase 2: The Injector.** Port `account` and `process` modules to enable real tokens.
3.  **Phase 3: The UI.** Build the Tauri frontend to control the engine.
