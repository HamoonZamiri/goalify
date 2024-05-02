package stores

import "goalify/models"

type UserStore interface {
	CreateUser(email, password string) error
	GetUserByEmail(email string) (models.User, error)
}
