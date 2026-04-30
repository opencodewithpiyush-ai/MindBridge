package usecases_test

import (
	"testing"

	"mindbridge/application/dto"
	"mindbridge/application/usecases"
	"mindbridge/application/usecases/mocks"
	"mindbridge/domain/entities"

	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func TestRegister_Success(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	mockAuth := new(mocks.IAuthService)
	uc := usecases.NewAuthUseCase(mockRepo, mockAuth, nil)

	mockRepo.On("FindByEmail", "test@example.com").Return(nil, nil)
	mockRepo.On("FindByUsername", "testuser").Return(nil, nil)
	mockAuth.On("HashPassword", "password123").Return("hashed", nil)
	mockRepo.On("Create", mock.AnythingOfType("*entities.User")).Return(nil)
	mockAuth.On("GenerateToken", mock.AnythingOfType("string")).Return("token123", nil)

	result := uc.Register(dto.RegisterRequestDTO{Name: "Test", Username: "testuser", Email: "test@example.com", Password: "password123"})
	assert.True(t, result.Success)
	assert.Equal(t, "token123", result.Data.Token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(mocks.IUserRepository)
	mockAuth := new(mocks.IAuthService)
	uc := usecases.NewAuthUseCase(mockRepo, mockAuth, nil)

	user := &entities.User{ID: "123", Email: "test@example.com", Password: "hashed"}
	mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)
	mockAuth.On("ComparePassword", "hashed", "password123").Return(true)
	mockAuth.On("GenerateToken", "123").Return("token123", nil)

	result := uc.Login(dto.LoginRequestDTO{Email: "test@example.com", Password: "password123"})
	assert.True(t, result.Success)
	assert.Equal(t, "token123", result.Data.Token)
}

func TestLogout_NoRedis(t *testing.T) {
	uc := usecases.NewAuthUseCase(nil, nil, nil)
	err := uc.Logout("user123", "token123")
	assert.NoError(t, err)
}
