package handlers

import (
	"encoding/json"
	"fmt"
	"mindbridge/application/dto"
	"mindbridge/application/usecases"
	"mindbridge/config"
	domainRepo "mindbridge/domain/repositories"
	"mindbridge/infrastructure/generators"
	"mindbridge/infrastructure/repositories"
	"mindbridge/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupChatRoutes(router *gin.RouterGroup, chatRepo domainRepo.IChatRepository, fileRepo domainRepo.IFileRepository) {
	idGenerator := generators.NewIDGenerator()
	emailGenerator := generators.NewEmailGenerator()

	router.POST("/chat/stream-raw", chatStreamRawHandler(chatRepo, idGenerator, emailGenerator))
	router.POST("/upload", fileUploadHandler(fileRepo))
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
			"GET /":                 "API info",
			"GET /models":           "List available AI models",
			"POST /auth/register":   "Register new user",
			"POST /auth/login":      "Login and get session",
			"POST /auth/logout":     "Logout and destroy session",
			"POST /chat/stream-raw": "Chat with full streaming (protected)",
			"POST /upload":          "Upload file (protected)",
		},
	})
}

func ListModelsHandler(c *gin.Context) {
	models := make([]gin.H, len(config.Models))
	for i, m := range config.Models {
		public := config.InternalToPublic[m.Name]
		if public == "" {
			public = m.Name
		}
		models[i] = gin.H{"id": m.ID, "name": public, "display": m.Display}
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"models":  models,
	})
}

func chatStreamRawHandler(
	chatRepo domainRepo.IChatRepository,
	idGenerator domainRepo.IIDGenerator,
	emailGenerator domainRepo.IEmailGenerator,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
	requestID, _ := c.Get("request_id")
	logger := utils.WithRequestID("ChatStreamRawHandler", fmt.Sprint(requestID))
		logger.Printf("Raw stream endpoint called | IP: %s", clientIP)

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
			request.Model = "claude-opus-4-1"
		}
		request.Model = config.ResolveModel(request.Model)

		userID := idGenerator.Generate()
		email := emailGenerator.Generate()
		if request.UserID != nil {
			userID = *request.UserID
		}
		if request.Email != nil {
			email = *request.Email
		}

		logger.Printf("Processing raw stream request | Model: %s | Query: %s...", request.Model, truncate(request.Query, 30))

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
			title, response, err = chatRepo.SendMessageWithFilesRaw(
				request.Query,
				request.Model,
				userID,
				email,
				request.Files,
				func(chunk map[string]interface{}) {
					dataBytes, _ := json.Marshal(chunk)
					sendEvent("chunk", string(dataBytes))
				},
			)
		} else {
			title, response, err = chatRepo.SendMessageStreamRaw(
				request.Query,
				request.Model,
				userID,
				email,
				func(chunk map[string]interface{}) {
					dataBytes, _ := json.Marshal(chunk)
					sendEvent("chunk", string(dataBytes))
				},
			)
		}

		if err != nil {
			logger.Printf("Raw stream error | Error: %s | IP: %s", err, clientIP)
			sendEvent("error", fmt.Sprintf(`{"error": %q}`, err.Error()))
			return
		}

		logger.Printf("Raw stream complete | Title: %s | IP: %s", title, clientIP)
		sendEvent("done", fmt.Sprintf(`{"title": %q, "response": %q}`, title, response))
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func fileUploadHandler(fileRepo domainRepo.IFileRepository) gin.HandlerFunc {
	const maxUploadSize int64 = 10 << 20 // 10 MB

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		requestID, _ := c.Get("request_id")
		logger := utils.WithRequestID("FileUploadHandler", fmt.Sprint(requestID))

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

		if header.Size <= 0 {
			logger.Printf("Invalid file size | IP: %s | Size: %d", clientIP, header.Size)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid file size",
			})
			return
		}

		if header.Size > maxUploadSize {
			logger.Printf("File too large | IP: %s | Size: %d | Max: %d", clientIP, header.Size, maxUploadSize)
			c.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"success": false,
				"error":   "File size exceeds maximum allowed limit",
			})
			return
		}

		maxInt := int64(^uint(0) >> 1)
		if header.Size > maxInt {
			logger.Printf("File size overflows int | IP: %s | Size: %d", clientIP, header.Size)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"error":   "Invalid file size",
			})
			return
		}

		fileData := make([]byte, int(header.Size))
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
