package handlers

import (
	"log"
	domainRepo "mindbridge/domain/repositories"
	"mindbridge/infrastructure/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService domainRepo.IAuthService, redisClient *repositories.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		logger := log.New(log.Writer(), "[AuthMiddleware] ", log.LstdFlags)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Printf("Missing Authorization header | IP: %s", clientIP)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Printf("Invalid Authorization header format | IP: %s", clientIP)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid Authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		if redisClient != nil {
			isValid, err := redisClient.IsSessionValid(token)
			if err == nil && !isValid {
				logger.Printf("Invalid session | IP: %s", clientIP)
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Session expired or invalid",
				})
				c.Abort()
				return
			}
		} else {
			userID, err := authService.ValidateToken(token)
			if err != nil {
				logger.Printf("Invalid token | IP: %s | Error: %v", clientIP, err)
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"error":   "Invalid or expired token",
				})
				c.Abort()
				return
			}
			c.Set("userID", userID)
			logger.Printf("Authenticated | UserID: %s | IP: %s", userID, clientIP)
			c.Next()
			return
		}

		userID, err := authService.ValidateToken(token)
		if err != nil {
			logger.Printf("Invalid token | IP: %s | Error: %v", clientIP, err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("userID", userID)
		logger.Printf("Authenticated | UserID: %s | IP: %s", userID, clientIP)
		c.Next()
	}
}
