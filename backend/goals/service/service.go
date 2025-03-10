package service

import (
	"database/sql"
	"errors"
	"fmt"
	"goalify/entities"
	"goalify/goals/stores"
	"goalify/utils/events"
	"goalify/utils/responses"
	"goalify/utils/stacktrace"
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
	_, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: invalid category id", responses.ErrBadRequest)
	}
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return nil, responses.ErrInternalServer
	}

	createdGoal, err := gs.goalStore.CreateGoal(title, description, userId, categoryId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.CreateGoal:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error creating goal", responses.ErrInternalServer)
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
	if errors.Is(err, sql.ErrNoRows) {
		return goals, nil
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalsByUserId:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error fetching goals", responses.ErrInternalServer)
	}
	return goals, nil
}

func (gs *goalService) GetGoalById(goalId uuid.UUID) (*entities.Goal, error) {
	funcStr := gs.traceLogger.GetTrace("service.GetGoalById")
	goal, err := gs.goalStore.GetGoalById(goalId)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: invalid goal id", responses.ErrNotFound)
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalById:", funcStr), "err", err)
		return nil, responses.ErrInternalServer
	}

	return goal, nil
}

func (gs *goalService) CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (*entities.GoalCategory, error) {
	funcStr := gs.traceLogger.GetTrace("service.CreateGoalCategory")

	cat, err := gs.goalCategoryStore.CreateGoalCategory(title, xpPerGoal, userId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.CreateGoalCategory:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error creating goal category", responses.ErrInternalServer)
	}

	e := events.NewEventWithUserId(events.GOAL_CATEGORY_CREATED, cat, cat.UserId.String())
	gs.eventPublisher.Publish(e)

	return cat, nil
}

func (gs *goalService) GetGoalCategoriesByUserId(userId uuid.UUID) ([]*entities.GoalCategory, error) {
	categories, err := gs.goalCategoryStore.GetGoalCategoriesByUserId(userId)
	if err != nil {
		return nil, responses.ErrInternalServer
	}
	return categories, nil
}

func (gs *goalService) GetGoalCategoryById(categoryId, userId uuid.UUID) (*entities.GoalCategory, error) {
	funcStr := gs.traceLogger.GetTrace("service.GetGoalCategoryById")

	invalidIdErr := fmt.Errorf("%w: invalid category id", responses.ErrBadRequest)

	gc, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, invalidIdErr
	}
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return nil, responses.ErrInternalServer
	}

	if gc.UserId != userId {
		slog.Error(fmt.Sprintf("%s: user does not own this category", funcStr), "userId", userId, "categoryId", categoryId, "ownerId", gc.UserId)
		return nil, invalidIdErr
	}

	return gc, nil
}

func (gs *goalService) UpdateGoalCategoryById(categoryId uuid.UUID, updates map[string]interface{}, userId uuid.UUID) (*entities.GoalCategory, error) {
	funcStr := gs.traceLogger.GetTrace("service.UpdateGoalCategoryById")
	badReqErr := fmt.Errorf("%w: invalid request", responses.ErrBadRequest)

	cat, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, badReqErr
	}
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return nil, responses.ErrInternalServer
	}

	if cat.UserId != userId {
		slog.Error("user does not own this category",
			"userId", userId,
			"categoryId", categoryId,
			"ownerId", cat.UserId,
			"trace", funcStr)
		return nil, badReqErr
	}

	updatedCat, err := gs.goalCategoryStore.UpdateGoalCategoryById(categoryId, updates)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.UpdateGoalCategoryById:", funcStr), "err", err)
		return nil, responses.ErrInternalServer
	}
	return updatedCat, nil
}

func (gs *goalService) DeleteGoalCategoryById(categoryId, userId uuid.UUID) error {
	funcStr := gs.traceLogger.GetTrace("service.DeleteGoalCategoryById")

	cat, err := gs.goalCategoryStore.GetGoalCategoryById(categoryId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return fmt.Errorf("%w: error fetching goal category", responses.ErrInternalServer)
	}

	if cat.UserId != userId {
		slog.Error(fmt.Sprintf("%s: user does not own this category", funcStr), "userId", userId, "categoryId", categoryId, "ownerId", cat.UserId)
		return fmt.Errorf("%w: user does not own this category", responses.ErrUnauthorized)
	}

	err = gs.goalCategoryStore.DeleteGoalCategoryById(categoryId)
	if err == sql.ErrNoRows {
		return fmt.Errorf("%w: category not found", responses.ErrNotFound)
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.DeleteGoalCategoryById:", funcStr), "err", err)
		return fmt.Errorf("%w: error deleting goal category", responses.ErrInternalServer)
	}

	return nil
}

func (gs *goalService) UpdateGoalById(goalId uuid.UUID, updates map[string]interface{}, userId uuid.UUID) (*entities.Goal, error) {
	funcStr := gs.traceLogger.GetTrace("service.UpdateGoalById")

	goal, err := gs.goalStore.GetGoalById(goalId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: invalid goal id", responses.ErrBadRequest)
	}
	if err != nil {
		slog.Error("service.handleUpdateGoalById: store.GetGoalById:", "err", err)
		return nil, responses.ErrInternalServer
	}
	if goal.UserId != userId {
		slog.Error(fmt.Sprintf("%s: user does not own this goal", funcStr), "userId", userId, "goalId", goalId, "ownerId", goal.UserId)
		return nil, fmt.Errorf("%w: invalid goal id", responses.ErrBadRequest)
	}

	categoryId := updates["category_id"]
	if categoryId != nil {
		categoryId, ok := categoryId.(string)
		if !ok {
			return nil, fmt.Errorf("%w: category_id must be a string", responses.ErrBadRequest)
		}
		parsedCategoryId, err := uuid.Parse(categoryId)
		if err != nil {
			slog.Error(fmt.Sprintf("%s: uuid.Parse(categoryId):", funcStr), "err", err)
			return nil, fmt.Errorf("%w: invalid category id", responses.ErrBadRequest)
		}

		cat, err := gs.goalCategoryStore.GetGoalCategoryById(parsedCategoryId)
		if err != nil {
			slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
			return nil, responses.ErrInternalServer
		}

		if cat.UserId != userId {
			slog.Error(fmt.Sprintf("%s: user does not own this category", funcStr), "userId", userId, "categoryId", parsedCategoryId, "ownerId", cat.UserId)
			return nil, fmt.Errorf("%w: invalid category id", responses.ErrBadRequest)
		}
	}

	updatedGoal, err := gs.goalStore.UpdateGoalById(goalId, updates)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: invalid goal id", responses.ErrBadRequest)
	}
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.UpdateGoalById:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error updating goal", responses.ErrInternalServer)
	}

	cat, err := gs.goalCategoryStore.GetGoalCategoryById(updatedGoal.CategoryId)
	skipEvent := false

	// if we can't fetch the goal category we don't know the correct xp per goal
	// skip publishing the event
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		skipEvent = true
	}
	if !skipEvent {
		eventData := &events.GoalUpdatedData{
			OldGoal: goal,
			NewGoal: updatedGoal,
			Xp:      cat.Xp_per_goal,
		}
		event := events.NewEventWithUserId(events.GOAL_UPDATED, eventData, userId.String())
		gs.eventPublisher.Publish(event)
	}
	return updatedGoal, nil
}

func (gs *goalService) DeleteGoalById(goalId, userId uuid.UUID) error {
	goal, err := gs.goalStore.GetGoalById(goalId)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: invalid goal id", responses.ErrBadRequest)
	}
	if err != nil {
		slog.Error("service.handleDeleteGoalById: store.GetGoalById:", "err", err)
		return responses.ErrInternalServer
	}
	if goal.UserId != userId {
		return fmt.Errorf("%w: invalid goal id", responses.ErrBadRequest)
	}

	err = gs.goalStore.DeleteGoalById(goalId)
	if err != nil {
		slog.Error("service.handleDeleteGoalById: store.DeleteGoalById:", "err", err)
		return responses.ErrInternalServer
	}
	return nil
}
