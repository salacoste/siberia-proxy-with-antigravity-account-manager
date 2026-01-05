package proxy

import (
	"context"
	"testing"
	"time"

	"github.com/salacoste/siberia/siberia/analytics"
	"github.com/salacoste/siberia/siberia/types"
)

func TestTelemetryChannelSaturation(t *testing.T) {
	// 1. Create Manager with small buffer
	engine := analytics.NewAnalyticsEngine()
	mgr := NewTelemetryManager(10, engine) // Buffer of 10

	// 2. Start Worker (skip wails events to avoid runtime dependency)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	mgr.Start(ctx, 1, true)

	// 3. Blast 100 events
	done := make(chan bool)
	go func() {
		for i := 0; i < 100; i++ {
			mgr.Emit(types.ProxyRequestEvent{
				URL: "http://fast.com",
			})
		}
		done <- true
	}()

	// 4. Verify it doesn't block (should finish instantly even if worker is slow)
	select {
	case <-done:
		// Success: Emit didn't block
	case <-time.After(1 * time.Second):
		t.Fatal("Emit blocked for too long! Channel saturation not handled gracefully.")
	}

	// 5. Verify stats (optional, but good to know if we dropped any)
	// In strict non-blocking mode, we might drop if buffer is full.
	// Or if `Emit` uses a buffered channel and select default, it drops.
	// Let's check implementation of Emit first.
}
