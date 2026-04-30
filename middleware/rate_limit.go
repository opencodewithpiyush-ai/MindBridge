package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// redisCounter is an interface for the Redis operations needed by rate limiting.
// Using an interface allows easy mocking in tests.
type redisCounter interface {
	Incr(key string) (int64, error)
	Expire(key string, expiry time.Duration)
}

// RateLimit middleware limits repeated requests from the same IP.
// It uses Redis INCR with TTL to count attempts within a sliding window.
func RateLimit(redisClient redisCounter, maxAttempts int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		if redisClient == nil {
			// No Redis – skip rate limiting
			c.Next()
			return
		}

		ip := c.ClientIP()
		key := "rate_limit:" + ip + ":" + c.FullPath()

		count, err := redisClient.Incr(key)
		if err != nil {
			// Redis error – allow request (fail open)
			c.Next()
			return
		}

		if count == 1 {
			redisClient.Expire(key, window)
		}

		if count > int64(maxAttempts) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Too many requests, please try again later.",
			})
			return
		}

		c.Next()
	}
}
