package service

import (
	"goalify/entities"
	"goalify/utils/events"
	"log/slog"
)

func (gs *GoalServiceImpl) HandleEvent(event events.Event) {
	switch event.EventType {
	case events.USER_CREATED:
		gs.handleUserCreatedEvent(event)
	case events.GOAL_CATEGORY_CREATED:
		gs.handleGoalCategoryCreatedEvent(event)
	default:
		slog.Error("service.HandleEvent: unknown event type", "eventType", event.EventType)
	}
}

func (gs *GoalServiceImpl) handleUserCreatedEvent(event events.Event) {
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

func (gs *GoalServiceImpl) handleGoalCategoryCreatedEvent(event events.Event) {
	category, err := events.ParseEventData[*entities.GoalCategory](event)
	if err != nil {
		slog.Error("service.handleGoalCategoryCreatedEvent: events.ParseEventData:", "err", err)
		return
	}

	_, err = gs.CreateGoal("example", "This is an example goal/task!", category.UserId, category.Id)
	if err != nil {
		slog.Error("service.handleGoalCategoryCreatedEvent: CreateGoal:", "err", err)
	}
}
