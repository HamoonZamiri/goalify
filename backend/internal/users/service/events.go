package service

import (
	"goalify/internal/events"
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
	eventData, err := events.ParseEventData[*events.GoalUpdatedData](event)
	if err != nil {
		slog.Error("service.handleGoalUpdatedEvent: events.ParseEventData:", "err", err)
		return
	}

	oldGoal := eventData.OldGoal
	newGoal := eventData.NewGoal
	xp := eventData.Xp

	if oldGoal.Status != newGoal.Status && newGoal.Status == "complete" {
		// we need to update the xp of the user
		user, err := us.userStore.GetUserById(oldGoal.UserId.String())
		if err != nil {
			slog.Error("service.handleGoalUpdatedEvent: store.GetUserById:", "err", err)
			return
		}
		level, err := us.userStore.GetLevelById(user.LevelId)
		if err != nil {
			slog.Error("service.handleGoalUpdatedEvent: store.GetLevelById:", "err", err)
			return
		}
		newXp := user.Xp + xp
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

		eventData := &events.XpUpdatedData{
			LevelId: newLevel,
			Xp:      newXp,
		}
		us.eventPublisher.Publish(events.NewEventWithUserId(events.XP_UPDATED, eventData, oldGoal.UserId.String()))

	}
}
