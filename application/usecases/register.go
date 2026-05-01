package usecases

import (
	"mindbridge/application/dto"
	"mindbridge/domain/entities"
	"net/mail"
	"strings"
	"time"
)

func (uc *AuthUseCase) Register(req dto.RegisterRequestDTO) RegisterResult {
	sanitizedEmail := strings.ToLower(strings.TrimSpace(req.Email))
	if _, err := mail.ParseAddress(sanitizedEmail); err != nil {
		return RegisterResult{Success: false, Error: "invalid email format"}
	}

	sanitizedUsername := strings.TrimSpace(req.Username)
	if sanitizedUsername == "" {
		return RegisterResult{Success: false, Error: "invalid username"}
	}

	existingUser, err := uc.userRepo.FindByEmail(sanitizedEmail)
	if err == nil && existingUser != nil {
		return RegisterResult{Success: false, Error: "email already registered"}
	}

	existingUserByUsername, err := uc.userRepo.FindByUsername(sanitizedUsername)
	if err == nil && existingUserByUsername != nil {
		return RegisterResult{Success: false, Error: "username already taken"}
	}

	hashedPassword, err := uc.authService.HashPassword(req.Password)
	if err != nil {
		return RegisterResult{Success: false, Error: "failed to hash password"}
	}

	user := &entities.User{
		Name:      req.Name,
		Username:  sanitizedUsername,
		Email:     sanitizedEmail,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.userRepo.Create(user); err != nil {
		return RegisterResult{Success: false, Error: "failed to create user"}
	}

	token, err := uc.authService.GenerateToken(user.ID)
	if err != nil {
		return RegisterResult{Success: false, Error: "failed to generate token"}
	}

	return RegisterResult{
		Success: true,
		Data: dto.AuthResponseDTO{
			Token: token,
			User:  dto.UserDTO{ID: user.ID, Name: user.Name, Username: user.Username, Email: user.Email},
		},
	}
}
