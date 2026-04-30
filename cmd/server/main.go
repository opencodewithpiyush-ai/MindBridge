package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mindbridge/config"
	"mindbridge/infrastructure/di"
	"mindbridge/middleware"
	"mindbridge/presentation/handlers"
	"mindbridge/utils"

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
