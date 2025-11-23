// Package service is the business logic layer for goals and goal categories
package service

import (
	"database/sql"
	"errors"
	"fmt"
	"goalify/internal/entities"
	"goalify/internal/events"
	"goalify/internal/goals/stores"
	"goalify/internal/responses"
	"goalify/pkg/stacktrace"
	"log/slog"
	"strings"

	"github.com/google/uuid"
)

const (
	XPPerGoalMax int = 100
)

var subscribedEvents = []string{events.GoalCategoryCreated, events.UserCreated}

type GoalService interface {
	// goals
	CreateGoal(title, description string, userID, categoryID uuid.UUID) (*entities.Goal, error)
	UpdateGoalStatus(status string, goalID, userID uuid.UUID) (*entities.Goal, error)
	GetGoalsByUserID(userID uuid.UUID) ([]*entities.Goal, error)
	GetGoalByID(goalID uuid.UUID) (*entities.Goal, error)
	UpdateGoalByID(
		goalID uuid.UUID,
		updates map[string]any,
		userID uuid.UUID,
	) (*entities.Goal, error)
	DeleteGoalByID(goalID, userID uuid.UUID) error

	// categories
	CreateGoalCategory(
		title string,
		xpPerGoal int,
		userID uuid.UUID,
	) (*entities.GoalCategory, error)
	GetGoalCategoriesByUserID(userID uuid.UUID) ([]*entities.GoalCategory, error)
	GetGoalCategoryByID(categoryID, userID uuid.UUID) (*entities.GoalCategory, error)
	UpdateGoalCategoryByID(
		categoryID uuid.UUID,
		updates map[string]any,
		userID uuid.UUID,
	) (*entities.GoalCategory, error)
	DeleteGoalCategoryByID(categoryID, userID uuid.UUID) error
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

func (gs *goalService) CreateGoal(
	title, description string,
	userID, categoryID uuid.UUID,
) (*entities.Goal, error) {
	funcStr := gs.traceLogger.GetTrace("service.CreateGoal")
	_, err := gs.goalCategoryStore.GetGoalCategoryByID(categoryID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: invalid category id", responses.ErrBadRequest)
	}
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return nil, responses.ErrInternalServer
	}

	createdGoal, err := gs.goalStore.CreateGoal(title, description, userID, categoryID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.CreateGoal:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error creating goal", responses.ErrInternalServer)
	}
	return createdGoal, nil
}

func (gs *goalService) UpdateGoalStatus(
	status string,
	goalID, userID uuid.UUID,
) (*entities.Goal, error) {
	cleanStatus := strings.ToLower(status)

	if cleanStatus != "completed" && cleanStatus != "not_completed" {
		return nil, errors.New("status must be either 'completed' or 'not_completed'")
	}

	goal, err := gs.goalStore.GetGoalByID(goalID)
	if err != nil {
		return nil, err
	}
	if goal.UserID != userID {
		return nil, errors.New("user does not own this goal")
	}

	return gs.goalStore.UpdateGoalStatus(goalID, status)
}

func (gs *goalService) GetGoalsByUserID(userID uuid.UUID) ([]*entities.Goal, error) {
	funcStr := gs.traceLogger.GetTrace("service.GetGoalsByUserId")
	goals, err := gs.goalStore.GetGoalsByUserID(userID)
	if errors.Is(err, sql.ErrNoRows) {
		return goals, nil
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalsByUserId:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error fetching goals", responses.ErrInternalServer)
	}
	return goals, nil
}

func (gs *goalService) GetGoalByID(goalID uuid.UUID) (*entities.Goal, error) {
	funcStr := gs.traceLogger.GetTrace("service.GetGoalById")
	goal, err := gs.goalStore.GetGoalByID(goalID)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: invalid goal id", responses.ErrNotFound)
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalById:", funcStr), "err", err)
		return nil, responses.ErrInternalServer
	}

	return goal, nil
}

func (gs *goalService) CreateGoalCategory(
	title string,
	xpPerGoal int,
	userID uuid.UUID,
) (*entities.GoalCategory, error) {
	funcStr := gs.traceLogger.GetTrace("service.CreateGoalCategory")

	cat, err := gs.goalCategoryStore.CreateGoalCategory(title, xpPerGoal, userID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.CreateGoalCategory:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error creating goal category", responses.ErrInternalServer)
	}

	e := events.NewEventWithUserID(events.GoalCategoryCreated, cat, cat.UserID.String())
	slog.Info("Publishing GOAL_CATEGORY_CREATED event",
		slog.String("categoryId", cat.ID.String()),
		slog.String("userId", cat.UserID.String()))
	gs.eventPublisher.Publish(e)

	return cat, nil
}

func (gs *goalService) GetGoalCategoriesByUserID(
	userID uuid.UUID,
) ([]*entities.GoalCategory, error) {
	categories, err := gs.goalCategoryStore.GetGoalCategoriesByUserID(userID)
	if err != nil {
		return nil, responses.ErrInternalServer
	}
	return categories, nil
}

func (gs *goalService) GetGoalCategoryByID(
	categoryID, userID uuid.UUID,
) (*entities.GoalCategory, error) {
	funcStr := gs.traceLogger.GetTrace("service.GetGoalCategoryById")

	invalidIDErr := fmt.Errorf("%w: invalid category id", responses.ErrBadRequest)

	gc, err := gs.goalCategoryStore.GetGoalCategoryByID(categoryID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, invalidIDErr
	}
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return nil, responses.ErrInternalServer
	}

	if gc.UserID != userID {
		slog.Error(
			fmt.Sprintf("%s: user does not own this category", funcStr),
			"userId",
			userID,
			"categoryId",
			categoryID,
			"ownerId",
			gc.UserID,
		)
		return nil, invalidIDErr
	}

	return gc, nil
}

func (gs *goalService) UpdateGoalCategoryByID(
	categoryID uuid.UUID,
	updates map[string]any,
	userID uuid.UUID,
) (*entities.GoalCategory, error) {
	funcStr := gs.traceLogger.GetTrace("service.UpdateGoalCategoryById")
	badReqErr := fmt.Errorf("%w: invalid request", responses.ErrBadRequest)

	cat, err := gs.goalCategoryStore.GetGoalCategoryByID(categoryID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, badReqErr
	}
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return nil, responses.ErrInternalServer
	}

	if cat.UserID != userID {
		slog.Error("user does not own this category",
			"userId", userID,
			"categoryId", categoryID,
			"ownerId", cat.UserID,
			"trace", funcStr)
		return nil, badReqErr
	}

	updatedCat, err := gs.goalCategoryStore.UpdateGoalCategoryByID(categoryID, updates)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.UpdateGoalCategoryById:", funcStr), "err", err)
		return nil, responses.ErrInternalServer
	}
	return updatedCat, nil
}

func (gs *goalService) DeleteGoalCategoryByID(categoryID, userID uuid.UUID) error {
	funcStr := gs.traceLogger.GetTrace("service.DeleteGoalCategoryById")

	cat, err := gs.goalCategoryStore.GetGoalCategoryByID(categoryID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.GetGoalCategoryById:", funcStr), "err", err)
		return fmt.Errorf("%w: error fetching goal category", responses.ErrInternalServer)
	}

	if cat.UserID != userID {
		slog.Error(
			fmt.Sprintf("%s: user does not own this category", funcStr),
			"userId",
			userID,
			"categoryId",
			categoryID,
			"ownerId",
			cat.UserID,
		)
		return fmt.Errorf("%w: user does not own this category", responses.ErrUnauthorized)
	}

	err = gs.goalCategoryStore.DeleteGoalCategoryByID(categoryID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: category not found", responses.ErrNotFound)
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.DeleteGoalCategoryById:", funcStr), "err", err)
		return fmt.Errorf("%w: error deleting goal category", responses.ErrInternalServer)
	}

	return nil
}

func (gs *goalService) checkCategoryForUpdate(catID any, userID uuid.UUID) bool {
	if catID == nil {
		return true
	}

	catIDString, ok := catID.(string)
	if !ok {
		slog.Error("service.checkCategoryForUpdate: invalid category id")
		return false
	}

	catIDUUID, err := uuid.Parse(catIDString)
	if err != nil {
		slog.Error("service.checkCategoryForUpdate: uuid.Parse:", "err", err)
		return false
	}

	cat, err := gs.goalCategoryStore.GetGoalCategoryByID(catIDUUID)
	if err != nil {
		return false
	}
	return cat.UserID.String() == userID.String()
}

func (gs *goalService) UpdateGoalByID(
	goalID uuid.UUID,
	updates map[string]any,
	userID uuid.UUID,
) (*entities.Goal, error) {
	funcStr := gs.traceLogger.GetTrace("service.UpdateGoalById")

	goal, err := gs.goalStore.GetGoalByID(goalID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: invalid goal id", responses.ErrBadRequest)
	}

	if err != nil {
		slog.Error("service.handleUpdateGoalById: store.GetGoalById:", "err", err)
		return nil, responses.ErrInternalServer
	}

	if goal.UserID != userID {
		slog.Error(
			fmt.Sprintf("%s: user does not own this goal", funcStr),
			"userID",
			userID,
			"goalID",
			goalID,
			"ownerID",
			goal.UserID,
		)
		return nil, fmt.Errorf("%w: invalid goal id", responses.ErrBadRequest)
	}

	categoryOk := gs.checkCategoryForUpdate(updates["category_id"], userID)
	if !categoryOk {
		return nil, fmt.Errorf("%w: invalid category id", responses.ErrBadRequest)
	}

	updatedGoal, err := gs.goalStore.UpdateGoalByID(goalID, updates)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: invalid goal id", responses.ErrBadRequest)
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: store.UpdateGoalById:", funcStr), "err", err)
		return nil, fmt.Errorf("%w: error updating goal", responses.ErrInternalServer)
	}

	cat, err := gs.goalCategoryStore.GetGoalCategoryByID(updatedGoal.CategoryID)
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
			Xp:      cat.XPPerGoal,
		}
		event := events.NewEventWithUserID(events.GoalUpdated, eventData, userID.String())
		gs.eventPublisher.Publish(event)
	}
	return updatedGoal, nil
}

func (gs *goalService) DeleteGoalByID(goalID, userID uuid.UUID) error {
	goal, err := gs.goalStore.GetGoalByID(goalID)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: invalid goal id", responses.ErrBadRequest)
	}
	if err != nil {
		slog.Error("service.handleDeleteGoalById: store.GetGoalById:", "err", err)
		return responses.ErrInternalServer
	}
	if goal.UserID != userID {
		return fmt.Errorf("%w: invalid goal id", responses.ErrBadRequest)
	}

	err = gs.goalStore.DeleteGoalByID(goalID)
	if err != nil {
		slog.Error("service.handleDeleteGoalById: store.DeleteGoalById:", "err", err)
		return responses.ErrInternalServer
	}
	return nil
}
