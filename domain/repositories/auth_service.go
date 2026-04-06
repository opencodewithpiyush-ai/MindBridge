package repositories

type IAuthService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (string, error)
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) bool
}
