package stores

import "goalify/models"

type UserStore interface {
	CreateUser(email, password string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	UpdateRefreshToken(id, refreshToken string) (models.User, error)
	GetUserById(id string) (models.User, error)
}
