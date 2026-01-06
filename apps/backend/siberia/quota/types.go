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

// Stats represents the stored JSON data in the Account struct
type Stats struct {
	Tier        string           `json:"tier"` // "free", "pro", "ultra" etc.
	ProjectID   string           `json:"project_id"`
	Models      map[string]int32 `json:"models"` // Percentage 0-100
	LastUpdated string           `json:"last_updated"`
}
