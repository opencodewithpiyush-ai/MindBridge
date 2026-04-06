package repositories

import "mindbridge/domain/entities"

type IUserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByUsername(username string) (*entities.User, error)
	FindByID(id string) (*entities.User, error)
}
