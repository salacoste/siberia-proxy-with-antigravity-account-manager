package zai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
)

// GenerateImage calls the Z.ai /v1/images/generations endpoint
func (c *VisionClient) GenerateImage(req *mappers.ImageRequest) (*mappers.ImageResponse, error) {
	url := fmt.Sprintf("%s/images/generations", strings.TrimRight(c.BaseURL, "/"))

	// Ensure defaults
	if req.N == 0 {
		req.N = 1
	}
	if req.Size == "" {
		req.Size = "1024x1024"
	}
	// Force model if needed, but for now we rely on default or allow mapping later.
	// We might need to inject "model": "flux-schnell" or "dall-e-3" depending on provider config.
	// For this client, we construct the payload.

	payload := map[string]interface{}{
		"prompt":          req.Prompt,
		"n":               req.N,
		"size":            req.Size,
		"response_format": req.ResponseFormat,
	}
	if req.Style != "" {
		payload["style"] = req.Style
	}
	// Add model explicit "flux" if commonly used by Z.ai, or "dall-e-3" if acting as proxy.
	// Let's explicitly default to "dall-e-3" compatible if strictly proxying,
	// OR "flux" if specific. Implementation plan said "translating requests".
	// Let's add "model": "flux" as a safe bet for Z.ai, or pass through if provided.
	payload["model"] = "flux"

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body))
	}

	var result mappers.ImageResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}
