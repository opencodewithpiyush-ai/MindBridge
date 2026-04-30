package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mindbridge/config"
	"mindbridge/infrastructure/repositories"
	"mindbridge/middleware"
	"mindbridge/presentation/handlers"
	"mindbridge/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	utils.SetupLogging()
	config.InitConfig()
	logger := log.New(os.Stdout, "[Main] ", log.LstdFlags)

	logger.Println("=================================================")
	logger.Println("MindBridge API Server Starting...")
	logger.Println("=================================================")

	logger.Printf("MongoDB URI: %s", config.GetMongoURI())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.GetMongoURI()))
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		logger.Fatalf("Failed to ping MongoDB: %v", err)
	}
	logger.Println("Connected to MongoDB")

	collection := client.Database(config.MongoDBName).Collection("users")

	userRepo := repositories.NewUserRepository(collection)

	redisClient, err := repositories.NewRedisClient()
	if err != nil {
		logger.Printf("Warning: Redis not connected: %v", err)
	} else {
		defer redisClient.Close()
	}

	jwtService := repositories.NewJWTService(config.JWTSecret, redisClient)

	router := gin.Default()

	router.Use(middleware.RequestID())

	// Auth routes with rate limiting
	authGroup := router.Group("")
	authGroup.Use(middleware.RateLimit(redisClient, config.RateLimitMax, config.RateLimitWindow))

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authUseCase := handlers.NewAuthUseCaseHandler(userRepo, jwtService, redisClient)
	handlers.SetupAuthRoutes(authGroup, authUseCase, jwtService, redisClient)

	chatRepo := repositories.NewChatRepository()
	fileRepo := repositories.NewFileRepository()

	protected := router.Group("/")
	protected.Use(handlers.AuthMiddleware(jwtService, redisClient))
	handlers.SetupChatRoutes(protected, chatRepo, fileRepo)

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
