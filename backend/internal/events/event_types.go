package events

import "goalify/internal/entities"

type GoalUpdatedData struct {
	OldGoal *entities.Goal
	NewGoal *entities.Goal
}

type XpUpdatedData struct {
	LevelID int `json:"level_id"`
	Xp      int `json:"xp"`
}
