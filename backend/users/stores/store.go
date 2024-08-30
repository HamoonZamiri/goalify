package stores

import (
	"goalify/db"
	"goalify/entities"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserStore interface {
	CreateUser(email, password string) (*entities.User, error)
	GetUserByEmail(email string) (*entities.User, error)
	UpdateRefreshToken(id, refreshToken string) (*entities.User, error)
	GetUserById(id string) (*entities.User, error)
	DeleteUserById(id string) error
	UpdateUserById(id uuid.UUID, updates map[string]any) (*entities.User, error)

	GetLevelById(id int) (*entities.Level, error)
}

const DEFAULT_LEVEL = 1

type userStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) UserStore {
	return &userStore{db: db}
}

func (s *userStore) CreateUser(email string, password string) (*entities.User, error) {
	query := `INSERT INTO users (email, password, refresh_token_expiry, level_id) VALUES ($1, $2, $3, $4) RETURNING *`

	var user entities.User
	expiry := time.Now().Add(time.Hour * 72)
	err := s.db.QueryRowx(query, email, password, expiry, DEFAULT_LEVEL).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userStore) GetUserByEmail(email string) (*entities.User, error) {
	var user entities.User

	// db.get returns an error if no rows are found
	err := s.db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userStore) UpdateRefreshToken(id string, refreshToken string) (*entities.User, error) {
	query := `UPDATE users 
    SET refresh_token = $1,
    refresh_token_expiry = $2
    WHERE id = $3 
    RETURNING *`

	var user entities.User
	err := s.db.QueryRowx(query, refreshToken, time.Now().Add(time.Hour*72), id).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userStore) GetUserById(id string) (*entities.User, error) {
	var user entities.User
	err := s.db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userStore) DeleteUserById(id string) error {
	_, err := s.db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

func (s *userStore) UpdateUserById(id uuid.UUID, updates map[string]any) (*entities.User, error) {
	query, args := db.BuildUpdateQuery("users", updates, id)
	query = query + " RETURNING *"

	var user entities.User
	err := s.db.QueryRowx(query, args...).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userStore) GetLevelById(id int) (*entities.Level, error) {
	var level entities.Level
	err := s.db.Get(&level, "SELECT * FROM levels WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &level, nil
}
