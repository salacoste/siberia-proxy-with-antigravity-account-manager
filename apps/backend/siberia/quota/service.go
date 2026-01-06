package quota

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/salacoste/siberia/siberia/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	CloudCodeBaseURL = "https://cloudcode-pa.googleapis.com"
	UserAgent        = "antigravity/1.11.3 Darwin/arm64"
)

type Service struct {
	client *http.Client
	ctx    context.Context
}

func NewService() *Service {
	return &Service{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// FetchAccountStats performs the full flow: LoadProject -> FetchQuota -> Aggregate Stats
func (s *Service) FetchAccountStats(accessToken string, email string) (*Stats, error) {
	// 1. Get Project & Tier
	projectID, tierID, err := s.loadProject(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to load project: %w", err)
	}

	// 2. Fetch Quota
	modelUsage, err := s.fetchQuotaInternal(accessToken, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quota: %w", err)
	}

	// 3. Construct Stats
	stats := &Stats{
		Tier:        tierID,
		ProjectID:   projectID,
		Models:      modelUsage,
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	logger.New("Quota").Info(fmt.Sprintf("[%s] Tier: %s, Models: %v", email, tierID, modelUsage))

	// Emit Event if context is available
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "quota:update", stats)
	}

	return stats, nil
}

// SetContext sets the context for event emission
func (s *Service) SetContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *Service) loadProject(accessToken string) (string, string, error) {
	url := fmt.Sprintf("%s/v1internal:loadCodeAssist", CloudCodeBaseURL)
	payload := map[string]interface{}{
		"metadata": map[string]string{
			"ideType": "ANTIGRAVITY",
		},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return "", "", fmt.Errorf("403 Forbidden")
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("API error: %s", resp.Status)
	}

	var data LoadProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", err
	}

	// Determine Tier: PaidTier > CurrentTier
	tierID := "free" // default
	if data.PaidTier != nil && data.PaidTier.ID != "" {
		tierID = data.PaidTier.ID
	} else if data.CurrentTier != nil && data.CurrentTier.ID != "" {
		tierID = data.CurrentTier.ID
	}

	return data.ProjectID, tierID, nil
}

func (s *Service) fetchQuotaInternal(accessToken, projectID string) (map[string]int32, error) {
	url := fmt.Sprintf("%s/v1internal:fetchAvailableModels", CloudCodeBaseURL)
	payload := map[string]string{
		"project": projectID,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Quota API error: %s", resp.Status)
	}

	var data FetchQuotaResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	usage := make(map[string]int32)
	for name, info := range data.Models {
		// Filter for relevant models
		if strings.Contains(name, "gemini") || strings.Contains(name, "claude") {
			if info.QuotaInfo != nil {
				percentage := int32(info.QuotaInfo.RemainingFraction * 100)
				usage[name] = percentage
			}
		}
	}

	return usage, nil
}
