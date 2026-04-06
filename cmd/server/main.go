package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mindbridge/config"
	"mindbridge/infrastructure/repositories"
	"mindbridge/presentation/handlers"
	"mindbridge/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	utils.SetupLogging()
	logger := log.New(os.Stdout, "[Main] ", log.LstdFlags)

	logger.Println("=================================================")
	logger.Println("MindBridge API Server Starting...")
	logger.Println("=================================================")

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	chatRepo := repositories.NewChatRepository()
	fileRepo := repositories.NewFileRepository()
	handlers.SetupRoutes(router, chatRepo, fileRepo)

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
