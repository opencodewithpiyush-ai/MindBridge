package repositories

type IChatRepository interface {
	SendMessage(query, model, userID, email string) (string, string, error)
	SendMessageStream(query, model, userID, email string, onChunk func(string)) (string, string, error)
	SendMessageWithFiles(query, model, userID, email string, files []map[string]string, onChunk func(string)) (string, string, error)
	SendMessageStreamRaw(query, model, userID, email string, onRawChunk func(map[string]interface{})) (string, string, error)
	SendMessageWithFilesRaw(query, model, userID, email string, files []map[string]string, onRawChunk func(map[string]interface{})) (string, string, error)
}

type IFileRepository interface {
	UploadFile(fileName, fileType string, fileData []byte) (string, string, error)
	GetFileURL(key string) string
}

type IIDGenerator interface {
	Generate() string
}

type IEmailGenerator interface {
	Generate() string
}
