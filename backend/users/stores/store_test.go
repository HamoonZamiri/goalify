package stores

import (
	"context"
	"database/sql"
	"goalify/db"
	"goalify/testsetup"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	dbConn       *sqlx.DB
	userStoreVar UserStore
	pgContainer  *postgres.PostgresContainer
)

func setup(ctx context.Context) {
	var err error

	pgContainer, err = testsetup.GetPgContainer()
	if err != nil {
		panic(err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	dbConn, err = db.NewWithConnString(connStr)
	if err != nil {
		panic(err)
	}

	userStoreVar = NewUserStore(dbConn)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	setup(ctx)

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("Failed to terminate container: %s", err)
		}
	}()
	code := m.Run()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	t.Parallel()
	user, err := userStoreVar.CreateUser("test1@mail.com", "password")
	assert.NoError(t, err)
	assert.Equal(t, "test1@mail.com", user.Email)
}

func TestGetUserById(t *testing.T) {
	t.Parallel()
	user, err := userStoreVar.CreateUser("test2@mail.com", "password")
	require.NoError(t, err)

	user, err = userStoreVar.GetUserById(user.Id.String())
	assert.NoError(t, err)
	assert.Equal(t, "test2@mail.com", user.Email)
}

func TestGetUserDoesNotExist(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	_, err := userStoreVar.GetUserById(id.String())
	assert.Error(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestUpdateRefreshToken(t *testing.T) {
	t.Parallel()
	user, err := userStoreVar.CreateUser("test3@mail.com", "password")
	require.NoError(t, err)
	oldExpiry := user.RefreshTokenExpiry
	oldToken := user.RefreshToken

	newToken := uuid.New()
	user, err = userStoreVar.UpdateRefreshToken(user.Id.String(), newToken.String())
	assert.NoError(t, err)
	assert.NotEqual(t, oldExpiry, user.RefreshTokenExpiry)
	assert.NotEqual(t, oldToken.String(), user.RefreshToken.String())
}

func TestIncorrectRefreshToken(t *testing.T) {
	t.Parallel()
	_, err := userStoreVar.UpdateRefreshToken(uuid.New().String(), uuid.New().String())
	assert.Error(t, err)
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	user, err := userStoreVar.CreateUser("test4@mail.com", "password")
	require.NoError(t, err)

	_, err = userStoreVar.GetUserById(user.Id.String())
	assert.NoError(t, err)

	err = userStoreVar.DeleteUserById(user.Id.String())
	assert.NoError(t, err)

	_, err = userStoreVar.GetUserById(user.Id.String())
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func DeleteNonExistentUser(t *testing.T) {
	t.Parallel()
	err := userStoreVar.DeleteUserById(uuid.New().String())
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	user, err := userStoreVar.CreateUser("test5@mail.com", "password")
	require.NoError(t, err)

	updates := map[string]any{
		"email": "test6@mail.com",
	}

	user, err = userStoreVar.UpdateUserById(user.Id, updates)
	assert.NoError(t, err)
	assert.Equal(t, "test6@mail.com", user.Email)
}

func TestGetLevelById(t *testing.T) {
	t.Parallel()
	level := 1
	expectedXp := 100
	expectedCash := 100
	for ; level <= 100; level++ {
		levelObj, err := userStoreVar.GetLevelById(level)
		assert.NoError(t, err)
		assert.Equal(t, level, levelObj.Id)
		assert.Equal(t, expectedXp, levelObj.LevelUpXp)
		assert.Equal(t, expectedCash, levelObj.CashReward)
		expectedXp += 10
		expectedCash += 10
	}
}

func TestLevelDoesNotExist(t *testing.T) {
	t.Parallel()
	_, err := userStoreVar.GetLevelById(0)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	_, err = userStoreVar.GetLevelById(101)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	_, err = userStoreVar.GetLevelById(-1000)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	_, err = userStoreVar.GetLevelById(1000)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
