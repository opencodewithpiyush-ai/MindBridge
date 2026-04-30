package utils

import (
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

var (
	logger     *log.Logger
	loggerMutex sync.RWMutex
)

func init() {
	logger = log.New(os.Stdout, "[MindBridge] ", log.LstdFlags|log.Lshortfile)
}

func SetupLogging() {
	f, err := os.OpenFile("mindbridge.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	logger.SetOutput(f)
}

func GetLogger(module string) *log.Logger {
	return log.New(os.Stdout, "["+module+"] ", log.LstdFlags)
}

// WithRequestID returns a new logger that includes the given request ID.
func WithRequestID(module, requestID string) *log.Logger {
	return log.New(os.Stdout, "["+module+"] ["+requestID+"] ", log.LstdFlags)
}

// NewRequestID generates a new UUID v4 string to be used as request ID.
func NewRequestID() string {
	return uuid.New().String()
}
