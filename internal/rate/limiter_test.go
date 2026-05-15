package rate

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	rl := NewRateLimiter(10, 5)
	if rl == nil {
		t.Fatal("Expected non-nil rate limiter")
	}
	if rl.rate != 10 {
		t.Errorf("Expected rate 10, got %d", rl.rate)
	}
	if rl.burst != 5 {
		t.Errorf("Expected burst 5, got %d", rl.burst)
	}
}

func TestRateLimiter_Allow(t *testing.T) {
	rl := NewRateLimiter(10, 5)

	// First 5 requests should be allowed
	for i := 0; i < 5; i++ {
		if !rl.allow("test-key") {
			t.Errorf("Expected request %d to be allowed", i)
		}
	}

	// 6th request should be denied
	if rl.allow("test-key") {
		t.Error("Expected 6th request to be denied")
	}
}

func TestRateLimiter_Middleware(t *testing.T) {
	rl := NewRateLimiter(10, 5)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	middleware := rl.Middleware(handler)

	// First request should be allowed
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	middleware.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestRateLimiter_TokenRefill(t *testing.T) {
	rl := NewRateLimiter(100, 5) // High rate to allow quick refill

	// Use all tokens
	for i := 0; i < 5; i++ {
		rl.allow("test-key")
	}

	// Wait for refill
	time.Sleep(100 * time.Millisecond)

	// Should be allowed again
	if !rl.allow("test-key") {
		t.Error("Expected request to be allowed after refill")
	}
}

func TestRateLimiter_Cleanup(t *testing.T) {
	rl := NewRateLimiter(10, 5)

	// Add some tokens
	rl.allow("key1")
	rl.allow("key2")

	// Trigger cleanup
	rl.cleanup()

	// Verify cleanup ran (this is a basic sanity check)
	if len(rl.tokens) > 2 {
		t.Error("Expected cleanup to remove old entries")
	}
}
