package usecases

import (
	"mindbridge/application/dto"
	"mindbridge/domain/repositories"
	"mindbridge/config"
	"strings"
)

type ProcessChatUseCase struct {
	chatRepository  repositories.IChatRepository
	idGenerator     repositories.IIDGenerator
	emailGenerator  repositories.IEmailGenerator
	availableModels []string
}

func NewProcessChatUseCase(
	chatRepo repositories.IChatRepository,
	idGen repositories.IIDGenerator,
	emailGen repositories.IEmailGenerator,
	availableModels []string,
) *ProcessChatUseCase {
	return &ProcessChatUseCase{
		chatRepository:  chatRepo,
		idGenerator:     idGen,
		emailGenerator:  emailGen,
		availableModels: availableModels,
	}
}

func (u *ProcessChatUseCase) Execute(request dto.ChatRequestDTO) dto.ChatResponseDTO {
	if request.Query == "" {
		return dto.ChatResponseDTO{Success: false, Error: "Query cannot be empty"}
	}

	validModel := false
	modelName := request.Model
	if modelName == "" {
		modelName = "claude-opus-4-1"
	}
	modelName = config.ResolveModel(modelName)
	for _, m := range u.availableModels {
		if m == modelName {
			validModel = true
			break
		}
	}
	if !validModel {
		return dto.ChatResponseDTO{Success: false, Error: "Invalid model. Choose from: " + strings.Join(u.availableModels, ", ")}
	}

	userID := request.UserID
	if userID == nil || *userID == "" {
		id := u.idGenerator.Generate()
		userID = &id
	}

	email := request.Email
	if email == nil || *email == "" {
		e := u.emailGenerator.Generate()
		email = &e
	}

	title, response, err := u.chatRepository.SendMessage(request.Query, modelName, *userID, *email)
	if err != nil {
		return dto.ChatResponseDTO{Success: false, Error: err.Error()}
	}

	chatID := u.idGenerator.Generate()

	return dto.ChatResponseDTO{
		Success: true,
		Data: map[string]interface{}{
			"chatId":   chatID,
			"title":    title,
			"model":    modelName,
			"query":    request.Query,
			"response": response,
		},
	}
}
