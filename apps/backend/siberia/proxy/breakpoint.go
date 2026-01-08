package proxy

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// BreakpointRule defines a condition to pause a request
type BreakpointRule struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
	Pattern string `json:"pattern"` // Simple substring match on URL for now
	Method  string `json:"method"`  // Optional method filter
}

// PendingRequest represents a paused request waiting for user action
type PendingRequest struct {
	ID        string            `json:"id"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Body      string            `json:"body"`
	CreatedAt time.Time         `json:"created_at"`

	// resumeChan is used to signal the blocking goroutine to proceed
	resumeChan chan ModifiedRequest
}

// ModifiedRequest holds the changes applied by the user (or empty if direct resume)
type ModifiedRequest struct {
	Drop    bool              `json:"drop"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type BreakpointManager struct {
	ctx     context.Context
	mu      sync.RWMutex
	Rules   []BreakpointRule          `json:"rules"`
	Pending map[string]PendingRequest `json:"pending"`
}

func NewBreakpointManager() *BreakpointManager {
	return &BreakpointManager{
		Pending: make(map[string]PendingRequest),
		Rules:   []BreakpointRule{},
	}
}

func (bm *BreakpointManager) SetContext(ctx context.Context) {
	bm.ctx = ctx
}

// AddRule adds a new breakpoint rule
func (bm *BreakpointManager) AddRule(pattern, method string) BreakpointRule {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	rule := BreakpointRule{
		ID:      uuid.New().String(),
		Enabled: true,
		Pattern: pattern,
		Method:  method,
	}
	bm.Rules = append(bm.Rules, rule)
	return rule
}

// DeleteRule removes a rule
func (bm *BreakpointManager) DeleteRule(id string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	newRules := []BreakpointRule{}
	for _, r := range bm.Rules {
		if r.ID != id {
			newRules = append(newRules, r)
		}
	}
	bm.Rules = newRules
}

// ShouldPause checks if a request matches any active rule
// Returns true only if a rule matches
func (bm *BreakpointManager) ShouldPause(req *http.Request) bool {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	for _, rule := range bm.Rules {
		if !rule.Enabled {
			continue
		}

		// Method check (if specified)
		if rule.Method != "" && rule.Method != "*" && !strings.EqualFold(req.Method, rule.Method) {
			continue
		}

		// Pattern check (URL substring)
		if strings.Contains(req.URL.String(), rule.Pattern) {
			return true
		}
	}
	return false
}

// HasActiveBreakpoints checks if there are any enabled rules
func (bm *BreakpointManager) HasActiveBreakpoints() bool {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	for _, rule := range bm.Rules {
		if rule.Enabled {
			return true
		}
	}
	return false
}

// GetRules returns a copy of current rules safely
func (bm *BreakpointManager) GetRules() []BreakpointRule {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	// Return a copy to avoid race conditions if caller modifies slice
	rules := make([]BreakpointRule, len(bm.Rules))
	copy(rules, bm.Rules)
	return rules
}

// PauseRequest blocks execution until Resumed.
// It broadcasts the event to Frontend.
func (bm *BreakpointManager) PauseRequest(req *http.Request, reqBody string) (ModifiedRequest, bool) {
	id := uuid.New().String()
	ch := make(chan ModifiedRequest)

	headers := make(map[string]string)
	for k, v := range req.Header {
		headers[k] = strings.Join(v, ", ")
	}

	pending := PendingRequest{
		ID:         id,
		Method:     req.Method,
		URL:        req.URL.String(),
		Headers:    headers,
		Body:       reqBody,
		CreatedAt:  time.Now(),
		resumeChan: ch,
	}

	bm.mu.Lock()
	bm.Pending[id] = pending
	bm.mu.Unlock()

	// Notify UI
	if bm.ctx != nil {
		runtime.EventsEmit(bm.ctx, "breakpoint:hit", pending)
	}

	// Block until Resume called
	select {
	case mod := <-ch:
		return mod, true
	case <-time.After(300 * time.Second): // 5 min timeout safety
		return ModifiedRequest{Drop: true}, false
	case <-req.Context().Done(): // Client gave up
		return ModifiedRequest{Drop: true}, false
	}
}

// ResumeRequest is called from Wails (Frontend) to release the block
func (bm *BreakpointManager) ResumeRequest(id string, mod ModifiedRequest) bool {
	bm.mu.Lock()
	pending, exists := bm.Pending[id]
	if exists {
		delete(bm.Pending, id)
	}
	bm.mu.Unlock()

	if !exists {
		return false
	}

	// Send modification to the waiting channel
	select {
	case pending.resumeChan <- mod:
		return true
	default:
		return false
	}
}
