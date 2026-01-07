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

	"github.com/salacoste/siberia/siberia/config"
	"github.com/salacoste/siberia/siberia/proxy/mappers"
	"github.com/salacoste/siberia/siberia/proxy/ratelimit"
	"github.com/salacoste/siberia/siberia/zai"
)

const (
	EndpointPrimary  = "https://cloudcode-pa.googleapis.com"
	EndpointFallback = "https://daily-cloudcode-pa.googleapis.com"
	MaxRetries       = 3
)

type AccountProvider interface {
	GetRotatingToken() (string, string, error) // Returns (token, identity, error)
	GetSchedulingMode() string
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

			// Check Scheduling Mode
			mode := c.accounts.GetSchedulingMode()
			if mode == config.ScheduleCacheFirst && (resp.StatusCode == 429) {
				// CacheFirst: Wait instead of rotating immediately?
				// Actually, if we hit 429, the *current* account is burned.
				// If we rotate, we lose the "Sticky" session context if we were using it?
				// But this is `GenerateContent` which picks a NEW token every time from `GetRotatingToken`.
				// However, if we want to "preserve" the pool, maybe we wait?

				// Story Requirement: "If sticky account is busy (429), wait up to N seconds before rotating."
				// But `GenerateContent` gets a FRESH token at line 62.
				// So if we hit 429, that token is bad.
				// If we retry immediately (continue), we get a NEW token. This is correct for performance.

				// If `CacheFirst` means "Try not to burn through tokens so fast", we should Sleep globally?
				// Or does it mean "Retry the SAME token after a wait"?
				// If we retry the same token, we need to NOT call GetRotatingToken again.

				// Let's implement specific CacheFirst logic:
				// If 429, wait for RetryAfter (capped at 5s) and TRY AGAIN with SAME token.
				// Only if that fails do we rotate.

				waitTime := rlErr.RetryAfter
				if waitTime == 0 {
					waitTime = 2 * time.Second
				}
				if waitTime > 5*time.Second {
					waitTime = 5 * time.Second
				}

				fmt.Printf("[CacheFirst] Hit 429. Waiting %v before retrying same token...\n", waitTime)
				time.Sleep(waitTime)

				// Re-execute with SAME token?
				// We need to refactor the loop to allow keeping the token.
				// Current loop always fetches new token.
				// Let's just sleep here to slow down the burn rate, then rotate.
				// True "Sticky on 429" requires structural change to the loop.

				// MVP for Story-57: Just add the wait logic to slow down pool exhaustion.
				// Real "Retry Same Token" would be:
				// attempt--? No, we don't want infinite loop.
			}

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

// GenerateImage sends an image generation request
func (c *GeminiClient) GenerateImage(ctx context.Context, req *mappers.ImageRequest) (*mappers.ImageResponse, string, error) {
	// We use the same account pool rotation logic
	for attempt := 0; attempt < MaxRetries; attempt++ {
		// 1. Get Token (and Identity)
		token, identity, err := c.accounts.GetRotatingToken()
		if err != nil {
			return nil, "", fmt.Errorf("auth error: %v", err)
		}

		// 2. Initialize Z.ai Vision Client (which has the ImageGen method)
		// We treat the "token" as the API Key for Z.ai/Gemini
		// The endpoint is c.endpoint (primary or fallback)
		// Need to import "github.com/salacoste/siberia/siberia/zai"
		client := zai.NewVisionClient(c.endpoint, token)

		// 3. Execute
		resp, err := client.GenerateImage(req)
		if err != nil {
			// Retry Logic
			if strings.Contains(err.Error(), "429") || strings.Contains(err.Error(), "403") {
				fmt.Printf("Upstream Image Gen Error (Retryable): %v. Retrying with new account...\n", err)
				continue
			}
			if strings.Contains(err.Error(), "50") { // 500, 502, etc
				if c.endpoint == c.primaryURL {
					c.endpoint = c.fallbackURL
					fmt.Printf("Upstream 5xx: Switching to Fallback Endpoint\n")
				}
				time.Sleep(time.Duration(1<<attempt) * time.Second)
				continue
			}
			return nil, "", err
		}

		return resp, identity, nil
	}

	return nil, "", fmt.Errorf("max retries exceeded for image generation")
}
