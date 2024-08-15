package stores

import (
	"context"
	"goalify/db"
	us "goalify/users/stores"
	"log"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const password = "Password123!"

var (
	dbConn      *sqlx.DB
	userStore   us.UserStore
	goalStore   GoalStore
	gcStore     GoalCategoryStore
	pgContainer *postgres.PostgresContainer
)

func setup(ctx context.Context) {
	dbName := "user_store"
	dbUser := "postgres"
	dbPass := "password"

	var err error

	pgContainer, err = postgres.Run(ctx, "docker.io/postgres:16-alpine",
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

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	dbConn, err = db.NewWithConnString(connStr)
	if err != nil {
		panic(err)
	}

	userStore = us.NewUserStore(dbConn)
	goalStore = NewGoalStore(dbConn)
	gcStore = NewGoalCategoryStore(dbConn)

	err = goose.UpContext(ctx, dbConn.DB, "../../db/migrations")
	if err != nil {
		panic(err)
	}
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

func TestCreateCategory(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, err := gcStore.CreateGoalCategory(t.Name(), 50, user.Id)
	assert.NoError(t, err)
	assert.Equal(t, t.Name(), category.Title)
	assert.Equal(t, 50, category.Xp_per_goal)
	assert.Equal(t, user.Id, category.UserId)
}

func TestGetGoalCategoriesByUser(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)

	category1, _ := gcStore.CreateGoalCategory(t.Name()+"1", 50, user.Id)
	category2, _ := gcStore.CreateGoalCategory(t.Name()+"2", 50, user.Id)
	category3, _ := gcStore.CreateGoalCategory(t.Name()+"3", 50, user.Id)
	category4, _ := gcStore.CreateGoalCategory(t.Name()+"4", 50, user.Id)

	ids := []string{category1.Id.String(), category2.Id.String(), category3.Id.String(), category4.Id.String()}
	foundIds := 0

	categories, err := gcStore.GetGoalCategoriesByUserId(user.Id)

	for c := range categories {
		if slices.Contains(ids, categories[c].Id.String()) {
			foundIds++
		}
	}
	assert.NoError(t, err)
	assert.Equal(t, 4, len(categories))
	assert.Equal(t, 4, foundIds)
}

func TestGetGoalCategoryById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.Id)

	foundCategory, err := gcStore.GetGoalCategoryById(category.Id)
	assert.NoError(t, err)
	assert.Equal(t, category.Id, foundCategory.Id)
	assert.Equal(t, category.Title, foundCategory.Title)
	assert.Equal(t, category.Xp_per_goal, foundCategory.Xp_per_goal)
	assert.Equal(t, category.UserId, foundCategory.UserId)
}

func TestUpdateGoalCategoryById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.Id)

	updates := map[string]any{
		"title":       "new title",
		"xp_per_goal": 100,
	}
	updated, err := gcStore.UpdateGoalCategoryById(category.Id, updates)
	assert.NoError(t, err)
	assert.Equal(t, "new title", updated.Title)
	assert.Equal(t, 100, updated.Xp_per_goal)
}

func TestDeleteGoalCategoryById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.Id)

	_, err = gcStore.GetGoalCategoryById(category.Id)
	assert.NoError(t, err)

	err = gcStore.DeleteGoalCategoryById(category.Id)
	assert.NoError(t, err)

	_, err = gcStore.GetGoalCategoryById(category.Id)
	assert.Error(t, err)
}
