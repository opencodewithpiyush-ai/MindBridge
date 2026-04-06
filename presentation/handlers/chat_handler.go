package handlers

import (
	"fmt"
	"log"
	"mindbridge/application/dto"
	"mindbridge/application/usecases"
	"mindbridge/config"
	domainRepo "mindbridge/domain/repositories"
	"mindbridge/infrastructure/generators"
	"mindbridge/infrastructure/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, chatRepo domainRepo.IChatRepository, fileRepo domainRepo.IFileRepository) {
	router.GET("/", IndexHandler)
	router.GET("/models", ListModelsHandler)
	router.POST("/upload", fileUploadHandler(fileRepo))
	router.GET("/file/:key", fileGetHandler(fileRepo))
}

func SetupChatRoutes(router *gin.RouterGroup, chatRepo domainRepo.IChatRepository, fileRepo domainRepo.IFileRepository) {
	idGenerator := generators.NewIDGenerator()
	emailGenerator := generators.NewEmailGenerator()

	availableModels := make([]string, len(config.Models))
	for i, m := range config.Models {
		availableModels[i] = m.Name
	}

	processChatUseCase := usecases.NewProcessChatUseCase(
		chatRepo,
		idGenerator,
		emailGenerator,
		availableModels,
	)

	router.POST("/chat", chatHandler(processChatUseCase))
	router.POST("/chat/stream", chatStreamHandler(chatRepo, idGenerator, emailGenerator, availableModels))
	router.POST("/upload", fileUploadHandler(fileRepo))
	router.GET("/file/:key", fileGetHandler(fileRepo))
}

func NewAuthUseCaseHandler(userRepo domainRepo.IUserRepository, authService domainRepo.IAuthService, redisClient *repositories.RedisClient) *usecases.AuthUseCase {
	return usecases.NewAuthUseCase(userRepo, authService, redisClient)
}

func IndexHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":     "MindBridge API Server",
		"version":     "1.1.2",
		"description": "AI Chat Gateway with Multi-Model Support",
		"owner": gin.H{
			"name":   "Piyush Makwana",
			"email":  "piyushmakwana@mindbridge.ai",
			"github": "https://github.com/piyushmakwana",
		},
		"endpoints": gin.H{
			"GET /":               "API info",
			"GET /models":         "List available AI models",
			"POST /auth/register": "Register new user",
			"POST /auth/login":    "Login and get session",
			"POST /auth/logout":   "Logout and destroy session",
			"POST /chat":          "Send chat message (protected)",
			"POST /chat/stream":   "Chat with streaming (protected)",
			"POST /upload":        "Upload file (protected)",
			"GET /file/:key":      "Get file URL (protected)",
		},
	})
}

func ListModelsHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"models":  config.Models,
	})
}

func chatHandler(useCase *usecases.ProcessChatUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		logger := log.New(log.Writer(), "[ChatHandler] ", log.LstdFlags)
		logger.Printf("Chat endpoint called | IP: %s", clientIP)

		var request dto.ChatRequestDTO
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Printf("Invalid request body | IP: %s | Error: %v", clientIP, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Request body is required",
			})
			return
		}

		if request.Model == "" {
			request.Model = "gateway-claude-opus-4-1"
		}

		logger.Printf("Processing request | Model: %s | Query: %s...", request.Model, truncate(request.Query, 30))

		result := useCase.Execute(request)

		if result.Success {
			logger.Printf("Request successful | Title: %v | IP: %s", result.Data.(map[string]interface{})["title"], clientIP)
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"data":    result.Data,
			})
		} else {
			logger.Printf("Request failed | Error: %s | IP: %s", result.Error, clientIP)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   result.Error,
			})
		}
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func chatStreamHandler(
	chatRepo domainRepo.IChatRepository,
	idGenerator domainRepo.IIDGenerator,
	emailGenerator domainRepo.IEmailGenerator,
	availableModels []string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		logger := log.New(log.Writer(), "[ChatStreamHandler] ", log.LstdFlags)
		logger.Printf("Stream endpoint called | IP: %s", clientIP)

		var request dto.ChatRequestDTO
		if err := c.ShouldBindJSON(&request); err != nil {
			logger.Printf("Invalid request body | IP: %s | Error: %v", clientIP, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Request body is required",
			})
			return
		}

		if request.Model == "" {
			request.Model = "gateway-claude-opus-4-1"
		}

		userID := idGenerator.Generate()
		email := emailGenerator.Generate()
		if request.UserID != nil {
			userID = *request.UserID
		}
		if request.Email != nil {
			email = *request.Email
		}

		logger.Printf("Processing stream request | Model: %s | Query: %s...", request.Model, truncate(request.Query, 30))

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")

		flusher, ok := c.Writer.(interface{ Flush() })
		if !ok {
			logger.Printf("Streaming not supported | IP: %s", clientIP)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Streaming not supported",
			})
			return
		}

		sendEvent := func(event string, data string) {
			fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, data)
			flusher.Flush()
		}

		sendEvent("connected", `{"status": "connected"}`)

		var title, response string
		var err error

		if len(request.Files) > 0 {
			title, response, err = chatRepo.SendMessageWithFiles(
				request.Query,
				request.Model,
				userID,
				email,
				request.Files,
				func(chunk string) {
					sendEvent("chunk", fmt.Sprintf(`{"chunk": %q}`, chunk))
				},
			)
		} else {
			title, response, err = chatRepo.SendMessageStream(
				request.Query,
				request.Model,
				userID,
				email,
				func(chunk string) {
					sendEvent("chunk", fmt.Sprintf(`{"chunk": %q}`, chunk))
				},
			)
		}

		if err != nil {
			logger.Printf("Stream error | Error: %s | IP: %s", err, clientIP)
			sendEvent("error", fmt.Sprintf(`{"error": %q}`, err.Error()))
			return
		}

		logger.Printf("Stream complete | Title: %s | IP: %s", title, clientIP)
		sendEvent("done", fmt.Sprintf(`{"title": %q, "response": %q}`, title, response))
	}
}

func fileUploadHandler(fileRepo domainRepo.IFileRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		logger := log.New(log.Writer(), "[FileUploadHandler] ", log.LstdFlags)

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			logger.Printf("Failed to get file | IP: %s | Error: %v", clientIP, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "File is required",
			})
			return
		}
		defer file.Close()

		fileData := make([]byte, header.Size)
		if _, err := file.Read(fileData); err != nil {
			logger.Printf("Failed to read file | IP: %s | Error: %v", clientIP, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Failed to read file",
			})
			return
		}

		fileName := c.PostForm("name")
		if fileName == "" {
			fileName = header.Filename
		}
		fileType := c.PostForm("type")
		if fileType == "" {
			fileType = header.Header.Get("Content-Type")
			if fileType == "" {
				fileType = "image/jpeg"
			}
		}

		logger.Printf("Uploading file | Name: %s | Type: %s | Size: %d", fileName, fileType, len(fileData))

		key, url, err := fileRepo.UploadFile(fileName, fileType, fileData)
		if err != nil {
			logger.Printf("Upload failed | Error: %s | IP: %s", err, clientIP)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		logger.Printf("Upload successful | Key: %s | IP: %s", key, clientIP)
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"key":     key,
			"url":     url,
		})
	}
}

func fileGetHandler(fileRepo domainRepo.IFileRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "File key is required",
			})
			return
		}

		fileURL := fileRepo.GetFileURL(key)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"url":     fileURL,
		})
	}
}
