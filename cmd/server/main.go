package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mindbridge/config"
	"mindbridge/infrastructure/di"
	"mindbridge/middleware"
	"mindbridge/presentation/handlers"
	"mindbridge/utils"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.SetupLogging()
	config.InitConfig()
	logger := log.New(os.Stdout, "[Main] ", log.LstdFlags)

	logger.Println("=================================================")
	logger.Println("MindBridge API Server Starting...")
	logger.Println("=================================================")

	container, err := di.NewContainer()
	if err != nil {
		logger.Fatalf("Failed to create container: %v", err)
	}
	defer func() {
		if container.RedisClient != nil {
			container.RedisClient.Close()
		}
	}()

	logger.Printf("MongoDB URI: %s", config.GetMongoURI())
	logger.Println("Connected to MongoDB")

	router := gin.Default()

	// CORS — must be first so OPTIONS preflight is handled before auth/rate-limit
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:3000",
			"https://test-mindbridge-v1-1.onrender.com",
			"https://*.onrender.com",
		},
		AllowOriginFunc: func(origin string) bool {
			// Allow all origins in development; lock down in production via env
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.RequestID())

	// Auth routes with rate limiting
	authGroup := router.Group("")
	authGroup.Use(middleware.RateLimit(container.RedisClient, config.RateLimitMax, config.RateLimitWindow))

	handlers.SetupAuthRoutes(authGroup, container.AuthUseCase, container.AuthService, container.RedisClient)

	protected := router.Group("/")
	protected.Use(handlers.AuthMiddleware(container.AuthService, container.RedisClient))
	handlers.SetupChatRoutes(protected, container.ChatRepo, container.FileRepo)

	router.GET("/", handlers.IndexHandler)
	router.GET("/models", handlers.ListModelsHandler)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	go func() {
		logger.Printf("MindBridge running at http://%s:%d", config.ServerHost, config.ServerPort)
		if err := router.Run(fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)); err != nil {
			logger.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")
}
