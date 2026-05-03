package handlers

import (
	"fmt"
	"mindbridge/application/dto"
	"mindbridge/application/usecases"
	domainRepo "mindbridge/domain/repositories"
	"mindbridge/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SetupAuthRoutes(router *gin.RouterGroup, authUseCase *usecases.AuthUseCase, authService domainRepo.IAuthService) {
	router.POST("/auth/register", registerHandler(authUseCase))
	router.POST("/auth/login", loginHandler(authUseCase))
	router.POST("/auth/logout", AuthMiddleware(authService), logoutHandler(authUseCase))
}

func registerHandler(useCase *usecases.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		requestID, _ := c.Get("request_id")
		logger := utils.WithRequestID("RegisterHandler", fmt.Sprint(requestID))
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
		requestID, _ := c.Get("request_id")
		logger := utils.WithRequestID("LoginHandler", fmt.Sprint(requestID))
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

func logoutHandler(authUseCase *usecases.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		requestID, _ := c.Get("request_id")
		logger := utils.WithRequestID("LogoutHandler", fmt.Sprint(requestID))

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
		if len(tokenParts) > 1 {
			token := tokenParts[1]

			// Parse token without verification to extract the jti (session ID)
			parsed, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
			if err != nil {
				logger.Printf("Failed to parse token for logout | UserID: %s | Error: %v", userID, err)
			} else {
				if claims, ok := parsed.Claims.(jwt.MapClaims); ok {
					jti, _ := claims["jti"].(string)
					if jti != "" {
						if err := authUseCase.Logout(userID.(string), jti); err != nil {
							logger.Printf("Logout error | UserID: %s | Error: %v", userID, err)
						}
					}
				}
			}
		}

		logger.Printf("Logout successful | UserID: %s | IP: %s", userID, clientIP)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Logout successful",
		})
	}
}
