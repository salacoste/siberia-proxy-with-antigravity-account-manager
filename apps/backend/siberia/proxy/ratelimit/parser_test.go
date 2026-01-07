package ratelimit

import (
	"testing"
	"time"
)

func TestParseError(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantDuration time.Duration
		wantQuota    bool
		wantDefault  bool // if true, accepts DefaultBackoff
	}{
		{
			name:         "NLP Minutes Seconds",
			body:         "Error: rate limit exceeded, try again in 2m 30s.",
			wantDuration: 2*time.Minute + 30*time.Second,
			wantQuota:    false,
		},
		{
			name:         "JSON Quota Delay",
			body:         `{"error": "Resource exhausted", "quotaResetDelay": "45s"}`,
			wantDuration: 45 * time.Second,
			wantQuota:    true,
		},
		{
			name:         "Raw Seconds",
			body:         "Retry After 60",
			wantDuration: 60 * time.Second,
			wantQuota:    false,
		},
		{
			name:         "Quota Keyword Fallback",
			body:         "The quota has been exhausted for this project.",
			wantQuota:    true,
			wantDuration: DefaultQuotaBackoff,
		},
		{
			name:         "Unknown Error Fallback",
			body:         "Too many requests",
			wantQuota:    false,
			wantDuration: DefaultBackoff,
		},
		{
			name:         "Complex NLP Mixed",
			body:         "Rate limit reached. Please retry in 5m.",
			wantDuration: 5 * time.Minute,
			wantQuota:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseError(tt.body)

			if got.IsQuota != tt.wantQuota {
				t.Errorf("IsQuota = %v, want %v", got.IsQuota, tt.wantQuota)
			}

			if got.RetryAfter != tt.wantDuration {
				t.Errorf("RetryAfter = %v, want %v", got.RetryAfter, tt.wantDuration)
			}
		})
	}
}
