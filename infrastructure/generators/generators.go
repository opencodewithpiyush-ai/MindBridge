package generators

import (
	"github.com/google/uuid"
	"mindbridge/domain/repositories"
)

type IDGenerator struct{}

func NewIDGenerator() repositories.IIDGenerator {
	return &IDGenerator{}
}

func (g *IDGenerator) Generate() string {
	return uuid.New().String()
}

type EmailGenerator struct{}

func NewEmailGenerator() repositories.IEmailGenerator {
	return &EmailGenerator{}
}

func (g *EmailGenerator) Generate() string {
	return "user_" + uuid.New().String()[:8] + "@example.com"
}
