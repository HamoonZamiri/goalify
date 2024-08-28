package service

import (
	"database/sql"
	"errors"
	"fmt"
	"goalify/entities"
	"goalify/goals/stores"
	"goalify/utils/events"
	"goalify/utils/stacktrace"
	"goalify/utils/svcerror"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

const (
	XP_PER_GOAL_MAX int = 100
)

var subscribedEvents = []string{events.GOAL_CATEGORY_CREATED, events.USER_CREATED}

type GoalService interface {
	// goals
	CreateGoal(title, description string, userId, categoryId uuid.UUID) (*entities.Goal, error)
	UpdateGoalStatus(status string, goalId, userId uuid.UUID) (*entities.Goal, error)
	GetGoalsByUserId(userId uuid.UUID) ([]*entities.Goal, error)
	GetGoalById(goalId uuid.UUID) (*entities.Goal, error)
	UpdateGoalById(goalId uuid.UUID, updates map[string]any, userId uuid.UUID) (*entities.Goal, error)
	DeleteGoalById(goalId, userId uuid.UUID) error

	// categories
	CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (*entities.GoalCategory, error)
	GetGoalCategoriesByUserId(userId uuid.UUID) ([]*entities.GoalCategory, error)
	GetGoalCategoryById(categoryId, userId uuid.UUID) (*entities.GoalCategory, error)
	UpdateGoalCategoryById(categoryId uuid.UUID, updates map[string]any, userId uuid.UUID) (*entities.GoalCategory, error)
	DeleteGoalCategoryById(categoryId, userId uuid.UUID) error
}

type goalService struct {
	goalStore         stores.GoalStore
	goalCategoryStore stores.GoalCategoryStore
	traceLogger       stacktrace.TraceLogger
	eventPublisher    events.EventPublisher
}

func NewGoalService(goalStore stores.GoalStore,
	goalCategoryStore stores.GoalCategoryStore,
	traceLogger stacktrace.TraceLogger, ep events.EventPublisher,
) GoalService {
	gs := &goalService{
		goalStore:         goalStore,
		goalCategoryStore: goalCategoryStore,
		traceLogger:       traceLogger,
		eventPublisher:    ep,
	}

	for _, event := range subscribedEvents {
		ep.Subscribe(event, gs)
	}

	return gs
}

func (gs *goalService) CreateGoal(title, description string, userId, categoryId uuid.UUID) (*entities.Goal, error) {
	funcStr := gs.traceLogger.GetTrace("service.CreateGoal")

	createdGoal, err := gs.goalStore.CreateGoal(title, description, userId, categoryId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.CreateGoal:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error creating goal", svcerror.ErrInternalServer)
	}
	return createdGoal, nil
}

func (gs *goalService) UpdateGoalStatus(status string, goalId, userId uuid.UUID) (*entities.Goal, error) {
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

func (gs *goalService) GetGoalsByUserId(userId uuid.UUID) ([]*entities.Goal, error) {
	funcStr := gs.traceLogger.GetTrace("service.GetGoalsByUserId")
	goals, err := gs.goalStore.GetGoalsByUserId(userId)
	if err == sql.ErrNoRows {
		return goals, nil
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalsByUserId:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error fetching goals", svcerror.ErrInternalServer)
	}
	return goals, nil
}

func (gs *goalService) GetGoalById(goalId uuid.UUID) (*entities.Goal, error) {
	funcStr := gs.traceLogger.GetTrace("service.GetGoalById")
	goal, err := gs.goalStore.GetGoalById(goalId)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("%w: goal not found", svcerror.ErrNotFound)
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalById:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error fetching goal", svcerror.ErrInternalServer)
	}

	return goal, nil
}

func (gs *goalService) CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (*entities.GoalCategory, error) {
	funcStr := gs.traceLogger.GetTrace("service.CreateGoalCategory")

	cat, err := gs.goalCategoryStore.CreateGoalCategory(title, xpPerGoal, userId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.CreateGoalCategory:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error creating goal category", svcerror.ErrInternalServer)
	}

	e := events.NewEventWithUserId(events.GOAL_CATEGORY_CREATED, cat, cat.UserId.String())
	gs.eventPublisher.Publish(e)

	return cat, nil
}

func (gs *goalService) GetGoalCategoriesByUserId(userId uuid.UUID) ([]*entities.GoalCategory, error) {
	funcStr := gs.traceLogger.GetTrace("service.GetGoalCategoriesByUserId")

	categories, err := gs.goalCategoryStore.GetGoalCategoriesByUserId(userId)
	if err == sql.ErrNoRows {
		return []*entities.GoalCategory{}, nil
	}
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoriesByUserId:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error fetching goal categories", svcerror.ErrInternalServer)
	}
	return categories, nil
}

func (gs *goalService) GetGoalCategoryById(categoryId, userId uuid.UUID) (*entities.GoalCategory, error) {
	funcStr := gs.traceLogger.GetTrace("service.GetGoalCategoryById")

	gc, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error fetching goal category", svcerror.ErrInternalServer)
	}

	if gc.UserId != userId {
		slog.Error(fmt.Sprintf("%s: user does not own this category", funcStr), "userId", userId, "categoryId", categoryId, "ownerId", gc.UserId)
		return nil, fmt.Errorf("%w: user does not own this category", svcerror.ErrBadRequest)
	}

	return gc, nil
}

func (gs *goalService) UpdateGoalCategoryById(categoryId uuid.UUID, updates map[string]interface{}, userId uuid.UUID) (*entities.GoalCategory, error) {
	funcStr := gs.traceLogger.GetTrace("service.UpdateGoalCategoryById")

	cat, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error fetching goal category", svcerror.ErrInternalServer)
	}

	if cat.UserId != userId {
		slog.Error(fmt.Sprintf("%s: user does not own this category", funcStr), "userId", userId, "categoryId", categoryId, "ownerId", cat.UserId)
		return nil, fmt.Errorf("%w: user does not own this category", svcerror.ErrBadRequest)
	}

	updatedCat, err := gs.goalCategoryStore.UpdateGoalCategoryById(categoryId, updates)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.UpdateGoalCategoryById:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error updating goal category", svcerror.ErrInternalServer)
	}
	return updatedCat, nil
}

func (gs *goalService) DeleteGoalCategoryById(categoryId, userId uuid.UUID) error {
	funcStr := gs.traceLogger.GetTrace("service.DeleteGoalCategoryById")

	cat, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return fmt.Errorf("%w: error fetching goal category", svcerror.ErrInternalServer)
	}

	if cat.UserId != userId {
		slog.Error(fmt.Sprintf("%s: user does not own this category", funcStr), "userId", userId, "categoryId", categoryId, "ownerId", cat.UserId)
		return fmt.Errorf("%w: user does not own this category", svcerror.ErrUnauthorized)
	}

	err = gs.goalCategoryStore.DeleteGoalCategoryById(categoryId)
	if err == sql.ErrNoRows {
		return fmt.Errorf("%w: category not found", svcerror.ErrNotFound)
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.DeleteGoalCategoryById:", funcStr), "err", err)
		return fmt.Errorf("%w: error deleting goal category", svcerror.ErrInternalServer)
	}

	return nil
}

func (gs *goalService) UpdateGoalById(goalId uuid.UUID, updates map[string]interface{}, userId uuid.UUID) (*entities.Goal, error) {
	funcStr := gs.traceLogger.GetTrace("service.UpdateGoalById")

	goal, err := gs.goalStore.GetGoalById(goalId)
	if err != nil {
		return nil, err
	}
	if goal.UserId != userId {
		slog.Error(fmt.Sprintf("%s: user does not own this goal", funcStr), "userId", userId, "goalId", goalId, "ownerId", goal.UserId)
		return nil, fmt.Errorf("%w: user does not own this goal", svcerror.ErrUnauthorized)
	}

	categoryId := updates["category_id"]
	if categoryId != nil {
		categoryId, ok := categoryId.(string)
		if !ok {
			return nil, fmt.Errorf("%w: category_id must be a string", svcerror.ErrBadRequest)
		}
		parsedCategoryId, err := uuid.Parse(categoryId)
		if err != nil {
			slog.Error(fmt.Sprintf("%s: uuid.Parse(categoryId):", funcStr), "err", err)
			return nil, fmt.Errorf("%w: error parsing category_id", svcerror.ErrBadRequest)
		}

		cat, err := gs.goalCategoryStore.GetGoalCategoryById(parsedCategoryId)
		if err != nil {
			slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
			return nil, fmt.Errorf("%w: error fetching goal category", svcerror.ErrInternalServer)
		}

		if cat.UserId != userId {
			slog.Error(fmt.Sprintf("%s: user does not own this category", funcStr), "userId", userId, "categoryId", parsedCategoryId, "ownerId", cat.UserId)
			return nil, fmt.Errorf("%w: user does not own this category", svcerror.ErrUnauthorized)
		}
	}

	updatedGoal, err := gs.goalStore.UpdateGoalById(goalId, updates)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.UpdateGoalById:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error updating goal", svcerror.ErrInternalServer)
	}

	cat, err := gs.goalCategoryStore.GetGoalCategoryById(updatedGoal.CategoryId)
	event := events.NewEventWithUserId(events.GOAL_UPDATED, map[string]any{
		"oldGoal": goal,
		"newGoal": updatedGoal,
		"xp":      cat.Xp_per_goal,
	},
		userId.String())
	gs.eventPublisher.Publish(event)
	return updatedGoal, nil
}

func (gs *goalService) DeleteGoalById(goalId, userId uuid.UUID) error {
	goal, err := gs.goalStore.GetGoalById(goalId)
	if err == sql.ErrNoRows {
		slog.Error("service.handleDeleteGoalById: store.GetGoalById:", "err", err)
		return fmt.Errorf("%w: goal not found", svcerror.ErrNotFound)
	}
	if err != nil {
		slog.Error("service.handleDeleteGoalById: store.GetGoalById:", "err", err)
		return fmt.Errorf("%w: error fetching goal", svcerror.ErrBadRequest)
	}
	if goal.UserId != userId {
		return fmt.Errorf("%w: user does not own this goal", svcerror.ErrUnauthorized)
	}
	err = gs.goalStore.DeleteGoalById(goalId)
	if err != nil {
		slog.Error("service.handleDeleteGoalById: store.DeleteGoalById:", "err", err)
		return fmt.Errorf("%w: error deleting goal", svcerror.ErrInternalServer)
	}
	return nil
}
