// Package stores is the repository layer package for the users domain
package stores

import (
	"context"
	"goalify/internal/entities"
	"time"

	sqlcdb "goalify/internal/db/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	UserStore interface {
		CreateUser(email, password string) (*entities.User, error)
		GetUserByEmail(email string) (*entities.User, error)
		UpdateRefreshToken(id, refreshToken string) (*entities.User, error)
		GetUserByID(id string) (*entities.User, error)
		DeleteUserByID(id string) error
		UpdateUserByID(id uuid.UUID, updates map[string]any) (*entities.User, error)

		GetLevelByID(id int) (*entities.Level, error)
	}
	userStore struct {
		queries *sqlcdb.Queries
	}
)

const DefaultLevel = 1

// Helper functions to convert between sqlc types and entity types
func pgxUserToEntity(u sqlcdb.User) *entities.User {
	return &entities.User{
		ID:                 uuid.UUID(u.ID.Bytes),
		Email:              u.Email,
		Password:           u.Password,
		Xp:                 int(u.Xp.Int32),
		LevelID:            int(u.LevelID.Int32),
		CashAvailable:      int(u.CashAvailable.Int32),
		RefreshToken:       uuid.UUID(u.RefreshToken.Bytes),
		RefreshTokenExpiry: u.RefreshTokenExpiry.Time,
		CreatedAt:          u.CreatedAt.Time,
		UpdatedAt:          u.UpdatedAt.Time,
	}
}

func pgxLevelToEntity(l sqlcdb.Level) *entities.Level {
	return &entities.Level{
		ID:         int(l.ID),
		LevelUpXp:  int(l.LevelUpXp),
		CashReward: int(l.CashReward),
		CreatedAt:  l.CreatedAt.Time,
		UpdatedAt:  l.UpdatedAt.Time,
	}
}

func NewUserStore(queries *sqlcdb.Queries) UserStore {
	return &userStore{
		queries: queries,
	}
}

func (s *userStore) CreateUser(email string, password string) (*entities.User, error) {
	expiry := time.Now().Add(time.Hour * 72)
	params := sqlcdb.CreateUserParams{
		Email:              email,
		Password:           password,
		RefreshTokenExpiry: pgtype.Timestamp{Time: expiry, Valid: true},
		LevelID:            pgtype.Int4{Int32: DefaultLevel, Valid: true},
	}

	user, err := s.queries.CreateUser(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxUserToEntity(user), nil
}

func (s *userStore) GetUserByEmail(email string) (*entities.User, error) {
	user, err := s.queries.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return pgxUserToEntity(user), nil
}

func (s *userStore) UpdateRefreshToken(id string, refreshToken string) (*entities.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	refreshTokenUUID, err := uuid.Parse(refreshToken)
	if err != nil {
		return nil, err
	}

	params := sqlcdb.UpdateRefreshTokenParams{
		RefreshToken:       pgtype.UUID{Bytes: refreshTokenUUID, Valid: true},
		RefreshTokenExpiry: pgtype.Timestamp{Time: time.Now().Add(time.Hour * 72), Valid: true},
		ID:                 pgtype.UUID{Bytes: userID, Valid: true},
	}

	user, err := s.queries.UpdateRefreshToken(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxUserToEntity(user), nil
}

func (s *userStore) GetUserByID(id string) (*entities.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, err := s.queries.GetUserById(
		context.Background(),
		pgtype.UUID{Bytes: userID, Valid: true},
	)
	if err != nil {
		return nil, err
	}

	return pgxUserToEntity(user), nil
}

func (s *userStore) DeleteUserByID(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return s.queries.DeleteUserById(context.Background(), pgtype.UUID{Bytes: userID, Valid: true})
}

func (s *userStore) UpdateUserByID(id uuid.UUID, updates map[string]any) (*entities.User, error) {
	params := sqlcdb.UpdateUserByIdParams{
		ID: pgtype.UUID{Bytes: id, Valid: true},
	}

	// Convert map updates to typed parameters
	if email, ok := updates["email"]; ok {
		if emailStr, ok := email.(string); ok {
			params.Email = pgtype.Text{String: emailStr, Valid: true}
		}
	}
	if password, ok := updates["password"]; ok {
		if passwordStr, ok := password.(string); ok {
			params.Password = pgtype.Text{String: passwordStr, Valid: true}
		}
	}
	if refreshToken, ok := updates["refresh_token"]; ok {
		if refreshTokenUUID, ok := refreshToken.(uuid.UUID); ok {
			params.RefreshToken = pgtype.UUID{Bytes: refreshTokenUUID, Valid: true}
		}
	}
	if refreshTokenExpiry, ok := updates["refresh_token_expiry"]; ok {
		if expiryTime, ok := refreshTokenExpiry.(time.Time); ok {
			params.RefreshTokenExpiry = pgtype.Timestamp{Time: expiryTime, Valid: true}
		}
	}
	if levelID, ok := updates["level_id"]; ok {
		if levelInt, ok := levelID.(int); ok {
			params.LevelID = pgtype.Int4{Int32: int32(levelInt), Valid: true}
		}
	}

	if xp, ok := updates["xp"]; ok {
		if xpInt, ok := xp.(int); ok {
			params.Xp = pgtype.Int4{Int32: int32(xpInt), Valid: true}
		}
	}

	if cash, ok := updates["cash_available"]; ok {
		if cashInt, ok := cash.(int); ok {
			params.CashAvailable = pgtype.Int4{Int32: int32(cashInt), Valid: true}
		}
	}

	user, err := s.queries.UpdateUserById(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxUserToEntity(user), nil
}

func (s *userStore) GetLevelByID(id int) (*entities.Level, error) {
	level, err := s.queries.GetLevelById(context.Background(), int32(id))
	if err != nil {
		return nil, err
	}

	return pgxLevelToEntity(level), nil
}
