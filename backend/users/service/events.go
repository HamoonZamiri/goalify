package service

import (
	"goalify/entities"
	"goalify/utils/events"
	"log/slog"
)

func (us *userService) HandleEvent(event events.Event) {
	switch event.EventType {
	case events.GOAL_UPDATED:
		us.handleGoalUpdatedEvent(event)
	default:
		slog.Error("service.HandleEvent: unknown event type", "eventType", event.EventType)
	}
}

func (us *userService) handleGoalUpdatedEvent(event events.Event) {
	eventData, err := events.ParseEventData[map[string]any](event)
	if err != nil {
		slog.Error("service.handleGoalUpdatedEvent: events.ParseEventData:", "err", err)
		return
	}

	oldGoal := eventData["oldGoal"]
	newGoal := eventData["newGoal"]
	goalXp := eventData["xp"]

	assertedOldGoal, ok1 := oldGoal.(*entities.Goal)
	assertedNewGoal, ok2 := newGoal.(*entities.Goal)
	assertedXp, ok3 := goalXp.(int)
	if !ok1 || !ok2 || !ok3 {
		slog.Error("service.handleGoalUpdatedEvent: type assertion failed")
		return
	}

	if assertedOldGoal.Status != assertedNewGoal.Status && assertedNewGoal.Status == "complete" {
		// we need to update the xp of the user
		user, err := us.userStore.GetUserById(assertedOldGoal.UserId.String())
		if err != nil {
			slog.Error("service.handleGoalUpdatedEvent: store.GetUserById:", "err", err)
			return
		}
		level, err := us.userStore.GetLevelById(user.LevelId)
		if err != nil {
			slog.Error("service.handleGoalUpdatedEvent: store.GetLevelById:", "err", err)
			return
		}
		newXp := user.Xp + assertedXp
		newLevel := user.LevelId
		if newXp >= level.LevelUpXp {
			newXp %= level.LevelUpXp
			newLevel += 1
		}

		_, err = us.userStore.UpdateUserById(user.Id, map[string]any{
			"xp":       newXp,
			"level_id": newLevel,
		})
		if err != nil {
			slog.Error("service.handleGoalUpdatedEvent: store.UpdateUserById:", "err", err)
			return
		}
	}
}
