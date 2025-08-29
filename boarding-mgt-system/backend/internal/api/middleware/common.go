package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// RateLimit implements a simple in-memory rate limiter
func RateLimit() gin.HandlerFunc {
	type client struct {
		limiter  *rateLimiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Clean up old entries periodically
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: newRateLimiter(100, 100), // 100 requests per minute
			}
		}
		clients[ip].lastSeen = time.Now()
		
		if !clients[ip].limiter.allow() {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		mu.Unlock()
		
		c.Next()
	}
}

// Simple token bucket rate limiter
type rateLimiter struct {
	tokens    float64
	capacity  float64
	rate      float64
	lastCheck time.Time
	mu        sync.Mutex
}

func newRateLimiter(rate, capacity float64) *rateLimiter {
	return &rateLimiter{
		tokens:    capacity,
		capacity:  capacity,
		rate:      rate,
		lastCheck: time.Now(),
	}
}

func (rl *rateLimiter) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(rl.lastCheck).Seconds()
	rl.lastCheck = now
	
	// Add tokens based on elapsed time
	rl.tokens += elapsed * (rl.rate / 60.0)
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}
	
	if rl.tokens >= 1.0 {
		rl.tokens--
		return true
	}
	
	return false
}

// ErrorHandler provides consistent error responses
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// Handle any errors that occurred during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Determine status code
			status := c.Writer.Status()
			if status == http.StatusOK {
				status = http.StatusInternalServerError
			}
			
			// Send error response
			c.JSON(status, gin.H{
				"error":      err.Error(),
				"request_id": c.GetString("request_id"),
			})
		}
	}
}

// Timeout adds a timeout to requests
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a copy of the context with timeout
		ctx, cancel := c.Copy(), func() {}
		c.Request = c.Request.WithContext(ctx.Request.Context())
		
		// Set up timeout
		timer := time.NewTimer(timeout)
		done := make(chan struct{})
		
		go func() {
			c.Next()
			close(done)
		}()
		
		select {
		case <-timer.C:
			cancel()
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"error": "Request timeout",
			})
			c.Abort()
		case <-done:
			timer.Stop()
		}
	}
}

// CORS handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}