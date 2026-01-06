# Architecture: Quota Service & Tier Detection

**Status:** Approved
**Date:** 2026-01-06
**Reference:** Based on Rust `modules/quota.rs` reference implementation.

## 1. Overview
The Quota Service is responsible for fetching usage limits and determining the "Paid Tier" of a Google account. This information allows the "Smart Pooling" logic to prioritize higher-tier accounts (e.g., Gemini Advanced) over Free accounts.

## 2. API Endpoints

### A. Load Project & Tier
**URL:** `https://cloudcode-pa.googleapis.com/v1internal:loadCodeAssist`
**Method:** POST
**Headers:**
- `Authorization`: `Bearer <access_token>`
- `Content-Type`: `application/json`
- `User-Agent`: `antigravity/1.11.3 Darwin/arm64`
**Body:**
```json
{"metadata": {"ideType": "ANTIGRAVITY"}}
```

**Response Schema:**
```json
{
  "cloudaicompanionProject": "project-id-string",
  "currentTier": { "id": "tier-id", "quotaTier": "tier-name" },
  "paidTier": { "id": "tier-id-preferred", "quotaTier": "tier-name" }
}
```

### B. Fetch Usage
**URL:** `https://cloudcode-pa.googleapis.com/v1internal:fetchAvailableModels`
**Method:** POST
**Body:**
```json
{"project": "project-id-string"}
```

**Response Schema:**
```json
{
  "models": {
    "gemini-pro": {
      "quotaInfo": {
        "remainingFraction": 0.9,
        "resetTime": "2024-01-06T12:00:00Z"
      }
    }
  }
}
```

## 3. Go Struct Definitions

```go
package quota

// LoadProjectResponse maps the initial handshake
type LoadProjectResponse struct {
	ProjectID   string `json:"cloudaicompanionProject"`
	CurrentTier *Tier  `json:"currentTier"`
	PaidTier    *Tier  `json:"paidTier"`
}

type Tier struct {
	ID        string `json:"id"`
	QuotaTier string `json:"quotaTier"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
}

// FetchQuotaResponse maps the models usage
type FetchQuotaResponse struct {
	Models map[string]ModelInfo `json:"models"`
}

type ModelInfo struct {
	QuotaInfo *QuotaInfo `json:"quotaInfo"`
}

type QuotaInfo struct {
	RemainingFraction float64 `json:"remainingFraction"` // 0.0 to 1.0
	ResetTime         string  `json:"resetTime"`
}
```

## 4. Logic Flow

1.  **Auth & Tier Check:**
    - Call `loadCodeAssist`.
    - **Tier Detection:** `PaidTier.ID` > `CurrentTier.ID`.
    - If `PaidTier` exists, mark account as **PRO/ULTRA**.
2.  **Quota Fetch:**
    - Use `projectID` from step 1.
    - Call `fetchAvailableModels`.
    - Filter keys containing "gemini" or "claude".
    - Store `RemainingFraction * 100` as percentage.
3.  **Error Handling:**
    - If **403 Forbidden**: Mark account as `Status: Invalid` (or `Forbidden`). Do not retry.

## 5. Storage
Store the result in the `Account` struct (SQLite):
-   `Stats` (JSON field):
    ```json
    {
      "tier": "gemini-advanced",
      "project_id": "bamboo-...",
      "models": { "gemini": 90, "claude": 100 },
      "last_updated": "2026-01-06T..."
    }
    ```
