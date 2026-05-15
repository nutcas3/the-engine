package rate

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	mu           sync.Mutex
	tokens       map[string]*tokenBucket
	rate         int // tokens per second
	burst        int // max tokens
	cleanupTimer *time.Timer
}

// tokenBucket represents a token bucket for a specific key
type tokenBucket struct {
	tokens     float64
	lastUpdate time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(map[string]*tokenBucket),
		rate:   rate,
		burst:  burst,
	}

	// Start cleanup routine
	rl.startCleanup()

	return rl
}

// Middleware returns a rate limiting middleware
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr // In production, use user ID or API key

		if !rl.allow(key) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// allow checks if a request is allowed for the given key
func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.tokens[key]

	if !exists {
		bucket = &tokenBucket{
			tokens:     float64(rl.burst),
			lastUpdate: now,
		}
		rl.tokens[key] = bucket
	}

	// Add tokens based on elapsed time
	elapsed := now.Sub(bucket.lastUpdate).Seconds()
	bucket.tokens += elapsed * float64(rl.rate)

	// Cap at burst size
	if bucket.tokens > float64(rl.burst) {
		bucket.tokens = float64(rl.burst)
	}

	bucket.lastUpdate = now

	// Check if we have enough tokens
	if bucket.tokens >= 1 {
		bucket.tokens -= 1
		return true
	}

	return false
}

// startCleanup periodically removes unused buckets
func (rl *RateLimiter) startCleanup() {
	rl.cleanupTimer = time.AfterFunc(5*time.Minute, func() {
		rl.cleanup()
		rl.startCleanup()
	})
}

// cleanup removes buckets that haven't been used recently
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, bucket := range rl.tokens {
		if now.Sub(bucket.lastUpdate) > 10*time.Minute {
			delete(rl.tokens, key)
		}
	}
}
