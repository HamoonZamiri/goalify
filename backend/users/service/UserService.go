package service

import "goalify/entities"

type UserService interface {
	SignUp(email, password string) (*entities.UserDTO, error)
	Login(email, password string) (*entities.UserDTO, error)
	Refresh(email, refreshToken string) (*entities.UserDTO, error)
	DeleteUserById(id string) error
}
