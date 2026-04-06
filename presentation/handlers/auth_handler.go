package handlers

import (
	"log"
	"mindbridge/application/dto"
	"mindbridge/application/usecases"
	domainRepo "mindbridge/domain/repositories"
	"mindbridge/infrastructure/repositories"
	"mindbridge/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(router *gin.Engine, authUseCase *usecases.AuthUseCase, authService domainRepo.IAuthService, redisClient *repositories.RedisClient) {
	router.POST("/auth/register", registerHandler(authUseCase))
	router.POST("/auth/login", loginHandler(authUseCase))
	router.POST("/auth/logout", AuthMiddleware(authService, redisClient), logoutHandler(redisClient, authService))
}

func registerHandler(useCase *usecases.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		logger := log.New(log.Writer(), "[RegisterHandler] ", log.LstdFlags)
		logger.Printf("Register endpoint called | IP: %s", clientIP)

		var request dto.RegisterRequestDTO
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Printf("Invalid request body | IP: %s | Error: %v", clientIP, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid request body: " + err.Error(),
			})
			return
		}

		validationErrors := utils.ValidateRegister(request.Name, request.Username, request.Email, request.Password)
		if len(validationErrors) > 0 {
			logger.Printf("Validation failed | IP: %s | Errors: %v", clientIP, validationErrors)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"errors":  validationErrors,
			})
			return
		}

		result := useCase.Register(request)

		if result.Success {
			logger.Printf("Registration successful | Email: %s | IP: %s", request.Email, clientIP)
			c.JSON(http.StatusCreated, gin.H{
				"success": true,
				"message": "User registered successfully",
				"data":    result.Data,
			})
		} else {
			logger.Printf("Registration failed | Error: %s | IP: %s", result.Error, clientIP)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   result.Error,
			})
		}
	}
}

func loginHandler(useCase *usecases.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		logger := log.New(log.Writer(), "[LoginHandler] ", log.LstdFlags)
		logger.Printf("Login endpoint called | IP: %s", clientIP)

		var request dto.LoginRequestDTO
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Printf("Invalid request body | IP: %s | Error: %v", clientIP, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid request body: " + err.Error(),
			})
			return
		}

		result := useCase.Login(request)

		if result.Success {
			logger.Printf("Login successful | Email: %s | IP: %s", request.Email, clientIP)
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Login successful",
				"data":    result.Data,
			})
		} else {
			logger.Printf("Login failed | Error: %s | IP: %s", result.Error, clientIP)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   result.Error,
			})
		}
	}
}

func logoutHandler(redisClient *repositories.RedisClient, authService domainRepo.IAuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		logger := log.New(log.Writer(), "[LogoutHandler] ", log.LstdFlags)

		userID, exists := c.Get("userID")
		if !exists {
			logger.Printf("User not authenticated | IP: %s", clientIP)
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Unauthorized",
			})
			return
		}

		authHeader := c.GetHeader("Authorization")
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) > 1 && redisClient != nil {
			token := tokenParts[1]
			redisClient.DeleteSession(token)
		}

		logger.Printf("Logout successful | UserID: %s | IP: %s", userID, clientIP)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Logout successful",
		})
	}
}
