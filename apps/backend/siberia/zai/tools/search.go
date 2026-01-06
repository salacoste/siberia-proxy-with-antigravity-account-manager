package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type SearchRequest struct {
	Query        string   `json:"query"`
	Recency      string   `json:"recency,omitempty"`       // "past_24h", "past_week", "past_month", "past_year"
	DomainFilter []string `json:"domain_filter,omitempty"` // List of domains to include
}

// Result mimics the Z.ai search result structure
type SearchResult struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
	Source  string `json:"source,omitempty"`
}

type SearchResponse struct {
	Results []SearchResult `json:"results"`
}

// CallWebSearchPrime executes a search query against Z.ai's Prime endpoint
func CallWebSearchPrime(apiKey string, baseURL string, req SearchRequest) ([]SearchResult, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("missing z.ai api key")
	}

	url := fmt.Sprintf("%s/search/prime", baseURL)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("X-Antigravity-Client", "siberia-proxy")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("upstream error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("upstream failed with status %d", resp.StatusCode)
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("decode error: %w", err)
	}

	return searchResp.Results, nil
}
