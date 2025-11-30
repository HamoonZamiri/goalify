package stores

import (
	"context"
	"fmt"
	"goalify/internal/db"
	"goalify/internal/entities"
	"goalify/internal/testsetup"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	sqlcdb "goalify/internal/db/generated"

	us "goalify/internal/users/stores"
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
	category, err := gcStore.CreateGoalCategory(t.Name(), 50, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, t.Name(), category.Title)
	assert.Equal(t, 50, category.XPPerGoal)
	assert.Equal(t, user.ID, category.UserID)
}

func TestGetGoalCategoriesByUser(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)

	category1, _ := gcStore.CreateGoalCategory(t.Name()+"1", 50, user.ID)
	category2, _ := gcStore.CreateGoalCategory(t.Name()+"2", 50, user.ID)
	category3, _ := gcStore.CreateGoalCategory(t.Name()+"3", 50, user.ID)
	category4, _ := gcStore.CreateGoalCategory(t.Name()+"4", 50, user.ID)

	ids := []string{
		category1.ID.String(),
		category2.ID.String(),
		category3.ID.String(),
		category4.ID.String(),
	}
	foundIds := 0

	categories, err := gcStore.GetGoalCategoriesByUserID(user.ID)

	for c := range categories {
		if slices.Contains(ids, categories[c].ID.String()) {
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
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.ID)

	foundCategory, err := gcStore.GetGoalCategoryByID(category.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, category.ID, foundCategory.ID)
	assert.Equal(t, category.Title, foundCategory.Title)
	assert.Equal(t, category.XPPerGoal, foundCategory.XPPerGoal)
	assert.Equal(t, category.UserID, foundCategory.UserID)
}

func TestUpdateGoalCategoryById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.ID)

	updates := map[string]any{
		"title":       "new title",
		"xp_per_goal": 100,
	}
	updated, err := gcStore.UpdateGoalCategoryByID(category.ID, user.ID, updates)
	assert.NoError(t, err)
	assert.Equal(t, "new title", updated.Title)
	assert.Equal(t, 100, updated.XPPerGoal)
}

func TestDeleteGoalCategoryById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.ID)

	_, err = gcStore.GetGoalCategoryByID(category.ID, user.ID)
	assert.NoError(t, err)

	err = gcStore.DeleteGoalCategoryByID(category.ID, user.ID)
	assert.NoError(t, err)

	_, err = gcStore.GetGoalCategoryByID(category.ID, user.ID)
	assert.Error(t, err)
}

func TestCreateGoal(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, err := gcStore.CreateGoalCategory(t.Name(), 50, user.ID)
	assert.NoError(t, err)

	goal, err := gStore.CreateGoal(t.Name(), "desc", user.ID, category.ID)
	assert.NoError(t, err)
	assert.Equal(t, t.Name(), goal.Title)
	assert.Equal(t, category.ID, goal.CategoryID)
	assert.Equal(t, user.ID, goal.UserID)
	assert.Equal(t, "desc", goal.Description)
}

func TestGetGoalById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.ID)
	goal, _ := gStore.CreateGoal(t.Name(), "desc", user.ID, category.ID)

	foundGoal, err := gStore.GetGoalByID(goal.ID, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, goal.ID, foundGoal.ID)
	assert.Equal(t, goal.Title, foundGoal.Title)
	assert.Equal(t, goal.Description, foundGoal.Description)
	assert.Equal(t, goal.UserID, foundGoal.UserID)
	assert.Equal(t, goal.CategoryID, foundGoal.CategoryID)
}

func TestUpdateGoalById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.ID)
	goal, _ := gStore.CreateGoal(t.Name(), "desc", user.ID, category.ID)

	updates := map[string]any{
		"title":       "new title",
		"description": "new desc",
	}
	updated, err := gStore.UpdateGoalByID(goal.ID, user.ID, updates)
	assert.NoError(t, err)
	assert.Equal(t, "new title", updated.Title)
	assert.Equal(t, "new desc", updated.Description)
}

func TestDeleteGoalById(t *testing.T) {
	t.Parallel()
	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)
	category, _ := gcStore.CreateGoalCategory(t.Name(), 50, user.ID)
	goal, _ := gStore.CreateGoal(t.Name(), "desc", user.ID, category.ID)

	_, err = gStore.GetGoalByID(goal.ID, user.ID)
	assert.NoError(t, err)

	err = gStore.DeleteGoalByID(goal.ID, user.ID)
	assert.NoError(t, err)

	_, err = gStore.GetGoalByID(goal.ID, user.ID)
	assert.Error(t, err)
}

func TestResetGoalsByCategoryID(t *testing.T) {
	t.Parallel()

	user, err := userStore.CreateUser(t.Name()+"@mail.com", password)
	assert.NoError(t, err)

	category, err := gcStore.CreateGoalCategory(t.Name(), 50, user.ID)
	assert.NoError(t, err)

	// Create multiple goals for this category
	numGoals := 5
	goals := make([]*entities.Goal, numGoals)
	for i := range numGoals {
		var goal *entities.Goal
		goal, err = gStore.CreateGoal(
			fmt.Sprintf("%s-goal-%d", t.Name(), i),
			"desc",
			user.ID,
			category.ID,
		)
		assert.NoError(t, err)
		goals[i] = goal
	}

	// Mark all goals as complete
	for i, goal := range goals {
		var updatedGoal *entities.Goal
		updatedGoal, err = gStore.UpdateGoalByID(
			goal.ID,
			user.ID,
			map[string]any{"status": "complete"},
		)
		assert.NoError(t, err)
		assert.Equal(t, "complete", updatedGoal.Status)
		goals[i] = updatedGoal // Update slice with latest state
	}

	// Reset all goals in the category
	err = gStore.ResetGoalsByCategoryID(category.ID, user.ID)
	assert.NoError(t, err)

	// Verify all goals are now not_complete
	for _, goal := range goals {
		var fetchedGoal *entities.Goal
		fetchedGoal, err = gStore.GetGoalByID(goal.ID, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "not_complete", fetchedGoal.Status,
			"Goal %s should be not_complete after reset", goal.ID)
	}

	// === Test ownership validation ===
	// Mark all goals complete again
	for i, goal := range goals {
		var updatedGoal *entities.Goal
		updatedGoal, err = gStore.UpdateGoalByID(
			goal.ID,
			user.ID,
			map[string]any{"status": "complete"},
		)
		assert.NoError(t, err)
		assert.Equal(t, "complete", updatedGoal.Status)
		goals[i] = updatedGoal
	}

	// Create another user who doesn't own this category
	user2, err := userStore.CreateUser(t.Name()+"1@mail.com", password)
	assert.NoError(t, err)

	// Attempt reset as unauthorized user (should not affect goals)
	err = gStore.ResetGoalsByCategoryID(category.ID, user2.ID)
	assert.NoError(t, err) // No error, but should be no-op

	// Verify all goals are STILL complete (not reset by unauthorized user)
	for _, goal := range goals {
		fetchedGoal, err := gStore.GetGoalByID(goal.ID, user.ID)
		assert.NoError(t, err)
		assert.Equal(t, "complete", fetchedGoal.Status,
			"Goal %s should still be complete (unauthorized reset attempt)", goal.ID)
	}
}
