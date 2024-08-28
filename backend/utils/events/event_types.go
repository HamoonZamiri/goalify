package events

import "goalify/entities"

type GoalUpdatedData struct {
	OldGoal *entities.Goal
	NewGoal *entities.Goal
	Xp      int
}
