package events

import "goalify/internal/entities"

type GoalUpdatedData struct {
	OldGoal *entities.Goal
	NewGoal *entities.Goal
	Xp      int
}

type XpUpdatedData struct {
	LevelId int `json:"level_id"`
	Xp      int `json:"xp"`
}
