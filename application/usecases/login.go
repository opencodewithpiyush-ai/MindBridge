package usecases

import (
	"mindbridge/application/dto"
	"time"
)

func (uc *AuthUseCase) Login(req dto.LoginRequestDTO) LoginResult {
	user, err := uc.userRepo.FindByEmail(req.Email)
	if err != nil || user == nil {
		return LoginResult{Success: false, Error: "invalid email or password"}
	}

	if !uc.authService.ComparePassword(user.Password, req.Password) {
		return LoginResult{Success: false, Error: "invalid email or password"}
	}

	token, err := uc.authService.GenerateToken(user.ID)
	if err != nil {
		return LoginResult{Success: false, Error: "failed to generate token"}
	}

	if uc.redisClient != nil {
		expiry := 7 * 24 * time.Hour
		uc.redisClient.CreateSession(token, user.ID, expiry)
	}

	return LoginResult{
		Success: true,
		Data: dto.AuthResponseDTO{
			Token: token,
			User:  dto.UserDTO{ID: user.ID, Name: user.Name, Username: user.Username, Email: user.Email},
		},
	}
}
