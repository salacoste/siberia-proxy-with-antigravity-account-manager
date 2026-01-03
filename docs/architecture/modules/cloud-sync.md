# Architecture: Cloud Profile Sync & E2EE

## 1. Overview
This document defines the architecture for **Story-33 (Cloud Profile Sync)**. The goal is to securely synchronize user settings and accounts across devices without exposing sensitive data (tokens, passwords) to the cloud provider.

## 2. Security Strategy: End-to-End Encryption (E2EE)
We adopt a **Zero-Knowledge** architecture. The cloud server acts as a "dumb store" for encrypted blobs and logic, with no ability to decrypt the data.

### 2.1 Encryption Primitives
*   **Algorithm:** AES-256-GCM (Authenticated Encryption).
*   **Key Derivation:** Argon2id (for deriving keys from user passwords).
*   **Serialization:** JSON -> Protocol Buffers (for compact binary storage) or Compressed JSON.

### 2.2 Key Management (Envelope Encryption)
To allow password changes without re-encrypting all data, we use Envelope Encryption:

1.  **Data Encryption Key (DEK)**: A random 32-byte key generated once per profile. This encrypts the actual data.
2.  **Key Encryption Key (KEK)**: Derived from the User's Password via Argon2id. This encrypts the DEK.

**Storage Structure:**
```json
{
  "profile_id": "uuid",
  "salt": "random_salt_for_kdf",
  "wrapped_dek": "AES(DEK, KEK)", // Encrypted Data Key
  "data_blob": "AES(Payload, DEK)", // Encrypted Profile Data
  "timestamp": 1678900000,
  "hash": "sha256_of_blob"
}
```

### 2.3 Workflow
*   **Encryption (Local)**: `Encrypt(Data, DEK) + Encrypt(DEK, Derived(Password))` -> Cloud.
*   **Decryption (Local)**: `Decrypt(WrappedDEK, Derived(Password))` -> `DEK` -> `Decrypt(Blob, DEK)` -> Data.
*   **Password Change**: Re-encrypt `DEK` with new `KEK` (Derived from new password). `data_blob` remains untouched.

## 3. Synchronization Protocol: "Simple-Sync v1"

For MVP (Story-33), we implement a **Last-Write-Wins (LWW)** strategy using highly reliable HTTP Push/Pull interactions. Real-time WebSockets are deferred to v2.

### 3.1 Data Model
The unit of sync is the **entire profile** (Snapshot Sync).
*   *Pros:* Simple, avoids complex merge conflicts in MVP.
*   *Cons:* Higher bandwidth usage (mitigated by compression).

### 3.2 API Endpoints
*   `POST /sync/push`: Client sends `{ timestamp, hash, encrypted_blob }`.
    *   Server accepts if `timestamp > current_server_timestamp`.
    *   Returns `409 Conflict` if server has newer data.
*   `GET /sync/pull`: Client sends `{ last_known_timestamp }`.
    *   Server returns `204 No Content` if client is up to date.
    *   Server returns `200 OK + { encrypted_blob }` if server has newer data.

### 3.3 Conflict Resolution (MVP)
*   **Strategy**: "Server Wins" on timestamp conflict, but client prompted.
*   If `Push` fails with 409:
    1.  Client performs `Pull`.
    2.  Client decrypts remote data.
    3.  Client (UI) asks user: "Conflict detected. Keep Local or Use Cloud?" (Or automerge if trivial).
    4.  Winner is re-encrypted and Pushed with `new_timestamp`.

## 4. Implementation Plan (Story-33)
1.  **Crypto Utils**: Implement `crypto/vault` package (Argon2id, AES-GCM).
2.  **Sync Service**: Implement `sync/manager` (HTTP Client, Polling logic).
3.  **UI**: Add "Sync Settings" panel with "Master Password" setup.
