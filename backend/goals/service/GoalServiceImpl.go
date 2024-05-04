package service

import (
	"errors"
	"fmt"
	"goalify/entities"
	"goalify/goals/stores"
	"strings"

	"github.com/google/uuid"
)

const (
	XP_PER_GOAL_MAX = 100
)

type GoalServiceImpl struct {
	goalStore         stores.GoalStore
	goalCategoryStore stores.GoalCategoryStore
}

func NewGoalService(goalStore stores.GoalStore, goalCategoryStore stores.GoalCategoryStore) *GoalServiceImpl {
	return &GoalServiceImpl{
		goalStore:         goalStore,
		goalCategoryStore: goalCategoryStore,
	}
}

func isInvalidUUID(id uuid.UUID) bool {
	return id == uuid.Nil || id.String() == ""
}

func isEmptyField(field string) bool {
	return field == ""
}

func (gs *GoalServiceImpl) CreateGoal(title, description string, userId, categoryId uuid.UUID) (entities.Goal, error) {
	if isEmptyField(title) {
		return entities.Goal{}, errors.New("title cannot be empty")
	}
	if isEmptyField(description) {
		return entities.Goal{}, errors.New("description cannot be empty")
	}

	if isInvalidUUID(categoryId) {
		return entities.Goal{}, errors.New("category id cannot be empty")
	}

	if isInvalidUUID(userId) {
		return entities.Goal{}, errors.New("invalid uuid")
	}

	return gs.goalStore.CreateGoal(title, description, userId, categoryId)
}

func (gs *GoalServiceImpl) UpdateGoalStatus(userId, goalId uuid.UUID, status string) (entities.Goal, error) {
	cleanStatus := strings.ToLower(status)

	if cleanStatus != "completed" && cleanStatus != "not_completed" {
		return entities.Goal{}, errors.New("status must be either 'completed' or 'not_completed'")
	}

	goal, err := gs.goalStore.GetGoalById(goalId)
	if err != nil {
		return entities.Goal{}, err
	}
	if goal.UserId != userId {
		return entities.Goal{}, errors.New("user does not own this goal")
	}

	return gs.goalStore.UpdateGoalStatus(goalId, status)
}

func (gs *GoalServiceImpl) GetGoalsByUserId(userId uuid.UUID) ([]entities.Goal, error) {
	return gs.goalStore.GetGoalsByUserId(userId)
}

func (gs *GoalServiceImpl) UpdateGoalTitle(title string, userId, goalId uuid.UUID) (entities.Goal, error) {
	if title == "" {
		return entities.Goal{}, errors.New("title cannot be empty")
	}
	goal, err := gs.goalStore.GetGoalById(goalId)
	if err != nil {
		return entities.Goal{}, err
	}
	if goal.UserId != userId {
		return entities.Goal{}, errors.New("user does not own this goal")
	}

	return gs.goalStore.UpdateGoalTitle(title, goalId)
}

func (gs *GoalServiceImpl) UpdateGoalDescription(description string, userId, goalId uuid.UUID) (entities.Goal, error) {
	if isEmptyField(description) {
		return entities.Goal{}, errors.New("description cannot be empty")
	}

	goal, err := gs.goalStore.GetGoalById(goalId)
	if err != nil {
		return entities.Goal{}, err
	}
	if goal.UserId != userId {
		return entities.Goal{}, errors.New("user does not own this goal")
	}
	return gs.goalStore.UpdateGoalDescription(description, goalId)
}

func (gs *GoalServiceImpl) GetGoalById(goalId uuid.UUID) (entities.Goal, error) {
	return gs.goalStore.GetGoalById(goalId)
}

func (gs *GoalServiceImpl) CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (entities.GoalCategory, error) {
	if isEmptyField(title) {
		return entities.GoalCategory{}, errors.New("title cannot be empty")
	}
	if xpPerGoal <= 0 || xpPerGoal > XP_PER_GOAL_MAX {
		return entities.GoalCategory{}, fmt.Errorf("xp per goal must be between 1 and %d", XP_PER_GOAL_MAX)
	}

	if isInvalidUUID(userId) {
		return entities.GoalCategory{}, errors.New("invalid uuid")
	}
	return gs.goalCategoryStore.CreateGoalCategory(title, xpPerGoal, userId)
}

func (gs *GoalServiceImpl) GetGoalCategoriesByUserId(userId uuid.UUID) ([]entities.GoalCategory, error) {
	if isInvalidUUID(userId) {
		return []entities.GoalCategory{}, errors.New("invalid uuid")
	}
	return gs.goalCategoryStore.GetGoalCategoriesByUserId(userId)
}

func (gs *GoalServiceImpl) UpdateGoalCategory(userId, categoryId, goalId uuid.UUID) (entities.GoalCategory, error) {
	if isInvalidUUID(categoryId) {
		return entities.GoalCategory{}, errors.New("invalid uuid category id")
	}
	if isInvalidUUID(goalId) {
		return entities.GoalCategory{}, errors.New("invalid uuid goal id")
	}
	gc, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if err != nil {
		return entities.GoalCategory{}, err
	}

	if gc.UserId != userId {
		return entities.GoalCategory{}, errors.New("user does not own this category")
	}

	return gs.goalCategoryStore.UpdateGoalCategory(categoryId, goalId)
}

func (gs *GoalServiceImpl) GetGoalCategoryById(categoryId uuid.UUID) (entities.GoalCategory, error) {
	if isInvalidUUID(categoryId) {
		return entities.GoalCategory{}, errors.New("invalid uuid category id")
	}
	return gs.goalCategoryStore.GetGoalCategoryById(categoryId)
}
