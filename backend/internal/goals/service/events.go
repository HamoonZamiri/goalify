package service

import (
	"goalify/internal/entities"
	"goalify/internal/events"
	"log/slog"
)

func (gs *goalService) HandleEvent(event events.Event) {
	switch event.EventType {
	case events.UserCreated:
		gs.handleUserCreatedEvent(event)
	case events.GoalCategoryCreated:
		gs.handleGoalCategoryCreatedEvent(event)
	default:
		slog.Error("service.HandleEvent: unknown event type", "eventType", event.EventType)
	}
}

func (gs *goalService) handleUserCreatedEvent(event events.Event) {
	user, err := events.ParseEventData[*entities.User](event)
	if err != nil {
		slog.Error("service.handleUserCreatedEvent: events.ParseEventData:", "err", err)
		return
	}

	_, err = gs.CreateGoalCategory("Daily", XPPerGoalMax, user.ID)
	if err != nil {
		slog.Error("service.handleUserCreatedEvent: CreateGoalCategory:", "err", err)
	}
}

func (gs *goalService) handleGoalCategoryCreatedEvent(event events.Event) {
	slog.Info("Handling GOAL_CATEGORY_CREATED event",
		slog.String("eventType", event.EventType),
		slog.String("userId", event.UserID.ValueOrZero()))
	category, err := events.ParseEventData[*entities.GoalCategory](event)
	if err != nil {
		slog.Error("service.handleGoalCategoryCreatedEvent: events.ParseEventData:", "err", err)
		return
	}

	defaultGoal, err := gs.CreateGoal(
		"example",
		"This is an example goal/task!",
		category.UserID,
		category.ID,
	)
	if err != nil {
		slog.Error("service.handleGoalCategoryCreatedEvent: CreateGoal:", "err", err)
		return
	}
	slog.Info("Publishing DEFAULT_GOAL_CREATED event",
		slog.String("goalId", defaultGoal.ID.String()),
		slog.String("userId", defaultGoal.UserID.String()))
	gs.eventPublisher.Publish(
		events.NewEventWithUserID(
			events.DefaultGoalCreated,
			defaultGoal,
			defaultGoal.UserID.String(),
		),
	)
}
