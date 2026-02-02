package middleware

import (
	"net/http"
	"sync"
	"time"
)

// rateLimiter implements token bucket algorithm
// WHY: Prevent abuse, protect against DoS attacks
type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int // Requests per minute
	burst    int // Burst capacity
}

// visitor represents a single IP's rate limit state
type visitor struct {
	tokens     int
	lastRefill time.Time
}

// newRateLimiter creates a new rate limiter
func newRateLimiter(rate, burst int) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}

	// Cleanup old visitors every 5 minutes
	go rl.cleanup()

	return rl
}

// cleanup removes old visitors
func (rl *rateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastRefill) > 10*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// allow checks if request is allowed
func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		// New visitor
		rl.visitors[ip] = &visitor{
			tokens:     rl.burst - 1,
			lastRefill: time.Now(),
		}
		return true
	}

	// Refill tokens based on time passed
	now := time.Now()
	elapsed := now.Sub(v.lastRefill)
	tokensToAdd := int(elapsed.Minutes()) * rl.rate

	if tokensToAdd > 0 {
		v.tokens = min(v.tokens+tokensToAdd, rl.burst)
		v.lastRefill = now
	}

	// Check if tokens available
	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

// min returns minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RateLimit creates rate limiting middleware
// rate: requests per minute
// burst: burst capacity
func RateLimit(rate, burst int) func(http.Handler) http.Handler {
	limiter := newRateLimiter(rate, burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get client IP
			ip := r.RemoteAddr

			// Check rate limit
			if !limiter.allow(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error":"rate limit exceeded","message":"too many requests","code":429}`))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
