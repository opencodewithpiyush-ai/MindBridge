package usecases

import (
	"errors"
	"mindbridge/application/dto"
	domainRepo "mindbridge/domain/repositories"
	"mindbridge/infrastructure/repositories"
)

type AuthUseCase struct {
	userRepo    domainRepo.IUserRepository
	authService domainRepo.IAuthService
	redisClient *repositories.RedisClient
}

func NewAuthUseCase(userRepo domainRepo.IUserRepository, authService domainRepo.IAuthService, redisClient *repositories.RedisClient) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		authService: authService,
		redisClient: redisClient,
	}
}

type RegisterResult struct {
	Success bool
	Data    dto.AuthResponseDTO
	Error   string
}

type LoginResult struct {
	Success bool
	Data    dto.AuthResponseDTO
	Error   string
}

var ErrInvalidToken = errors.New("invalid token")
