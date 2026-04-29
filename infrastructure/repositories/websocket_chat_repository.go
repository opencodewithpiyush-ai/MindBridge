package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mindbridge/config"
	domainRepo "mindbridge/domain/repositories"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type ChatRepository struct {
	idGenerator domainRepo.IIDGenerator
	headers     http.Header
}

func NewChatRepository() domainRepo.IChatRepository {
	return &ChatRepository{
		idGenerator: &IDGeneratorImpl{},
		headers: http.Header{
			"User-Agent": []string{config.UserAgent},
			"Origin":     []string{config.Origin},
		},
	}
}

func (r *ChatRepository) SendMessage(query, model, userID, email string) (string, string, error) {
	return r.SendMessageStream(query, model, userID, email, nil)
}

func (r *ChatRepository) SendMessageStream(query, model, userID, email string, onChunk func(string)) (string, string, error) {
	chatID := r.idGenerator.Generate()
	deviceID := r.idGenerator.Generate()

	uri := fmt.Sprintf("%s/%s?userId=%s&userType=regular&userEmail=%s&planType=free&isTestUser=false",
		config.WebSocketURL, chatID, userID, email)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, uri, r.headers)
	if err != nil {
		return "", "", err
	}
	defer conn.Close()

	prewarm := map[string]interface{}{"type": "prewarm", "chatId": chatID}
	if err := conn.WriteJSON(prewarm); err != nil {
		return "", "", err
	}

	chatMsg := r.buildChatMessage(chatID, userID, email, deviceID, query, model)
	if err := conn.WriteJSON(chatMsg); err != nil {
		return "", "", err
	}

	fullResponse := ""
	chatTitle := ""
	chunkCount := 0

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			continue
		}

		dataType, _ := data["type"].(string)

		if dataType == "stream-start" {
			log.Printf("Stream started | ChatID: %s", chatID)
		}

		if dataType == "data-chat-title-update" {
			if titleData, ok := data["data"].(map[string]interface{}); ok {
				if title, ok := titleData["title"].(string); ok {
					chatTitle = title
				}
			}
		}

		if chunk, ok := data["chunk"].(map[string]interface{}); ok {
			if chunkType, ok := chunk["type"].(string); ok {
				if chunkType == "text-delta" {
					if delta, ok := chunk["delta"].(string); ok {
						fullResponse += delta
						chunkCount++
						if onChunk != nil {
							onChunk(fullResponse)
						}
					}
				} else if chunkType == "finish" {
					return chatTitle, fullResponse, nil
				}
			}
		}
	}

	return chatTitle, fullResponse, nil
}

func (r *ChatRepository) SendMessageWithFiles(query, model, userID, email string, files []map[string]string, onChunk func(string)) (string, string, error) {
	chatID := r.idGenerator.Generate()
	deviceID := r.idGenerator.Generate()

	uri := fmt.Sprintf("%s/%s?userId=%s&userType=regular&userEmail=%s&planType=free&isTestUser=false",
		config.WebSocketURL, chatID, userID, email)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, uri, r.headers)
	if err != nil {
		return "", "", err
	}
	defer conn.Close()

	prewarm := map[string]interface{}{"type": "prewarm", "chatId": chatID}
	if err := conn.WriteJSON(prewarm); err != nil {
		return "", "", err
	}

	chatMsg := r.buildChatMessageWithFiles(chatID, userID, email, deviceID, query, model, files)
	if err := conn.WriteJSON(chatMsg); err != nil {
		return "", "", err
	}

	fullResponse := ""
	chatTitle := ""
	chunkCount := 0

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			continue
		}

		dataType, _ := data["type"].(string)

		if dataType == "stream-start" {
			log.Printf("Stream started with files | ChatID: %s", chatID)
		}

		if dataType == "data-chat-title-update" {
			if titleData, ok := data["data"].(map[string]interface{}); ok {
				if title, ok := titleData["title"].(string); ok {
					chatTitle = title
				}
			}
		}

		if chunk, ok := data["chunk"].(map[string]interface{}); ok {
			if chunkType, ok := chunk["type"].(string); ok {
				if chunkType == "text-delta" {
					if delta, ok := chunk["delta"].(string); ok {
						fullResponse += delta
						chunkCount++
						if onChunk != nil {
							onChunk(fullResponse)
						}
					}
				} else if chunkType == "finish" {
					return chatTitle, fullResponse, nil
				}
			}
		}
	}

	return chatTitle, fullResponse, nil
}

func (r *ChatRepository) buildChatMessage(chatID, userID, email, deviceID, query, model string) map[string]interface{} {
	return map[string]interface{}{
		"abortSignal":           map[string]interface{}{},
		"chatId":                chatID,
		"userId":                userID,
		"email":                 email,
		"userType":              "regular",
		"userEmail":             email,
		"planType":              "free",
		"subscriptionStatus":    "active",
		"isFreemium":            false,
		"isTestUser":            false,
		"mixpanelUserId":        "",
		"deviceId":              deviceID,
		"isMobile":              false,
		"isWebSearchMode":       false,
		"isDeepResearchMode":    false,
		"isIncognito":           true,
		"isImageGenerationMode": false,
		"isStandaloneImageMode": false,
		"needsBlurPreview":      false,
		"deepResearchProcessor": "pro-fast",
		"selectedModel":         model,
		"disableReasoning":      false,
		"source":                "incognito",
		"locale":                "en",
		"messages": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"type": "text", "text": query},
				},
				"id":   "msg-001",
				"role": "user",
				"metadata": map[string]interface{}{
					"isDeepResearchMode":    false,
					"isWebSearchMode":       false,
					"isImageGenerationMode": false,
					"needsBlurPreview":      false,
					"deepResearchProcessor": "pro-fast",
				},
			},
		},
		"trigger": "submit-message",
	}
}

func (r *ChatRepository) buildChatMessageWithFiles(chatID, userID, email, deviceID, query, model string, files []map[string]string) map[string]interface{} {
	parts := []map[string]interface{}{}

	for _, file := range files {
		parts = append(parts, map[string]interface{}{
			"type":      "file",
			"filename":  file["name"],
			"mediaType": file["type"],
			"url":       file["url"],
		})
	}

	parts = append(parts, map[string]interface{}{
		"type": "text",
		"text": query,
	})

	return map[string]interface{}{
		"abortSignal":           map[string]interface{}{},
		"chatId":                chatID,
		"userId":                userID,
		"email":                 email,
		"userType":              "regular",
		"userEmail":             email,
		"planType":              "free",
		"subscriptionStatus":    "active",
		"isFreemium":            false,
		"isTestUser":            false,
		"mixpanelUserId":        "",
		"deviceId":              deviceID,
		"isMobile":              false,
		"isWebSearchMode":       false,
		"isDeepResearchMode":    false,
		"isImageGenerationMode": false,
		"isStandaloneImageMode": false,
		"needsBlurPreview":      false,
		"deepResearchProcessor": "pro-fast",
		"selectedModel":         model,
		"disableReasoning":      false,
		"locale":                "en",
		"messages": []map[string]interface{}{
			{
				"parts": parts,
				"id":    "msg-001",
				"role":  "user",
				"metadata": map[string]interface{}{
					"isDeepResearchMode":    false,
					"isWebSearchMode":       false,
					"isImageGenerationMode": false,
					"needsBlurPreview":      false,
					"deepResearchProcessor": "pro-fast",
				},
			},
		},
		"trigger": "submit-message",
		"source":  "chat_page",
	}
}

type IDGeneratorImpl struct{}

func (g *IDGeneratorImpl) Generate() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (r *ChatRepository) SendMessageStreamRaw(query, model, userID, email string, onRawChunk func(map[string]interface{})) (string, string, error) {
	return r.sendMessageStreamRaw(query, model, userID, email, nil, onRawChunk)
}

func (r *ChatRepository) SendMessageWithFilesRaw(query, model, userID, email string, files []map[string]string, onRawChunk func(map[string]interface{})) (string, string, error) {
	return r.sendMessageStreamRaw(query, model, userID, email, files, onRawChunk)
}

func (r *ChatRepository) sendMessageStreamRaw(query, model, userID, email string, files []map[string]string, onRawChunk func(map[string]interface{})) (string, string, error) {
	chatID := r.idGenerator.Generate()
	deviceID := r.idGenerator.Generate()

	uri := fmt.Sprintf("%s/%s?userId=%s&userType=regular&userEmail=%s&planType=free&isTestUser=false",
		config.WebSocketURL, chatID, userID, email)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, uri, r.headers)
	if err != nil {
		return "", "", err
	}
	defer conn.Close()

	prewarm := map[string]interface{}{"type": "prewarm", "chatId": chatID}
	if err := conn.WriteJSON(prewarm); err != nil {
		return "", "", err
	}

	if files != nil {
		chatMsg := r.buildChatMessageWithFiles(chatID, userID, email, deviceID, query, model, files)
		if err := conn.WriteJSON(chatMsg); err != nil {
			return "", "", err
		}
	} else {
		chatMsg := r.buildChatMessage(chatID, userID, email, deviceID, query, model)
		if err := conn.WriteJSON(chatMsg); err != nil {
			return "", "", err
		}
	}

	fullResponse := ""
	chatTitle := ""
	var lastToolCall map[string]interface{}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var data map[string]interface{}
		if err := json.Unmarshal(message, &data); err != nil {
			continue
		}

		if onRawChunk != nil {
			onRawChunk(data)
		}

		dataType, _ := data["type"].(string)

		if dataType == "data-chat-title-update" {
			if titleData, ok := data["data"].(map[string]interface{}); ok {
				if title, ok := titleData["title"].(string); ok {
					chatTitle = title
				}
			}
		}

		if chunk, ok := data["chunk"].(map[string]interface{}); ok {
			chunkType, _ := chunk["type"].(string)
			if chunkType == "text-delta" {
				if delta, ok := chunk["delta"].(string); ok {
					fullResponse += delta
				}
			} else if chunkType == "tool-input-available" {
				lastToolCall = chunk
				if toolInput, ok := chunk["input"].(map[string]interface{}); ok {
					if result, ok := toolInput["resultURL"].(string); ok {
						fullResponse = "[Image generated] " + result
					} else if prompt, ok := toolInput["prompt"].(string); ok {
						fullResponse = "Editing: " + prompt
					}
				}
			} else if chunkType == "finish" {
				break
			}
		}
	}

	if fullResponse == "" && lastToolCall != nil {
		dataBytes, _ := json.Marshal(lastToolCall)
		fullResponse = string(dataBytes)
	}

	return chatTitle, fullResponse, nil
}
