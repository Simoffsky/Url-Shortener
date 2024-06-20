package auth

import "url-shorter/internal/models"

type UserRepository interface {
	CreateUser(*models.User) error
	Exists(login string) (bool, error)
	GetUser(login string) (*models.User, error)
}
