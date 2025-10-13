package stores

import (
	"context"
	"goalify/internal/db"
	sqlcdb "goalify/internal/db/generated"
	"goalify/internal/testsetup"
	us "goalify/internal/users/stores"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const password = "Password123!"

var (
	userStore   us.UserStore
	gStore      GoalStore
	gcStore     GoalCategoryStore
	pgContainer *postgres.PostgresContainer
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
	userStore = us.NewUserStore(queries)
	gStore = NewGoalStore(queries)
	gcStore = NewGoalCategoryStore(queries)
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

func TestCreateGoal(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, err := gcStore.CreateGoalCategory(t.Name(), 50, user.Id)
	assert.NoError(t, err)

	goal, err := gStore.CreateGoal(t.Name(), "desc", user.Id, category.Id)
	assert.NoError(t, err)
	assert.Equal(t, t.Name(), goal.Title)
	assert.Equal(t, category.Id, goal.CategoryId)
	assert.Equal(t, user.Id, goal.UserId)
	assert.Equal(t, "desc", goal.Description)
}

func TestGetGoalById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.Id)
	goal, _ := gStore.CreateGoal(t.Name(), "desc", user.Id, category.Id)

	foundGoal, err := gStore.GetGoalById(goal.Id)
	assert.NoError(t, err)
	assert.Equal(t, goal.Id, foundGoal.Id)
	assert.Equal(t, goal.Title, foundGoal.Title)
	assert.Equal(t, goal.Description, foundGoal.Description)
	assert.Equal(t, goal.UserId, foundGoal.UserId)
	assert.Equal(t, goal.CategoryId, foundGoal.CategoryId)
}

func TestUpdateGoalById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.Id)
	goal, _ := gStore.CreateGoal(t.Name(), "desc", user.Id, category.Id)

	updates := map[string]any{
		"title":       "new title",
		"description": "new desc",
	}
	updated, err := gStore.UpdateGoalById(goal.Id, updates)
	assert.NoError(t, err)
	assert.Equal(t, "new title", updated.Title)
	assert.Equal(t, "new desc", updated.Description)
}

func TestDeleteGoalById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.Id)
	goal, _ := gStore.CreateGoal(t.Name(), "desc", user.Id, category.Id)

	_, err = gStore.GetGoalById(goal.Id)
	assert.NoError(t, err)

	err = gStore.DeleteGoalById(goal.Id)
	assert.NoError(t, err)

	_, err = gStore.GetGoalById(goal.Id)
	assert.Error(t, err)
}
