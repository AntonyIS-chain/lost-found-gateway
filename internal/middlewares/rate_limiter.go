package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Rate limiting configuration
const (
	limit     = 5                   // Max requests per period
	interval  = 60 * time.Second     // Time window (1 minute)
)

// ✅ In-memory rate limiter (Thread-Safe)
var requestCounts = make(map[string]int)
var mutex = sync.Mutex{}

// ✅ RateLimiterMiddleware restricts requests per IP
func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mutex.Lock()
		count := requestCounts[ip]
		if count >= limit {
			mutex.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests, try again later"})
			return
		}

		// Increment request count
		requestCounts[ip]++
		mutex.Unlock()

		// Reset count after interval
		go func() {
			time.Sleep(interval)
			mutex.Lock()
			requestCounts[ip]--
			mutex.Unlock()
		}()

		c.Next()
	}
}
