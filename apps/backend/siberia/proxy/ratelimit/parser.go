package ratelimit

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	// "Try again in 2m 30s" or "2m 30s"
	reDurationNLP = regexp.MustCompile(`(?i)(?:try again in|retry in|after)\s*((?:\d+\s*[hms]\s*)+)`)
	// "quotaResetDelay": "42s" (JSON fields often used)
	reJsonDelay = regexp.MustCompile(`(?i)"?(?:quotaResetDelay|retryAfter)"?\s*[:=]\s*"?(\d+[hms])"?`)
	// "retry after 50" (seconds implicitly)
	// Go regexp doesn't support lookaround, so we match basic form and rely on priority.
	reStandardSeconds = regexp.MustCompile(`(?i)(?:retry after|wait)\s*(\d+)(?:\s|$)`)
)

// DefaultBackoff is used when parsing fails but we know it's a rate limit
const DefaultBackoff = 30 * time.Second
const DefaultQuotaBackoff = 1 * time.Hour

// ParseError analyzes the error body and returns a RateLimitError
// valid even if no duration is found (returns defaults)
func ParseError(body string) *RateLimitError {
	body = strings.ToLower(body)
	isQuota := strings.Contains(body, "quota") || strings.Contains(body, "exhausted")

	retryAfter := parseDuration(body)
	if retryAfter == 0 {
		if isQuota {
			retryAfter = DefaultQuotaBackoff
		} else {
			retryAfter = DefaultBackoff
		}
	}

	return &RateLimitError{
		RetryAfter: retryAfter,
		IsQuota:    isQuota,
		Original:   truncate(body, 50),
	}
}

func parseDuration(s string) time.Duration {
	// 1. NLP "Try again in 2m 30s"
	if matches := reDurationNLP.FindStringSubmatch(s); len(matches) > 1 {
		dur, err := time.ParseDuration(strings.ReplaceAll(matches[1], " ", ""))
		if err == nil {
			return dur
		}
	}

	// 2. JSON/Field "42s"
	if matches := reJsonDelay.FindStringSubmatch(s); len(matches) > 1 {
		dur, err := time.ParseDuration(matches[1])
		if err == nil {
			return dur
		}
	}

	// 3. Raw Seconds "retry after 60"
	if matches := reStandardSeconds.FindStringSubmatch(s); len(matches) > 1 {
		if sec, err := strconv.Atoi(matches[1]); err == nil {
			return time.Duration(sec) * time.Second
		}
	}

	return 0
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
