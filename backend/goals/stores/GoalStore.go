package stores

import (
	"goalify/entities"

	"github.com/google/uuid"
)

type GoalStore interface {
	CreateGoal(title, description string, userId, categoryId uuid.UUID) (*entities.Goal, error)
	UpdateGoalStatus(goalId uuid.UUID, status string) (*entities.Goal, error)
	GetGoalsByUserId(userId uuid.UUID) ([]*entities.Goal, error)
	UpdateGoalTitle(title string, goalId uuid.UUID) (*entities.Goal, error)
	UpdateGoalDescription(description string, goalId uuid.UUID) (*entities.Goal, error)
	GetGoalById(goalId uuid.UUID) (*entities.Goal, error)
	UpdateGoalById(goalId uuid.UUID, updates map[string]interface{}) (*entities.Goal, error)
}
