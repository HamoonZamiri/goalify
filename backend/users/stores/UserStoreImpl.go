package stores

import (
	"goalify/entities"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserStoreImpl struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) *UserStoreImpl {
	return &UserStoreImpl{db: db}
}

func (s *UserStoreImpl) CreateUser(email string, password string) (entities.User, error) {
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING *`

	var user entities.User
	err := s.db.QueryRowx(query, email, password).StructScan(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserStoreImpl) GetUserByEmail(email string) (entities.User, error) {
	var user entities.User

	// db.get returns an error if no rows are found
	err := s.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *UserStoreImpl) UpdateRefreshToken(id string, refreshToken string) (entities.User, error) {
	query := `UPDATE users 
    SET refresh_token = $1,
    refresh_token_expiry = $2
    WHERE id = $3 
    RETURNING *`

	var user entities.User
	err := s.db.QueryRowx(query, refreshToken, time.Now().Add(time.Hour*72), id).StructScan(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserStoreImpl) GetUserById(id string) (entities.User, error) {
	var user entities.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return user, err
	}
	return user, nil
}
