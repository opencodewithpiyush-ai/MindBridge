package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Model struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Display string `json:"display"`
}

var (
	Models          []Model
	WebSocketURL    string
	FileUploadURL   string
	FileBaseURL     string
	UserAgent       string
	Origin          string
	ServerHost      string
	ServerPort      int
	MongoDBUsername string
	MongoDBPassword string
	MongoDBCluster  string
	MongoDBName     string
	JWTSecret       string
	RedisHost       string
	RedisPort       int
	RedisUsername   string
	RedisPassword   string
)

func InitConfig() {
	LoadEnv()

	Models = []Model{
		{ID: 1, Name: "gateway-gpt-5-4", Display: "GPT-5.4 (Latest)"},
		{ID: 2, Name: "gateway-gpt-5-3", Display: "GPT-5.3"},
		{ID: 3, Name: "gateway-gpt-5-1", Display: "GPT-5.1"},
		{ID: 4, Name: "gateway-gpt-5", Display: "GPT-5"},
		{ID: 5, Name: "gateway-gpt-4o", Display: "GPT-4o"},
		{ID: 6, Name: "gateway-gpt-4o-mini", Display: "GPT-4o Mini"},
		{ID: 7, Name: "gateway-grok-4", Display: "Grok-4 (xAI)"},
		{ID: 8, Name: "gateway-claude-sonnet-4-6", Display: "Claude Sonnet 4.6"},
		{ID: 9, Name: "gateway-claude-opus-4-5", Display: "Claude Opus 4.5"},
		{ID: 10, Name: "gateway-claude-opus-4-1", Display: "Claude Opus 4.1"},
		{ID: 11, Name: "gateway-deepseek-r1", Display: "DeepSeek R1"},
		{ID: 12, Name: "gateway-gemini-3-1-pro", Display: "Gemini 3.1 Pro"},
		{ID: 13, Name: "gateway-gemini-3-pro", Display: "Gemini 3 Pro"},
		{ID: 14, Name: "gateway-gemini-2.5-flash", Display: "Gemini 2.5 Flash"},
		{ID: 15, Name: "gateway-qwen-3-max", Display: "Qwen 3 Max"},
		{ID: 16, Name: "gateway-llama-3-3-70b-versatile", Display: "Llama 3.3 70B"},
		{ID: 17, Name: "gateway-deepinfra-kimi-k2", Display: "Kimi K2"},
	}

	WebSocketURL = getEnv("WEBSOCKET_URL", "wss://agents.use.ai/agents/budget-agent")
	FileUploadURL = getEnv("FILE_UPLOAD_URL", "https://files.use.ai/upload")
	FileBaseURL = getEnv("FILE_BASE_URL", "https://files.use.ai")
	UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"
	Origin = "https://use.ai"

	ServerHost = getEnv("SERVER_HOST", "127.0.0.1")
	ServerPort = getEnvInt("SERVER_PORT", 5000)

	MongoDBUsername = getEnv("MONGO_USERNAME", "")
	MongoDBPassword = getEnv("MONGO_PASSWORD", "")
	MongoDBCluster = getEnv("MONGO_CLUSTER", "")
	MongoDBName = getEnv("MONGO_DB", "mindbridge")
	JWTSecret = getEnv("JWT_SECRET", "mindbridge-secret-key")

	RedisHost = getEnv("REDIS_HOST", "localhost")
	RedisPort = getEnvInt("REDIS_PORT", 6379)
	RedisUsername = getEnv("REDIS_USERNAME", "")
	RedisPassword = getEnv("REDIS_PASSWORD", "")

	fmt.Printf("[Config] Loaded - Username: %s, Cluster: %s\n", MongoDBUsername, MongoDBCluster)
}

func LoadEnv() {
	envPath := filepath.Join(".", ".env")

	_, err := os.Stat(envPath)
	if err != nil {
		fmt.Println("Warning: .env file not found:", err)
		return
	}

	err = godotenv.Overload(envPath)
	if err != nil {
		fmt.Println("Warning: .env file not loaded:", err)
		return
	}
}

func GetMongoURI() string {
	if MongoDBUsername != "" && MongoDBPassword != "" && MongoDBCluster != "" {
		return fmt.Sprintf("mongodb+srv://%s:%s@%s/?appName=MindBridge", MongoDBUsername, MongoDBPassword, MongoDBCluster)
	}
	return "mongodb://localhost:27017"
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intVal int
		fmt.Sscanf(value, "%d", &intVal)
		return intVal
	}
	return defaultValue
}
