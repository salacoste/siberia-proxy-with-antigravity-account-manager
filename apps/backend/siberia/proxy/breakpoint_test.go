package proxy

import (
	"net/http/httptest"
	"testing"
	"time"
)

func TestBreakpointManager_Rules(t *testing.T) {
	bm := NewBreakpointManager()

	// 1. Add Rule
	rule := bm.AddRule("*api/v1*", "POST")
	if len(bm.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(bm.Rules))
	}
	if rule.Pattern != "*api/v1*" || rule.Method != "POST" {
		t.Errorf("Rule content mismatch: %+v", rule)
	}

	// 2. Get Rules (Verify Copy)
	rules := bm.GetRules()
	if len(rules) != 1 {
		t.Errorf("Expected 1 fetched rule, got %d", len(rules))
	}
	// Modify copy should not affect original
	rules[0].Enabled = false
	if !bm.Rules[0].Enabled { // Original should still be true (default)
		t.Errorf("GetRules() did not return a deep enough copy (or strictly a value copy of struct)")
	}
	// Wait, struct in slice is value type, so copy is value copy. safe.

	// 3. Delete Rule
	bm.DeleteRule(rule.ID)
	if len(bm.Rules) != 0 {
		t.Errorf("Expected 0 rules after delete, got %d", len(bm.Rules))
	}
}

func TestBreakpointManager_ShouldPause(t *testing.T) {
	bm := NewBreakpointManager()
	bm.AddRule("api/checkout", "POST")
	bm.AddRule("api/public", "*") // Wildcard method

	tests := []struct {
		url    string
		method string
		want   bool
	}{
		{"http://example.com/api/checkout", "POST", true},
		{"http://example.com/api/checkout", "GET", false}, // Wrong method
		{"http://example.com/api/public/list", "GET", true},
		{"http://example.com/other", "POST", false},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.url, nil)
		if got := bm.ShouldPause(req); got != tt.want {
			t.Errorf("ShouldPause(%s %s) = %v, want %v", tt.method, tt.url, got, tt.want)
		}
	}
}

func TestBreakpointManager_PendingConcurrency(t *testing.T) {
	bm := NewBreakpointManager()

	// Simulate concurrent resume requests for non-existent IDs
	go func() {
		for i := 0; i < 100; i++ {
			bm.ResumeRequest("fake-id", ModifiedRequest{})
		}
	}()

	// Simulate concurrent adds
	go func() {
		for i := 0; i < 100; i++ {
			bm.AddRule("test", "GET")
		}
	}()

	// Wait a bit
	time.Sleep(100 * time.Millisecond)

	// Check stability
	if len(bm.GetRules()) != 100 {
		t.Errorf("Expected 100 rules, got %d", len(bm.GetRules()))
	}
}
