package stores

import (
	"goalify/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserStoreImpl struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStoreImpl {
	return &UserStoreImpl{db: db}
}

func (s *UserStoreImpl) CreateUser(email string, password string) (models.User, error) {
	_, err := s.db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, password)
	if err != nil {
		return models.User{}, err
	}
	var user models.User
	err = s.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserStoreImpl) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	// db.get returns an error if no rows are found
	err := s.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *UserStoreImpl) UpdateRefreshToken(id string, refreshToken string) (models.User, error) {
	_, err := s.db.Exec("UPDATE users SET refresh_token = $1 WHERE id = $2", refreshToken, id)
	if err != nil {
		return models.User{}, err
	}

	_, err = s.db.Exec("UPDATE users SET refresh_token_expiry = $1 WHERE id = $2", time.Now().Add(time.Hour*72), id)
	if err != nil {
		return models.User{}, err
	}

	var user models.User
	err = s.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *UserStoreImpl) GetUserById(id string) (models.User, error) {
	var user models.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
