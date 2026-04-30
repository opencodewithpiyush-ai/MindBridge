package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"mindbridge/infrastructure/repositories"
)

// RateLimit middleware limits repeated requests from the same IP.
// It uses Redis INCR with TTL to count attempts within a sliding window.
func RateLimit(redisClient *repositories.RedisClient, maxAttempts int, window time.Duration) gin.HandlerFunc {
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
