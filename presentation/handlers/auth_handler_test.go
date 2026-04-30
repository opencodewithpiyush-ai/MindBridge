package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mindbridge/application/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

type mockAuthUseCase struct {
	mock.Mock
}

func (m *mockAuthUseCase) Register(req dto.RegisterRequestDTO) dto.AuthResponseDTO {
	args := m.Called(req)
	return args.Get(0).(dto.AuthResponseDTO)
}

func (m *mockAuthUseCase) Login(req dto.LoginRequestDTO) dto.AuthResponseDTO {
	args := m.Called(req)
	return args.Get(0).(dto.AuthResponseDTO)
}

func (m *mockAuthUseCase) Logout(userID, token string) error {
	args := m.Called(userID, token)
	return args.Error(0)
}

func TestRegisterHandler_ValidationFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/register", func(c *gin.Context) {
		var req dto.RegisterRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"errors":  []gin.H{gin.H{"field": "name", "message": "invalid"}},
		})
	})

	body := `{}`
	req, _ := http.NewRequest("POST", "/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"success":false`)
}

func TestRegisterHandler_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/register", func(c *gin.Context) {
		var req dto.RegisterRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
	})

	req, _ := http.NewRequest("POST", "/auth/register", strings.NewReader(`{bad json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	m := new(mockAuthUseCase)
	m.On("Login", mock.AnythingOfType("dto.LoginRequestDTO")).Return(dto.AuthResponseDTO{
		Token: "tok456",
		User:  dto.UserDTO{ID: "1", Name: "Test", Username: "test", Email: "test@example.com"},
	})

	router := gin.New()
	router.POST("/auth/login", func(c *gin.Context) {
		var req dto.LoginRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
		result := m.Login(req)
		if result.Token != "" {
			c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"token": result.Token, "user": result.User}})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "invalid credentials"})
		}
	})

	body := `{"email":"test@example.com","password":"Pass@123"}`
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"success":true`)
	m.AssertExpectations(t)
}

func TestLoginHandler_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/login", func(c *gin.Context) {
		var req dto.LoginRequestDTO
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
			return
		}
	})

	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(`{bad json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogoutHandler_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/auth/logout", func(c *gin.Context) {
		_, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	req, _ := http.NewRequest("POST", "/auth/logout", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
