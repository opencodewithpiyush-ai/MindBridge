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

func (uc *AuthUseCase) GetProfile(userID string) (dto.UserDTO, error) {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil || user == nil {
		return dto.UserDTO{}, errors.New("user not found")
	}
	return dto.UserDTO{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
