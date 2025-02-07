package middleware

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// getLimiter returns a limiter for a given IP address.
// If the limiter does not exist - it creates a new one.
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if limiter, exists := limiters[ip]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(1, 5)
	limiters[ip] = limiter

	go func(ip string) {
		time.Sleep(10 * time.Minute)
		mu.Lock()
		delete(limiters, ip)
		mu.Unlock()
	}(ip)

	return limiter
}

// RateLimitMiddleware limits the number of requests per second for a given IP.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		limiter := getLimiter(ip)
		if !limiter.Allow() {
			writeJSONError(w, http.StatusTooManyRequests, "Too many requests")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// writeJSONError returns an error message in JSON format.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
