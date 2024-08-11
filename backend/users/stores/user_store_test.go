package stores

import (
	"context"
	"database/sql"
	"goalify/db"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	dbConn    *sqlx.DB
	userStore UserStore
)

const MIGRATION_STR string = `
CREATE TABLE levels  (
    id SERIAL PRIMARY KEY,
    level_up_xp INTEGER NOT NULL,
    cash_reward INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE users  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    xp INTEGER DEFAULT 0,
    level_id SERIAL REFERENCES levels(id),
    cash_available INTEGER DEFAULT 0,
    refresh_token UUID DEFAULT gen_random_uuid(),
    refresh_token_expiry TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE goal_categories  (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    xp_per_goal INTEGER NOT NULL,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TYPE goal_status AS ENUM ('complete', 'not_complete');

CREATE TABLE goals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) DEFAULT '',
    user_id UUID REFERENCES users(id),
    category_id UUID REFERENCES goal_categories(id) ON DELETE CASCADE,
    status goal_status DEFAULT 'not_complete',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);  

-- Insert Default Levels
INSERT INTO levels (id, level_up_xp, cash_reward) VALUES (1, 100, 10);`

func setup() {
	dbConn.MustExec(MIGRATION_STR)
	userStore = NewUserStore(dbConn)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	dbName := "user_store"
	dbUser := "postgres"
	dbPass := "password"

	pgContainer, err := postgres.Run(ctx, "docker.io/postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	dbConn, err = db.NewWithConnString(connStr)
	if err != nil {
		panic(err)
	}

	setup()

	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser("test1@mail.com", "password")
	assert.NoError(t, err)
	assert.Equal(t, "test1@mail.com", user.Email)
}

func TestGetUserById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser("test2@mail.com", "password")
	require.NoError(t, err)

	user, err = userStore.GetUserById(user.Id.String())
	assert.NoError(t, err)
	assert.Equal(t, "test2@mail.com", user.Email)
}

func TestGetUserDoesNotExist(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	_, err := userStore.GetUserById(id.String())
	assert.Error(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestUpdateRefreshToken(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser("test3@mail.com", "password")
	require.NoError(t, err)
	oldExpiry := user.RefreshTokenExpiry
	oldToken := user.RefreshToken

	newToken := uuid.New()
	user, err = userStore.UpdateRefreshToken(user.Id.String(), newToken.String())
	assert.NoError(t, err)
	assert.NotEqual(t, oldExpiry, user.RefreshTokenExpiry)
	assert.Greater(t, user.RefreshTokenExpiry.Unix(), time.Now().Unix())
	assert.NotEqual(t, oldToken.String(), user.RefreshToken.String())
}

func TestIncorrectRefreshToken(t *testing.T) {
	t.Parallel()
	_, err := userStore.UpdateRefreshToken(uuid.New().String(), uuid.New().String())
	assert.Error(t, err)
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser("test4@mail.com", "password")
	require.NoError(t, err)

	_, err = userStore.GetUserById(user.Id.String())
	assert.NoError(t, err)

	err = userStore.DeleteUserById(user.Id.String())
	assert.NoError(t, err)

	_, err = userStore.GetUserById(user.Id.String())
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func DeleteNonExistentUser(t *testing.T) {
	t.Parallel()
	err := userStore.DeleteUserById(uuid.New().String())
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser("test5@mail.com", "password")
	require.NoError(t, err)

	updates := map[string]any{
		"email": "test6@mail.com",
	}

	user, err = userStore.UpdateUserById(user.Id, updates)
	assert.NoError(t, err)
	assert.Equal(t, "test6@mail.com", user.Email)
}
