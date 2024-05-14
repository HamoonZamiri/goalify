package service

import (
	"database/sql"
	"errors"
	"fmt"
	"goalify/entities"
	"goalify/goals/stores"
	"goalify/svcerror"
	"log/slog"
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

func (gs *GoalServiceImpl) CreateGoal(title, description string, userId, categoryId uuid.UUID) (*entities.Goal, error) {
	if isEmptyField(title) {
		return nil, errors.New("title cannot be empty")
	}
	if isEmptyField(description) {
		return nil, errors.New("description cannot be empty")
	}

	if isInvalidUUID(categoryId) {
		return nil, errors.New("category id cannot be empty")
	}

	if isInvalidUUID(userId) {
		return nil, errors.New("invalid uuid")
	}

	createdGoal, err := gs.goalStore.CreateGoal(title, description, userId, categoryId)
	if err != nil {
		slog.Error("error creating goal", "err", err)
		return nil, fmt.Errorf("%w: error creating goal", svcerror.ErrInternalServer)
	}
	return createdGoal, nil
}

func (gs *GoalServiceImpl) UpdateGoalStatus(status string, goalId, userId uuid.UUID) (*entities.Goal, error) {
	cleanStatus := strings.ToLower(status)

	if cleanStatus != "completed" && cleanStatus != "not_completed" {
		return nil, errors.New("status must be either 'completed' or 'not_completed'")
	}

	goal, err := gs.goalStore.GetGoalById(goalId)
	if err != nil {
		return nil, err
	}
	if goal.UserId != userId {
		return nil, errors.New("user does not own this goal")
	}

	return gs.goalStore.UpdateGoalStatus(goalId, status)
}

func (gs *GoalServiceImpl) GetGoalsByUserId(userId uuid.UUID) ([]*entities.Goal, error) {
	goals, err := gs.goalStore.GetGoalsByUserId(userId)
	if err == sql.ErrNoRows {
		return goals, nil
	}
	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (gs *GoalServiceImpl) UpdateGoalTitle(title string, goalId, userId uuid.UUID) (*entities.Goal, error) {
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	goal, err := gs.goalStore.GetGoalById(goalId)
	if err != nil {
		return nil, err
	}
	if goal.UserId != userId {
		return nil, errors.New("user does not own this goal")
	}

	return gs.goalStore.UpdateGoalTitle(title, goalId)
}

func (gs *GoalServiceImpl) UpdateGoalDescription(description string, goalId, userId uuid.UUID) (*entities.Goal, error) {
	if isEmptyField(description) {
		return nil, errors.New("description cannot be empty")
	}

	goal, err := gs.goalStore.GetGoalById(goalId)
	if err != nil {
		return nil, err
	}
	if goal.UserId != userId {
		return nil, errors.New("user does not own this goal")
	}
	return gs.goalStore.UpdateGoalDescription(description, goalId)
}

func (gs *GoalServiceImpl) GetGoalById(goalId uuid.UUID) (*entities.Goal, error) {
	return gs.goalStore.GetGoalById(goalId)
}

func (gs *GoalServiceImpl) CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (*entities.GoalCategory, error) {
	if isEmptyField(title) {
		return nil, errors.New("title cannot be empty")
	}
	if xpPerGoal <= 0 || xpPerGoal > XP_PER_GOAL_MAX {
		return nil, fmt.Errorf("xp per goal must be between 1 and %d", XP_PER_GOAL_MAX)
	}

	if isInvalidUUID(userId) {
		return nil, errors.New("invalid uuid")
	}
	cat, err := gs.goalCategoryStore.CreateGoalCategory(title, xpPerGoal, userId)
	if err != nil {
		slog.Error("error creating goal category", "err", err)
		return nil, fmt.Errorf("%w: error creating goal category", svcerror.ErrInternalServer)
	}
	return cat, nil
}

func (gs *GoalServiceImpl) GetGoalCategoriesByUserId(userId uuid.UUID) ([]*entities.GoalCategory, error) {
	if isInvalidUUID(userId) {
		return []*entities.GoalCategory{}, errors.New("invalid uuid")
	}
	return gs.goalCategoryStore.GetGoalCategoriesByUserId(userId)
}

func (gs *GoalServiceImpl) UpdateGoalCategory(categoryId, goalId, userId uuid.UUID) (*entities.GoalCategory, error) {
	if isInvalidUUID(categoryId) {
		return nil, errors.New("invalid uuid category id")
	}
	if isInvalidUUID(goalId) {
		return nil, errors.New("invalid uuid goal id")
	}
	gc, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if err != nil {
		return nil, err
	}

	if gc.UserId != userId {
		return nil, errors.New("user does not own this category")
	}

	return gs.goalCategoryStore.UpdateGoalCategory(categoryId, goalId)
}

func (gs *GoalServiceImpl) GetGoalCategoryById(categoryId uuid.UUID) (*entities.GoalCategory, error) {
	if isInvalidUUID(categoryId) {
		return nil, errors.New("invalid uuid category id")
	}
	return gs.goalCategoryStore.GetGoalCategoryById(categoryId)
}
