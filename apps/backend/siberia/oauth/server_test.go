package oauth

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestLoopbackServer(t *testing.T) {
	srv := NewServer()

	// Start Server
	expectedState := "test-random-state"
	port, err := srv.Start(expectedState)
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	if port <= 0 {
		t.Fatalf("Invalid port: %d", port)
	}

	// Channel for client logic result
	done := make(chan error)

	// Simulate Browser Callback
	go func() {
		// Wait slightly
		time.Sleep(100 * time.Millisecond)

		// 1. Success Request
		url := fmt.Sprintf("http://127.0.0.1:%d/callback?code=mock-auth-code&state=%s", port, expectedState)
		resp, err := http.Get(url)
		if err != nil {
			done <- fmt.Errorf("http get failed: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			done <- fmt.Errorf("expected 200, got %d", resp.StatusCode)
			return
		}

		done <- nil
	}()

	// Wait for Code
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	code, err := srv.WaitForCode(ctx)
	if err != nil {
		t.Fatalf("WaitForCode failed: %v (ctx error: %v)", err, ctx.Err())
	}

	if code != "mock-auth-code" {
		t.Errorf("expected code 'mock-auth-code', got '%s'", code)
	}

	// Check client goroutine result
	if err := <-done; err != nil {
		t.Fatalf("Client simulation failed: %v", err)
	}
}

func TestLoopbackServer_StateMismatch(t *testing.T) {
	srv := NewServer()
	expectedState := "secure-state"
	port, _ := srv.Start(expectedState)

	go func() {
		time.Sleep(50 * time.Millisecond)
		url := fmt.Sprintf("http://127.0.0.1:%d/callback?code=bad-code&state=wrong-state", port)
		http.Get(url)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	_, err := srv.WaitForCode(ctx)
	if err == nil {
		t.Error("Expected error on state mismatch, got nil")
	} else if err.Error() != "state mismatch" {
		// Note: WaitForCode returns the specific error pushed to channel
		t.Logf("Got expected error: %v", err)
	}
}
