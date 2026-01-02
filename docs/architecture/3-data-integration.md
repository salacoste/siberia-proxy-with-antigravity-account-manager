# Architecture Part 3: Data & Integration

## 1. Local Database
*   **Engine:** `SQLite` (via `cgo` or `modernc.org/sqlite`).
*   **Encryption:** `SQLCipher` is recommended for `accounts.db`.
*   **Schema (Accounts):**
    ```sql
    CREATE TABLE accounts (
        id TEXT PRIMARY KEY,
        email TEXT,
        refresh_token TEXT, -- Encrypted
        proxy_id TEXT,
        status TEXT
    );
    ```

## 2. External Integration (Critical)
*   **The "Switch" Mechanism:**
    This is the most fragile part of the architecture and requires robust error handling.
    *   **Module:** `siberia/modules/injection`
    *   **Dependency:** `protobuf` (Go generator).
    *   **Flow:**
        1.  User clicks "Switch".
        2.  Frontend calls `SwitchAccount(id)`.
        3.  Backend performs:
            *   Atomic DB read of Account credentials.
            *   Process discovery (ps list) -> Kill.
            *   File Lock -> Backup -> Edit `state.vscdb`.
            *   Process Spawn.

## 3. Configuration
*   **Viper:** Use `spf13/viper` for managing `config.json`.
*   **Config Structure:**
    ```json
    {
        "port": 8080,
        "theme": "dark",
        "proxy": {
             "upstream": "...",
             "mode": "strict"
        }
    }
    ```
