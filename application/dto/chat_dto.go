package dto

type ChatRequestDTO struct {
	Query  string              `json:"query" binding:"required"`
	Model  string              `json:"model"`
	UserID *string             `json:"userId,omitempty"`
	Email  *string             `json:"email,omitempty"`
	Files  []map[string]string `json:"files,omitempty"`
}

type ChatResponseDTO struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
