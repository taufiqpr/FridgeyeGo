package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type tokenBucket struct {
	capacity    int
	tokens      int
	refillEvery time.Duration
	lastRefill  time.Time
}

type limiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
	cap     int
	rate    time.Duration
}

func newLimiter(capacity int, refillEvery time.Duration) *limiter {
	return &limiter{
		buckets: make(map[string]*tokenBucket),
		cap:     capacity,
		rate:    refillEvery,
	}
}

func (l *limiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	b, ok := l.buckets[key]
	now := time.Now()
	if !ok {
		b = &tokenBucket{capacity: l.cap, tokens: l.cap, refillEvery: l.rate, lastRefill: now}
		l.buckets[key] = b
	}
	// refill
	elapsed := now.Sub(b.lastRefill)
	if elapsed >= b.refillEvery {
		n := int(elapsed / b.refillEvery)
		b.tokens += n
		if b.tokens > b.capacity {
			b.tokens = b.capacity
		}
		b.lastRefill = now
	}
	if b.tokens <= 0 {
		return false
	}
	b.tokens--
	return true
}

func clientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func RateLimit(reqPerWindow int, window time.Duration) func(http.Handler) http.Handler {
	l := newLimiter(reqPerWindow, window)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := clientIP(r)
			if !l.allow(key) {
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
