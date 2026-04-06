package usecases

import (
	"errors"
	"mindbridge/application/dto"
	domainRepo "mindbridge/domain/repositories"
	"mindbridge/infrastructure/repositories"
	"time"

	"mindbridge/domain/entities"
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

func (uc *AuthUseCase) Register(req dto.RegisterRequestDTO) RegisterResult {
	existingUser, err := uc.userRepo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return RegisterResult{
			Success: false,
			Error:   "email already registered",
		}
	}

	existingUserByUsername, err := uc.userRepo.FindByUsername(req.Username)
	if err == nil && existingUserByUsername != nil {
		return RegisterResult{
			Success: false,
			Error:   "username already taken",
		}
	}

	hashedPassword, err := uc.authService.HashPassword(req.Password)
	if err != nil {
		return RegisterResult{
			Success: false,
			Error:   "failed to hash password",
		}
	}

	user := &entities.User{
		Name:      req.Name,
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return RegisterResult{
			Success: false,
			Error:   "failed to create user",
		}
	}

	token, err := uc.authService.GenerateToken(user.ID)
	if err != nil {
		return RegisterResult{
			Success: false,
			Error:   "failed to generate token",
		}
	}

	if uc.redisClient != nil {
		expiry := 7 * 24 * time.Hour
		uc.redisClient.CreateSession(token, user.ID, expiry)
	}

	return RegisterResult{
		Success: true,
		Data: dto.AuthResponseDTO{
			Token: token,
			User: dto.UserDTO{
				ID:       user.ID,
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
			},
		},
	}
}

type LoginResult struct {
	Success bool
	Data    dto.AuthResponseDTO
	Error   string
}

func (uc *AuthUseCase) Login(req dto.LoginRequestDTO) LoginResult {
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return LoginResult{
			Success: false,
			Error:   "invalid email or password",
		}
	}

	if !uc.authService.ComparePassword(user.Password, req.Password) {
		return LoginResult{
			Success: false,
			Error:   "invalid email or password",
		}
	}

	token, err := uc.authService.GenerateToken(user.ID)
	if err != nil {
		return LoginResult{
			Success: false,
			Error:   "failed to generate token",
		}
	}

	if uc.redisClient != nil {
		expiry := 7 * 24 * time.Hour
		uc.redisClient.CreateSession(token, user.ID, expiry)
	}

	return LoginResult{
		Success: true,
		Data: dto.AuthResponseDTO{
			Token: token,
			User: dto.UserDTO{
				ID:       user.ID,
				Name:     user.Name,
				Username: user.Username,
				Email:    user.Email,
			},
		},
	}
}

var ErrInvalidToken = errors.New("invalid token")
