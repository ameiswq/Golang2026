package utils

import (
	"net/http"
	"sync"
	"time"
	"github.com/gin-gonic/gin"
)

type clientData struct {
	Count int
	ExpiresAt time.Time
}

type RateLimiter struct {
	mu sync.Mutex
	clients map[string]*clientData
	limit int
	window time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*clientData),
		limit: limit,
		window: window,
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if userID, exists := c.Get("userID"); exists {
			if s, ok := userID.(string); ok && s != "" {
				key = "user:" + s
			}
		}
		now := time.Now()
		rl.mu.Lock()
		defer rl.mu.Unlock()
		data, exists := rl.clients[key]
		if !exists || now.After(data.ExpiresAt) {
			rl.clients[key] = &clientData{
				Count:     1,
				ExpiresAt: now.Add(rl.window),
			}
			c.Next()
			return
		}

		if data.Count >= rl.limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
			})
			return
		}

		data.Count++
		c.Next()
	}
}