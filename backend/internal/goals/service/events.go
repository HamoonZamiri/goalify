package service

import (
	"goalify/internal/entities"
	"goalify/internal/events"
	"log/slog"
)

func (gs *goalService) HandleEvent(event events.Event) {
	switch event.EventType {
	case events.USER_CREATED:
		gs.handleUserCreatedEvent(event)
	case events.GOAL_CATEGORY_CREATED:
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

	_, err = gs.CreateGoalCategory("Daily", XP_PER_GOAL_MAX, user.Id)
	if err != nil {
		slog.Error("service.handleUserCreatedEvent: CreateGoalCategory:", "err", err)
	}
}

func (gs *goalService) handleGoalCategoryCreatedEvent(event events.Event) {
	slog.Info("Handling GOAL_CATEGORY_CREATED event",
		slog.String("eventType", event.EventType),
		slog.String("userId", event.UserId.ValueOrZero()))
	category, err := events.ParseEventData[*entities.GoalCategory](event)
	if err != nil {
		slog.Error("service.handleGoalCategoryCreatedEvent: events.ParseEventData:", "err", err)
		return
	}

	defaultGoal, err := gs.CreateGoal("example", "This is an example goal/task!", category.UserId, category.Id)
	if err != nil {
		slog.Error("service.handleGoalCategoryCreatedEvent: CreateGoal:", "err", err)
		return
	}
	slog.Info("Publishing DEFAULT_GOAL_CREATED event",
		slog.String("goalId", defaultGoal.Id.String()),
		slog.String("userId", defaultGoal.UserId.String()))
	gs.eventPublisher.Publish(events.NewEventWithUserId(events.DEFAULT_GOAL_CREATED, defaultGoal, defaultGoal.UserId.String()))
}
