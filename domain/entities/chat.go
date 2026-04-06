package entities

type ChatMessage struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequestEntity struct {
	Query  string  `json:"query"`
	Model  string  `json:"model"`
	UserID *string `json:"user_id,omitempty"`
	Email  *string `json:"email,omitempty"`
}

type ChatResponseEntity struct {
	ChatID   string `json:"chat_id"`
	Title    string `json:"title"`
	Model    string `json:"model"`
	Query    string `json:"query"`
	Response string `json:"response"`
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
}
