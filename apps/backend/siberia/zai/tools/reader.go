package tools

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
)

// CallWebReader fetches a URL and converts it to Markdown.
// It automatically strips tracking parameters.
func CallWebReader(targetURL string) (string, error) {
	// 1. Normalize URL (Strip Tracking)
	cleanURL, err := normalizeURL(targetURL)
	if err != nil {
		return "", fmt.Errorf("invalid url: %w", err)
	}

	// 2. Fetch HTML
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(cleanURL)
	if err != nil {
		return "", fmt.Errorf("fetch error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch url: status %d", resp.StatusCode)
	}

	// 3. Convert to Markdown
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("conversion error: %w", err)
	}

	return markdown.String(), nil
}

// normalizeURL removes tracking parameters like utm_, gclid, etc.
func normalizeURL(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for key := range q {
		lowerKey := strings.ToLower(key)
		if strings.HasPrefix(lowerKey, "utm_") ||
			strings.HasPrefix(lowerKey, "gclid") ||
			strings.HasPrefix(lowerKey, "fbclid") ||
			strings.HasPrefix(lowerKey, "yclid") ||
			strings.HasPrefix(lowerKey, "mc_eid") {
			q.Del(key)
		}
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}
