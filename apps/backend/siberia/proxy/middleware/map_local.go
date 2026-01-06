package middleware

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/elazarl/goproxy"
)

// MapLocalRule defines a rule for mapping a URL pattern to a local file
type MapLocalRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	UrlRegex    string `json:"url_regex"`
	LocalPath   string `json:"local_path"`
	ContentType string `json:"content_type"` // Optional override
	Status      int    `json:"status"`       // HTTP Status code (default 200)
}

// MapLocalMiddleware handles serving local files for matched requests
type MapLocalMiddleware struct {
	rules sync.Map // map[string]*MapLocalRule (key is ID)
}

func NewMapLocalMiddleware() *MapLocalMiddleware {
	return &MapLocalMiddleware{}
}

// AddRule adds or updates a rule
func (m *MapLocalMiddleware) AddRule(rule MapLocalRule) error {
	// Validate Regex
	if _, err := regexp.Compile(rule.UrlRegex); err != nil {
		return fmt.Errorf("invalid regex: %w", err)
	}
	// Validate Path (basic check)
	if rule.LocalPath == "" {
		return fmt.Errorf("local path cannot be empty")
	}
	m.rules.Store(rule.ID, &rule)
	return nil
}

// RemoveRule deletes a rule by ID
func (m *MapLocalMiddleware) RemoveRule(id string) {
	m.rules.Delete(id)
}

// RangeRules iterates over all rules
func (m *MapLocalMiddleware) RangeRules(f func(rule MapLocalRule)) {
	m.rules.Range(func(key, value interface{}) bool {
		if rule, ok := value.(*MapLocalRule); ok {
			f(*rule)
		}
		return true
	})
}

// ClearRules removes all rules

func (m *MapLocalMiddleware) ClearRules() {
	m.rules.Range(func(key, value interface{}) bool {
		m.rules.Delete(key)
		return true
	})
}

// HandleRequest is the goproxy middleware function
func (m *MapLocalMiddleware) HandleRequest(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// 1. Iterate over rules (Snapshot for safety)
	// optimization: we could cache compiled regexes, but for <100 rules iteration is fast enough
	var matchedRule *MapLocalRule

	m.rules.Range(func(key, value interface{}) bool {
		rule := value.(*MapLocalRule)
		if !rule.Enabled {
			return true
		}

		matched, err := regexp.MatchString(rule.UrlRegex, req.URL.String())
		if err == nil && matched {
			matchedRule = rule
			return false // Stop iteration
		}
		return true
	})

	// 2. If no match, pass through
	if matchedRule == nil {
		return req, nil
	}

	// 3. Serve Local File
	// Check if file exists
	info, err := os.Stat(matchedRule.LocalPath)
	if os.IsNotExist(err) || info.IsDir() {
		return req, goproxy.NewResponse(req,
			goproxy.ContentTypeText, http.StatusNotFound,
			fmt.Sprintf("Siberia Map Local: File not found: %s", matchedRule.LocalPath))
	}

	// Read file
	file, err := os.Open(matchedRule.LocalPath)
	if err != nil {
		return req, goproxy.NewResponse(req,
			goproxy.ContentTypeText, http.StatusInternalServerError,
			fmt.Sprintf("Siberia Map Local: Read error: %v", err))
	}
	defer file.Close()

	// Determine Content-Type
	contentType := matchedRule.ContentType
	if contentType == "" {
		// Detect from file extension/content
		contentType = "application/octet-stream"
		// Simple extension check
		ext := filepath.Ext(matchedRule.LocalPath)
		switch ext {
		case ".json":
			contentType = "application/json"
		case ".js":
			contentType = "application/javascript"
		case ".css":
			contentType = "text/css"
		case ".html":
			contentType = "text/html"
		case ".txt":
			contentType = "text/plain"
		case ".png":
			contentType = "image/png"
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		}
	}

	// Determine Status Code
	statusCode := matchedRule.Status
	if statusCode == 0 {
		statusCode = http.StatusOK
	}

	// Create Response
	// resp := goproxy.NewResponse(req, contentType, statusCode, "")

	// We need to write the body manually because goproxy.NewResponse takes a string/byte slice
	// but we want to potentially stream or just read valid bytes.
	// Actually goproxy.NewResponse sets the body. Let's read file fully.
	// For huge files this is bad, but for mocks it's fine.
	content, _ := os.ReadFile(matchedRule.LocalPath)

	// Create a new response with the content
	resp := goproxy.NewResponse(req, contentType, statusCode, string(content))

	// Add a header to indicate mapping happened
	resp.Header.Set("X-Siberia-Map-Local", matchedRule.ID)
	resp.Header.Set("Date", time.Now().Format(time.RFC1123))

	fmt.Printf("[MapLocal] Serving %s -> %s\n", req.URL.String(), matchedRule.LocalPath)

	return req, resp
}
