package analytics

import (
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/salacoste/siberia/siberia/types"
)

// AnalyticsSnapshot represents the current state of traffic stats
type AnalyticsSnapshot struct {
	RPS               float64        `json:"rps"`
	BandwidthInSpeed  float64        `json:"bandwidth_in_speed"`  // Bytes/sec
	BandwidthOutSpeed float64        `json:"bandwidth_out_speed"` // Bytes/sec
	ActiveConnections int64          `json:"active_connections"`
	TopDomains        []DomainStat   `json:"top_domains"`
	ResponseCodes     map[string]int `json:"response_codes"` // "2xx", "4xx", "5xx"
	ProtocolBreakdown map[string]int `json:"protocol_breakdown"`
	TotalRequests     int64          `json:"total_requests"`
}

type DomainStat struct {
	Domain string `json:"domain"`
	Count  int64  `json:"count"`
}

// AnalyticsEngine aggregates traffic data in real-time
type AnalyticsEngine struct {
	mu sync.RWMutex

	// Counters
	totalRequests     int64
	activeConnections int64
	domainHits        map[string]int64
	responseCodes     map[string]int
	protocolStats     map[string]int

	// Rate Calculation (Sliding Window approx)
	// Rate Calculation (Sliding Window approx)
	requestBuckets      map[int64]int
	domainBuckets       map[int64]map[string]int // timestamp -> domain -> count
	bandwidthInBuckets  map[int64]int64          // timestamp -> bytes (Request)
	bandwidthOutBuckets map[int64]int64          // timestamp -> bytes (Response)

	// Configuration
	maxTopDomains int
	windowSize    int64 // seconds to keep data
}

func NewAnalyticsEngine() *AnalyticsEngine {
	return &AnalyticsEngine{
		responseCodes:       make(map[string]int),
		protocolStats:       make(map[string]int),
		requestBuckets:      make(map[int64]int),
		domainBuckets:       make(map[int64]map[string]int),
		bandwidthInBuckets:  make(map[int64]int64),
		bandwidthOutBuckets: make(map[int64]int64),
		maxTopDomains:       10,
		windowSize:          300, // Keep 5 minutes of history
	}
}

// Track processes a new proxy event
func (ae *AnalyticsEngine) Track(event types.ProxyRequestEvent) {
	ae.mu.Lock()
	defer ae.mu.Unlock()

	ae.totalRequests++

	// Sliding Window for Rates (1-second buckets)
	now := time.Now().Unix()

	// Domain Stats (Windowed)
	if u, err := url.Parse(event.URL); err == nil {
		if ae.domainBuckets[now] == nil {
			ae.domainBuckets[now] = make(map[string]int)
		}
		ae.domainBuckets[now][u.Hostname()]++
	}

	// Status Codes
	status := event.Status
	group := "Other"
	if status >= 200 && status < 300 {
		group = "2xx"
	} else if status >= 300 && status < 400 {
		group = "3xx"
	} else if status >= 400 && status < 500 {
		group = "4xx"
	} else if status >= 500 {
		group = "5xx"
	}
	ae.responseCodes[group]++

	// Sliding Window for Rates (1-second buckets)
	ae.requestBuckets[now]++

	// Estimate Request Size (In)
	// Start with rough estimation: Method + URL + Body
	reqSize := int64(len(event.Method) + len(event.URL) + len(event.ReqBody))
	// Add Headers size estimation
	for k, v := range event.ReqHeaders {
		reqSize += int64(len(k) + len(v) + 4) // +4 for ": " and "\r\n"
	}
	ae.bandwidthInBuckets[now] += reqSize

	// Protocol Stats
	proto := "HTTP/1.1" // Default assumption if not present (goproxy usually handles 1.1)
	if event.ReqHeaders["Upgrade"] == "websocket" {
		proto = "WebSocket"
	}
	// goproxy doesn't easily expose the actual wire protocol (HTTP/2 vs 1.1) in the DoFunc context unless we dig deeper.
	// But let's check if we can infer or if we simply rely on 1.1 for now.
	// Ideally 'event' should have a 'Protocol' field.
	// For now, let's just track WebSocket vs HTTP.
	// FIXME: Real HTTP/2 detection requires TLS connection state inspection which is hard in this layer.
	ae.protocolStats[proto]++

	// Response Size (Out) - event.Size is typically the response content length + headers
	ae.bandwidthOutBuckets[now] += event.Size

	// Cleanup old buckets (> 5 seconds ago for Rate, > windowSize for Domains)
	if len(ae.requestBuckets) > 10 {
		for k := range ae.requestBuckets {
			if k < now-5 {
				delete(ae.requestBuckets, k)
				delete(ae.bandwidthInBuckets, k)
				delete(ae.bandwidthOutBuckets, k)
			}
		}
	}
	// Cleanup Domain Buckets (check less frequently or just separate loop)
	if len(ae.domainBuckets) > 0 { // Simple check
		cutoff := now - ae.windowSize
		for k := range ae.domainBuckets {
			if k < cutoff {
				delete(ae.domainBuckets, k)
			}
		}
	}
}

// TrackActive increments active connection count
func (ae *AnalyticsEngine) IncrementActive() {
	ae.mu.Lock()
	ae.activeConnections++
	ae.mu.Unlock()
}

// TrackInactive decrements active connection count
func (ae *AnalyticsEngine) DecrementActive() {
	ae.mu.Lock()
	if ae.activeConnections > 0 {
		ae.activeConnections--
	}
	ae.mu.Unlock()
}

// GetSnapshot returns the current calculated stats
func (ae *AnalyticsEngine) GetSnapshot() AnalyticsSnapshot {
	ae.mu.RLock()
	defer ae.mu.RUnlock()
	// ... existing snapshot logic ...

	// Calculate RPS (Average of last 5 seconds)
	now := time.Now().Unix()
	var totalReqsInWindow int
	var totalBytesInWindow int64
	var totalBytesOutWindow int64

	for i := 0; i < 5; i++ {
		ts := now - int64(i)
		if count, ok := ae.requestBuckets[ts]; ok {
			totalReqsInWindow += count
		}
		if bytes, ok := ae.bandwidthInBuckets[ts]; ok {
			totalBytesInWindow += bytes
		}
		if bytes, ok := ae.bandwidthOutBuckets[ts]; ok {
			totalBytesOutWindow += bytes
		}
	}

	// Top Domains (Windowed Aggregation)
	domainCounts := make(map[string]int64)
	domainCutoff := now - ae.windowSize
	for ts, bucket := range ae.domainBuckets {
		if ts >= domainCutoff {
			for domain, count := range bucket {
				domainCounts[domain] += int64(count)
			}
		}
	}

	topDomains := make([]DomainStat, 0)
	for d, c := range domainCounts {
		topDomains = append(topDomains, DomainStat{Domain: d, Count: c})
	}
	// Sort by count descending
	sort.Slice(topDomains, func(i, j int) bool {
		return topDomains[i].Count > topDomains[j].Count
	})

	if len(topDomains) > ae.maxTopDomains {
		topDomains = topDomains[:ae.maxTopDomains]
	}

	return AnalyticsSnapshot{
		RPS:               float64(totalReqsInWindow) / 5.0,
		BandwidthInSpeed:  float64(totalBytesInWindow) / 5.0,
		BandwidthOutSpeed: float64(totalBytesOutWindow) / 5.0,
		ActiveConnections: ae.activeConnections,
		ResponseCodes:     ae.deepCopyMap(ae.responseCodes),
		ProtocolBreakdown: ae.deepCopyMap(ae.protocolStats), // Include broken down map
		TopDomains:        topDomains,
		TotalRequests:     ae.totalRequests,
	}
}

func (ae *AnalyticsEngine) deepCopyMap(m map[string]int) map[string]int {
	c := make(map[string]int)
	for k, v := range m {
		c[k] = v
	}
	return c
}
