package middleware

import (
	"github.com/gin-gonic/gin"
	"mindbridge/utils"
)

// RequestID middleware generates a unique request ID and sets it as header X-Request-ID,
// and injects it into the context for use downstream.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := utils.NewRequestID()
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		// Also add to logger context if needed
		c.Next()
	}
}
