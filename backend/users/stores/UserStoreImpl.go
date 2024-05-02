package stores

import (
	"goalify/models"

	"github.com/jmoiron/sqlx"
)

type UserStoreImpl struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStoreImpl {
	return &UserStoreImpl{db: db}
}

func (s *UserStoreImpl) CreateUser(email string, password string) error {
	_, err := s.db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, password)
	if err != nil {
		return err
	}
	return nil
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
