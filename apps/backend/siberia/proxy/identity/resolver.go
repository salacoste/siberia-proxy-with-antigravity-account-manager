package identity

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const UpstreamURL = "https://cloudcode-pa.googleapis.com/v1internal:loadCodeAssist"

type loadCodeAssistRequest struct {
	Metadata struct {
		IdeType string `json:"ideType"`
	} `json:"metadata"`
}

type loadCodeAssistResponse struct {
	CloudAICompanionProject string `json:"cloudaicompanionProject"`
}

// Resolver handles fetching the internal Project ID.
type Resolver struct {
	Client *http.Client
}

// NewResolver creates a new Resolver with a default timeout client.
func NewResolver() *Resolver {
	return &Resolver{
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchProjectID calls the Google internal API to get the cloudaicompanionProject ID.
func (r *Resolver) FetchProjectID(accessToken string) (string, error) {
	reqBody := loadCodeAssistRequest{}
	reqBody.Metadata.IdeType = "ANTIGRAVITY"

	// Create JSON body (simplified for brevity, normally marshal)
	bodyStr := `{"metadata":{"ideType":"ANTIGRAVITY"}}`

	req, err := http.NewRequest("POST", UpstreamURL, strings.NewReader(bodyStr))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "antigravity/1.0.0")

	resp, err := r.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("upstream error: %s", resp.Status)
	}

	var data loadCodeAssistResponse
	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return "", err
	}

	if data.CloudAICompanionProject == "" {
		return "", fmt.Errorf("no project id in response")
	}

	return data.CloudAICompanionProject, nil
}

// GenerateMockProjectID creates a fake project ID for fallback/shadow mode.
// Format: {adj}-{noun}-{5chars}
func GenerateMockProjectID() string {
	adjectives := []string{"swift", "calm", "bold", "bright"}
	nouns := []string{"core", "flow", "wave", "spark"}

	rand.Seed(time.Now().UnixNano())
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]

	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	suffix := make([]byte, 5)
	for i := range suffix {
		suffix[i] = chars[rand.Intn(len(chars))]
	}

	return fmt.Sprintf("%s-%s-%s", adj, noun, string(suffix))
}
