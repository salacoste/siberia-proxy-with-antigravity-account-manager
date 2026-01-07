package ratelimit

import (
	"fmt"
	"time"
)

// RateLimitError represents a parsed rate limit error
type RateLimitError struct {
	RetryAfter time.Duration
	IsQuota    bool // True if long-term quota exhaustion, False if short-term rate limit
	Original   string
}

func (e *RateLimitError) Error() string {
	msg := "rate limit exceeded"
	if e.IsQuota {
		msg = "quota exhausted"
	}
	return fmt.Sprintf("%s: retry after %v (original: %s)", msg, e.RetryAfter, e.Original)
}
