package zai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// VisionClient handles the communication with Z.ai Vision API
type VisionClient struct {
	BaseURL string
	ApiKey  string
	HTTP    *http.Client
}

// NewVisionClient creates a new client
func NewVisionClient(baseURL, apiKey string) *VisionClient {
	return &VisionClient{
		BaseURL: baseURL,
		ApiKey:  apiKey,
		HTTP: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Message types for GLM-4v
type Message struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

type ContentPart struct {
	Type     string    `json:"type"`                // "text" or "image_url"
	Text     string    `json:"text,omitempty"`      // for text
	ImageURL *ImageURL `json:"image_url,omitempty"` // for image_url
}

type ImageURL struct {
	URL string `json:"url"` // base64 or http url
}

// Response structures
type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *APIError `json:"error,omitempty"`
}

type APIError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// VisionChatCompletion sends a multimodal request to Z.ai
func (c *VisionClient) VisionChatCompletion(messages []Message) (string, error) {
	url := fmt.Sprintf("%s/chat/completions", strings.TrimRight(c.BaseURL, "/"))

	payload := map[string]interface{}{
		"model":    "glm-4v", // Vision model
		"messages": messages,
		"stream":   false,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("api error (status %d): %s", resp.StatusCode, string(body))
	}

	var result ChatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("api error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices returned")
	}

	return result.Choices[0].Message.Content, nil
}

// ImageSourceToContent helper converts a file path or URL to a message content part
func ImageSourceToContent(src string) (ContentPart, error) {
	// 1. Check if URL
	if strings.HasPrefix(src, "http") {
		return ContentPart{
			Type: "image_url",
			ImageURL: &ImageURL{
				URL: src,
			},
		}, nil
	}

	// 2. Assume local file
	data, err := os.ReadFile(src)
	if err != nil {
		return ContentPart{}, fmt.Errorf("read file: %w", err)
	}

	// Check size (simple check, 5MB limit for example)
	if len(data) > 5*1024*1024 {
		return ContentPart{}, fmt.Errorf("image too large (>5MB)")
	}

	// Detect mime type (basic)
	mimeType := http.DetectContentType(data)
	base64Data := base64.StdEncoding.EncodeToString(data)

	return ContentPart{
		Type: "image_url",
		ImageURL: &ImageURL{
			URL: fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data),
		},
	}, nil
}
