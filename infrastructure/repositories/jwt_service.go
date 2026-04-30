package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"mindbridge/config"
)

type JWTService struct {
	jwtSecret  []byte
	tokenTTL  time.Duration
	bcryptCost int
	redisClient *RedisClient
}

func NewJWTService(secret string, redisClient *RedisClient) *JWTService {
	ttl := time.Duration(config.TokenTTLHours) * time.Hour
	if ttl == 0 {
		ttl = 168 * time.Hour // default 7 days
	}
	cost := config.BcryptCost
	if cost == 0 {
		cost = 12
	}
	return &JWTService{
		jwtSecret:  []byte(secret),
		tokenTTL:  ttl,
		bcryptCost: cost,
		redisClient: redisClient,
	}
}

func (s *JWTService) GenerateToken(userID string) (string, error) {
	jti := uuid.New().String()
	claims := jwt.MapClaims{
		"user_id": userID,
		"jti":     jti,
		"exp":     time.Now().Add(s.tokenTTL).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	if s.redisClient != nil {
		s.redisClient.CreateSession(jti, userID, s.tokenTTL)
	}

	return signed, nil
}

func (s *JWTService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}
		jti, _ := claims["jti"].(string)
		if jti != "" && s.redisClient != nil {
			exists, err := s.redisClient.IsSessionValid(jti)
			if err != nil || !exists {
				return "", errors.New("token revoked or expired")
			}
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}

func (s *JWTService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.bcryptCost)
	return string(bytes), err
}

func (s *JWTService) ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
