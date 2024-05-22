package stores

import (
	"goalify/entities"

	"github.com/google/uuid"
)

type UserStore interface {
	CreateUser(email, password string) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	UpdateRefreshToken(id, refreshToken string) (*entities.User, error)
	GetUserById(id string) (*entities.User, error)
	DeleteUserById(id string) error
	UpdateUserById(id uuid.UUID, updates map[string]any) (*entities.User, error)
}
