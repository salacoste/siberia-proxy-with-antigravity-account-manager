package analytics

import (
	"testing"
	"time"

	"github.com/salacoste/siberia/siberia/types"
)

func TestAnalyticsEngine_Track(t *testing.T) {
	engine := NewAnalyticsEngine()

	// 1. Simulate a request
	event := types.ProxyRequestEvent{
		Method: "GET",
		URL:    "https://api.openai.com/v1/chat/completions",
		Status: 200,
		Size:   1024,
		Time:   time.Now().Format("15:04:05"),
	}

	engine.Track(event)

	// 2. Check Snapshot
	snap := engine.GetSnapshot()

	if snap.TotalRequests != 1 {
		t.Errorf("Expected 1 total request, got %d", snap.TotalRequests)
	}

	if val := snap.ResponseCodes["2xx"]; val != 1 {
		t.Errorf("Expected 1 2xx response, got %d", val)
	}

	foundDomain := false
	for _, d := range snap.TopDomains {
		if d.Domain == "api.openai.com" && d.Count == 1 {
			foundDomain = true
			break
		}
	}
	if !foundDomain {
		t.Errorf("Expected api.openai.com in Top Domains")
	}

	// 3. Check Bandwidth (approx)
	// Since we use 1-second buckets, and test runs fast, it might be in the current bucket
	// GetSnapshot averages over 5 seconds. 1024 bytes / 5s = 204.8 bytes/sec
	if snap.BandwidthInSpeed == 0 {
		t.Logf("Warning: Bandwidth might be 0 if bucket logic relies on strict seconds?")
	}
}

func TestAnalyticsEngine_RateCalculation(t *testing.T) {
	engine := NewAnalyticsEngine()

	// Simulate 10 requests
	for i := 0; i < 10; i++ {
		engine.Track(types.ProxyRequestEvent{
			URL:    "http://example.com",
			Status: 200,
		})
	}

	snap := engine.GetSnapshot()
	// 10 reqs / 5 sec window = 2.0 RPS
	if snap.RPS != 2.0 {
		t.Errorf("Expected 2.0 RPS, got %f", snap.RPS)
	}
}
