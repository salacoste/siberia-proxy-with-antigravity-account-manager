package proxy

import (
	"context"
	"sync"
	"time"

	"github.com/salacoste/siberia/siberia/analytics"
	"github.com/salacoste/siberia/siberia/logger"
	"github.com/salacoste/siberia/siberia/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TelemetryManager handles async logging and event emission
type TelemetryManager struct {
	eventsChan    chan types.ProxyRequestEvent
	ctx           context.Context
	skipWails     bool
	wg            sync.WaitGroup
	bufferSize    int
	droppedEvents int64
	mu            sync.Mutex
	analytics     *analytics.AnalyticsEngine
}

// NewTelemetryManager creates a new manager with specified buffer size
func NewTelemetryManager(bufferSize int, analyticsEngine *analytics.AnalyticsEngine) *TelemetryManager {
	return &TelemetryManager{
		eventsChan: make(chan types.ProxyRequestEvent, bufferSize),
		bufferSize: bufferSize,
		analytics:  analyticsEngine,
	}
}

// Start spawns the worker goroutines
func (tm *TelemetryManager) Start(ctx context.Context, workers int, skipWails bool) {
	tm.ctx = ctx
	tm.skipWails = skipWails

	for i := 0; i < workers; i++ {
		tm.wg.Add(1)
		go tm.worker()
	}
}

// Stop waits for workers to drain (with timeout)
func (tm *TelemetryManager) Stop() {
	close(tm.eventsChan)
	// Wait logic could go here, or we rely on main app context cancellation
	tm.wg.Wait()
}

// Emit queues an event. If buffer is full, it drops the event to preserve performance.
func (tm *TelemetryManager) Emit(event types.ProxyRequestEvent) {
	select {
	case tm.eventsChan <- event:
		// Queued successfully
	default:
		// Buffer full, drop event
		tm.incrementDropped()
	}
}

func (tm *TelemetryManager) worker() {
	defer tm.wg.Done()

	// Batching disk writes could be an optimization, but for now we just process async
	for event := range tm.eventsChan {
		// 1. Emit to Frontend (Wails)
		if !tm.skipWails {
			runtime.EventsEmit(tm.ctx, "proxy:log", event)
		}

		// 2. Log to Disk (Simulated "Heavy I/O")
		// In a real high-throughput system, we might want to batch these writes.
		// For Story-35, just moving it off the request goroutine is the win.
		logger.LogAccess(logger.AccessEntry{
			Time:       event.Time,
			Timestamp:  time.Now().Unix(),
			Method:     event.Method,
			URL:        event.URL,
			Status:     event.Status,
			DurationMs: event.Duration,
			Size:       event.Size,
		})

		// 3. Update Analytics
		if tm.analytics != nil {
			tm.analytics.Track(event)
		}
	}
}

func (tm *TelemetryManager) incrementDropped() {
	tm.mu.Lock()
	tm.droppedEvents++
	tm.mu.Unlock()
}

func (tm *TelemetryManager) GetDroppedCount() int64 {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.droppedEvents
}
