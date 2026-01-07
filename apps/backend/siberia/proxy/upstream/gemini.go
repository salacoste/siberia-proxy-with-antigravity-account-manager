package upstream

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/salacoste/siberia/siberia/proxy/mappers"
	"github.com/salacoste/siberia/siberia/proxy/ratelimit"
)

const (
	EndpointPrimary  = "https://cloudcode-pa.googleapis.com"
	EndpointFallback = "https://daily-cloudcode-pa.googleapis.com"
	MaxRetries       = 3
)

type AccountProvider interface {
	GetRotatingToken() (string, string, error) // Returns (token, identity, error)
}

type GeminiClient struct {
	client      *http.Client
	accounts    AccountProvider
	endpoint    string
	primaryURL  string
	fallbackURL string
}

func NewGeminiClient(accProvider AccountProvider) *GeminiClient {
	// Task 1: Optimized Client
	hc := &http.Client{
		Timeout: 60 * time.Second, // Global timeout
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return &GeminiClient{
		client:      hc,
		accounts:    accProvider,
		endpoint:    EndpointPrimary,
		primaryURL:  EndpointPrimary,
		fallbackURL: EndpointFallback,
	}
}

// GenerateContent - Unary
func (c *GeminiClient) GenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (*mappers.GeminiResponse, string, error) {
	urlPath := fmt.Sprintf("/v1/%s:generateContent", model)

	for attempt := 0; attempt < MaxRetries; attempt++ {
		// 1. Get Token
		token, identity, err := c.accounts.GetRotatingToken()
		if err != nil {
			return nil, "", fmt.Errorf("auth error: %v", err)
		}

		// 2. Build Request
		targetURL := c.endpoint + urlPath
		body, _ := json.Marshal(req)
		httpReq, _ := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+token)

		// 3. Execute
		resp, err := c.client.Do(httpReq)
		if err != nil {
			// Network error -> Maybe Fallback endpoint?
			if attempt == MaxRetries-1 {
				return nil, "", err
			}
			time.Sleep(time.Duration(attempt*100) * time.Millisecond) // Linear backoff
			continue
		}
		defer resp.Body.Close()

		// 4. Handle Codes
		if resp.StatusCode == 200 {
			var gResp mappers.GeminiResponse
			if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
				return nil, "", fmt.Errorf("decode error: %v", err)
			}
			return &gResp, identity, nil
		}

		// Retry Logic
		if resp.StatusCode == 429 || resp.StatusCode == 403 || resp.StatusCode == 401 {
			// Read body for intelligent parsing
			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close() // Close early to reuse connection

			rlErr := ratelimit.ParseError(string(bodyBytes))

			// Quota/Auth -> Retry with NEW token
			// Ideally we would respect rlErr.RetryAfter if it's a short rate limit,
			// but for now we just log it and rotate account immediately.
			fmt.Printf("Upstream %d: %s (Wait: %v). Retrying with new account...\n",
				resp.StatusCode, rlErr.Original, rlErr.RetryAfter)

			continue
		}

		if resp.StatusCode >= 500 {
			// Service error -> Switch Endpoint?
			if c.endpoint == c.primaryURL {
				c.endpoint = c.fallbackURL
				fmt.Printf("Upstream 5xx: Switching to Fallback Endpoint\n")
			}
			time.Sleep(time.Duration(1<<attempt) * time.Second) // Expo backoff
			continue
		}

		// Other errors (400 Bad Request) -> Fail immediately
		return nil, "", fmt.Errorf("upstream error %d: %s", resp.StatusCode, resp.Status)
	}

	return nil, "", fmt.Errorf("max retries exceeded")
}

// StreamGenerateContent - Streaming
func (c *GeminiClient) StreamGenerateContent(ctx context.Context, model string, req *mappers.GeminiRequest) (<-chan *mappers.GeminiResponse, <-chan error) {
	ch := make(chan *mappers.GeminiResponse)
	errCh := make(chan error, 1)

	go func() {
		defer close(ch)
		defer close(errCh)

		urlPath := fmt.Sprintf("/v1/%s:streamGenerateContent?alt=sse", model)

		for attempt := 0; attempt < MaxRetries; attempt++ {
			token, _, err := c.accounts.GetRotatingToken()
			if err != nil {
				errCh <- err
				return
			}

			targetURL := c.endpoint + urlPath
			body, _ := json.Marshal(req)
			httpReq, _ := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewBuffer(body))
			httpReq.Header.Set("Content-Type", "application/json")
			httpReq.Header.Set("Authorization", "Bearer "+token)

			resp, err := c.client.Do(httpReq)
			if err != nil {
				if attempt == MaxRetries-1 {
					errCh <- err
					return
				}
				continue
			}

			if resp.StatusCode != 200 {
				resp.Body.Close()
				// Retry logic same as Unary
				if resp.StatusCode == 429 || resp.StatusCode == 403 {
					continue
				}
				if resp.StatusCode >= 500 {
					if c.endpoint == c.primaryURL {
						c.endpoint = c.fallbackURL
					}
					continue
				}
				errCh <- fmt.Errorf("upstream error %d", resp.StatusCode)
				return
			}

			// Success - Parse SSE
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "data: ") {
					data := strings.TrimPrefix(line, "data: ")
					var gResp mappers.GeminiResponse
					if err := json.Unmarshal([]byte(data), &gResp); err == nil {
						ch <- &gResp
					}
				}
			}
			resp.Body.Close()
			return // Done
		}
		errCh <- fmt.Errorf("max retries exceeded")
	}()

	return ch, errCh
}
