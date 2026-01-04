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
	// GetSnapshot averages over 5 seconds.
	// Out size = 1024 bytes
	// In size = Method(3) + URL(42) + headers(0) + body(0) = 45 bytes
	// specific expected values:
	// InSpeed = 45 / 5.0 = 9.0
	// OutSpeed = 1024 / 5.0 = 204.8

	if snap.BandwidthOutSpeed != 204.8 {
		t.Errorf("Expected 204.8 BandwidthOut, got %f", snap.BandwidthOutSpeed)
	}

	if snap.BandwidthInSpeed != 9.0 {
		t.Errorf("Expected 9.0 BandwidthIn, got %f", snap.BandwidthInSpeed)
	}

	// 4. Check Protocols
	if snap.ProtocolBreakdown["HTTP/1.1"] != 1 {
		t.Errorf("Expected 1 HTTP/1.1 request, got %v", snap.ProtocolBreakdown)
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

func TestAnalyticsEngine_ActiveConnections(t *testing.T) {
	engine := NewAnalyticsEngine()

	if engine.GetSnapshot().ActiveConnections != 0 {
		t.Errorf("Expected 0 active connections initially")
	}

	engine.IncrementActive()
	engine.IncrementActive()

	if val := engine.GetSnapshot().ActiveConnections; val != 2 {
		t.Errorf("Expected 2 active connections, got %d", val)
	}

	engine.DecrementActive()
	if val := engine.GetSnapshot().ActiveConnections; val != 1 {
		t.Errorf("Expected 1 active connection, got %d", val)
	}
}
