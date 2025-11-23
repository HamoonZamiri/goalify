package service

import (
	"goalify/internal/events"
	"log/slog"
)

func (s *userService) HandleEvent(event events.Event) {
	switch event.EventType {
	case events.GoalUpdated:
		s.handleGoalUpdatedEvent(event)
	default:
		slog.Error("service.HandleEvent: unknown event type", "eventType", event.EventType)
	}
}

func (s *userService) handleGoalUpdatedEvent(event events.Event) {
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
		user, err := s.userStore.GetUserByID(oldGoal.UserID.String())
		if err != nil {
			slog.Error("service.handleGoalUpdatedEvent: store.GetUserById:", "err", err)
			return
		}
		level, err := s.userStore.GetLevelByID(user.LevelID)
		if err != nil {
			slog.Error("service.handleGoalUpdatedEvent: store.GetLevelById:", "err", err)
			return
		}
		newXp := user.Xp + xp
		newLevel := user.LevelID
		if newXp >= level.LevelUpXp {
			newXp %= level.LevelUpXp
			newLevel += 1
		}

		_, err = s.userStore.UpdateUserByID(user.ID, map[string]any{
			"xp":       newXp,
			"level_id": newLevel,
		})
		if err != nil {
			slog.Error("service.handleGoalUpdatedEvent: store.UpdateUserById:", "err", err)
			return
		}

		eventData := &events.XpUpdatedData{
			LevelID: newLevel,
			Xp:      newXp,
		}
		s.eventPublisher.Publish(
			events.NewEventWithUserID(events.XPUpdated, eventData, oldGoal.UserID.String()),
		)

	}
}
