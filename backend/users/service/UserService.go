package service

import (
	"goalify/entities"

	"github.com/google/uuid"
)

type UserService interface {
	SignUp(email, password string) (*entities.UserDTO, error)
	Login(email, password string) (*entities.UserDTO, error)
	Refresh(email, refreshToken string) (*entities.UserDTO, error)
	DeleteUserById(id string) error
	UpdateUserById(id uuid.UUID, updates map[string]interface{}) (*entities.UserDTO, error)
}
