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

func (gs *GoalServiceImpl) CreateGoal(title, description string, userId, categoryId uuid.UUID) (*entities.Goal, error) {
	createdGoal, err := gs.goalStore.CreateGoal(title, description, userId, categoryId)
	if err != nil {
		slog.Error("service.CreateGoal: store.CreateGoal:", "err", err)
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
		slog.Error("service.GetGoalsByUserId: store.GetGoalsByUserId:", "err", err)
		return nil, fmt.Errorf("%w: error fetching goals", svcerror.ErrInternalServer)
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
	cat, err := gs.goalCategoryStore.CreateGoalCategory(title, xpPerGoal, userId)
	if err != nil {
		slog.Error("error creating goal category", "err", err)
		return nil, fmt.Errorf("%w: error creating goal category", svcerror.ErrInternalServer)
	}
	return cat, nil
}

func (gs *GoalServiceImpl) GetGoalCategoriesByUserId(userId uuid.UUID) ([]*entities.GoalCategory, error) {
	categories, err := gs.goalCategoryStore.GetGoalCategoriesByUserId(userId)
	if err == sql.ErrNoRows {
		return []*entities.GoalCategory{}, nil
	}
	if err != nil {
		slog.Error("service.GetGoalCategoriesByUserId: store.GetGoalCategoriesByUserId:", "err", err)
		return nil, fmt.Errorf("%w: error fetching goal categories", svcerror.ErrInternalServer)
	}
	return categories, nil
}

func (gs *GoalServiceImpl) UpdateGoalCategory(categoryId, goalId, userId uuid.UUID) (*entities.GoalCategory, error) {
	gc, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if err != nil {
		slog.Error("service.UpdateGoalCategory: store.GetGoalCategoryById:", "err", err)
		return nil, fmt.Errorf("%w: error fetching goal category", svcerror.ErrInternalServer)
	}

	if gc.UserId != userId {
		slog.Error("service.UpdateGoalCategory: user does not own this category", "userId", userId, "categoryId", categoryId, "ownerId", gc.UserId)
		return nil, errors.New("user does not own this category")
	}

	return gs.goalCategoryStore.UpdateGoalCategory(categoryId, goalId)
}

func (gs *GoalServiceImpl) GetGoalCategoryById(categoryId, userId uuid.UUID) (*entities.GoalCategory, error) {
	gc, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if err != nil {
		slog.Error("service.GetGoalCategoryById: store.GetGoalCategoryById:", "err", err)
		return nil, fmt.Errorf("%w: error fetching goal category", svcerror.ErrInternalServer)
	}

	if gc.UserId != userId {
		slog.Error("service.GetGoalCategoryById:", "err", "user does not own this category", "userId", userId, "categoryId", categoryId, "ownerId", gc.UserId)
		return nil, fmt.Errorf("%w: user does not own this category", svcerror.ErrBadRequest)
	}

	return gc, nil
}

func (gs *GoalServiceImpl) UpdateGoalCategoryById(categoryId uuid.UUID, updates map[string]interface{}, userId uuid.UUID) (*entities.GoalCategory, error) {
	cat, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if err != nil {
		slog.Error("service.UpdateGoalCategoryById: store.GetGoalCategoryById:", "err", err)
		return nil, fmt.Errorf("%w: error fetching goal category", svcerror.ErrInternalServer)
	}

	if cat.UserId != userId {
		slog.Error("service.UpdateGoalCategoryById:", "err", "user does not own this category")
		return nil, fmt.Errorf("%w: user does not own this category", svcerror.ErrBadRequest)
	}

	updatedCat, err := gs.goalCategoryStore.UpdateGoalCategoryById(categoryId, updates)
	if err != nil {
		slog.Error("service.UpdateGoalCategoryById: store.UpdateGoalCategoryById:", "err", err)
		return nil, fmt.Errorf("%w: error updating goal category", svcerror.ErrInternalServer)
	}
	return updatedCat, nil
}

func (gs *GoalServiceImpl) DeleteGoalCategoryById(categoryId, userId uuid.UUID) error {
	cat, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if err != nil {
		slog.Error("service.DeleteGoalCategoryById: store.GetGoalCategoryById:", "err", err)
		return fmt.Errorf("%w: error fetching goal category", svcerror.ErrInternalServer)
	}
	if cat.UserId != userId {
		slog.Error("service.DeleteGoalCategoryById:", "err", "user does not own this category", "userId", userId, "categoryId", categoryId, "ownerId", cat.UserId)
		return fmt.Errorf("%w: user does not own this category", svcerror.ErrUnauthorized)
	}

	err = gs.goalCategoryStore.DeleteGoalCategoryById(categoryId)
	if err == sql.ErrNoRows {
		slog.Error("service.DeleteGoalCategoryById: store.DeleteGoalCategoryById:", "err", "category not found")
		return fmt.Errorf("%w: category not found", svcerror.ErrNotFound)
	}

	if err != nil {
		slog.Error("service.DeleteGoalCategoryById: store.DeleteGoalCategoryById:", "err", err)
		return fmt.Errorf("%w: error deleting goal category", svcerror.ErrInternalServer)
	}

	return nil
}

func (gs *GoalServiceImpl) UpdateGoalById(goalId uuid.UUID, updates map[string]interface{}, userId uuid.UUID) (*entities.Goal, error) {
	goal, err := gs.goalStore.GetGoalById(goalId)
	if err != nil {
		return nil, err
	}
	if goal.UserId != userId {
		slog.Error("service.UpdateGoalById: user does not own this goal", "userId", userId, "goalId", goalId, "ownerId", goal.UserId)
		return nil, fmt.Errorf("%w: user does not own this goal", svcerror.ErrUnauthorized)
	}

	updatedGoal, err := gs.goalStore.UpdateGoalById(goalId, updates)
	if err != nil {
		slog.Error("service.UpdateGoalById: store.UpdateGoalById: ", "err", err)
		return nil, fmt.Errorf("%w: error updating goal", svcerror.ErrInternalServer)
	}
	return updatedGoal, nil
}
