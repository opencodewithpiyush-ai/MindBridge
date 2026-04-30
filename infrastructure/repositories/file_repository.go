package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"mindbridge/config"
	domainRepo "mindbridge/domain/repositories"
	"net/http"
	"time"
)

type FileRepository struct{}

func NewFileRepository() domainRepo.IFileRepository {
	return &FileRepository{}
}

func (r *FileRepository) GetFileURL(key string) string {
	return fmt.Sprintf("%s/files/%s", config.FileBaseURL, key)
}

func (r *FileRepository) UploadFile(fileName, fileType string, fileData []byte) (string, string, error) {
	logger := log.New(log.Writer(), "[FileRepository] ", log.LstdFlags)
	logger.Printf("Uploading file | Name: %s | Type: %s | Size: %d bytes", fileName, fileType, len(fileData))

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	namePart, err := writer.CreateFormField("name")
	if err != nil {
		return "", "", err
	}
	if _, err := namePart.Write([]byte(fileName)); err != nil {
		return "", "", err
	}

	typePart, err := writer.CreateFormField("type")
	if err != nil {
		return "", "", err
	}
	if _, err := typePart.Write([]byte(fileType)); err != nil {
		return "", "", err
	}

	filePart, err := writer.CreatePart(map[string][]string{
		"Content-Type":        {fileType},
		"Content-Disposition": {fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName)},
	})
	if err != nil {
		return "", "", err
	}
	if _, err := filePart.Write(fileData); err != nil {
		return "", "", err
	}

	writer.Close()

	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest("POST", config.FileUploadURL, &buf)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Origin", "https://use.ai")
	req.Header.Set("Referer", "https://use.ai/")

	resp, err := client.Do(req)
	if err != nil {
		logger.Printf("Request failed: %v", err)
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logger.Printf("Status: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("upload failed: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Printf("[FileRepository] json unmarshal error: %v", err)
		return "", "", err
	}

	key, _ := result["key"].(string)
	fileURL, _ := result["url"].(string)

	return key, fileURL, nil
}
