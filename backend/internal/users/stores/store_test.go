package stores

import (
	"context"
	"database/sql"
	"goalify/internal/db"
	"goalify/internal/testsetup"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	sqlcdb "goalify/internal/db/generated"
)

var (
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

	pgxPool, err := db.NewPgxPoolWithConnString(ctx, connStr)
	if err != nil {
		panic(err)
	}

	queries := sqlcdb.New(pgxPool)
	userStoreVar = NewUserStore(queries)
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	setup(ctx)

	code := m.Run()

	if err := pgContainer.Terminate(ctx); err != nil {
		log.Fatalf("Failed to terminate container: %s", err)
	}

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

	user, err = userStoreVar.GetUserByID(user.ID.String())
	assert.NoError(t, err)
	assert.Equal(t, "test2@mail.com", user.Email)
}

func TestGetUserDoesNotExist(t *testing.T) {
	t.Parallel()
	id := uuid.New()
	_, err := userStoreVar.GetUserByID(id.String())
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
	user, err = userStoreVar.UpdateRefreshToken(user.ID.String(), newToken.String())
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

	_, err = userStoreVar.GetUserByID(user.ID.String())
	assert.NoError(t, err)

	err = userStoreVar.DeleteUserByID(user.ID.String())
	assert.NoError(t, err)

	_, err = userStoreVar.GetUserByID(user.ID.String())
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func DeleteNonExistentUser(t *testing.T) {
	t.Parallel()
	err := userStoreVar.DeleteUserByID(uuid.New().String())
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	user, err := userStoreVar.CreateUser("test5@mail.com", "password")
	require.NoError(t, err)

	updates := map[string]any{
		"email": "test6@mail.com",
	}

	user, err = userStoreVar.UpdateUserByID(user.ID, updates)
	assert.NoError(t, err)
	assert.Equal(t, "test6@mail.com", user.Email)
}

func TestGetLevelById(t *testing.T) {
	t.Parallel()
	level := 1
	expectedXp := 100
	expectedCash := 100
	for ; level <= 100; level++ {
		levelObj, err := userStoreVar.GetLevelByID(level)
		assert.NoError(t, err)
		assert.Equal(t, level, levelObj.ID)
		assert.Equal(t, expectedXp, levelObj.LevelUpXp)
		assert.Equal(t, expectedCash, levelObj.CashReward)
		expectedXp += 10
		expectedCash += 10
	}
}

func TestLevelDoesNotExist(t *testing.T) {
	t.Parallel()
	_, err := userStoreVar.GetLevelByID(0)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	_, err = userStoreVar.GetLevelByID(101)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	_, err = userStoreVar.GetLevelByID(-1000)
	assert.ErrorIs(t, err, sql.ErrNoRows)

	_, err = userStoreVar.GetLevelByID(1000)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}
