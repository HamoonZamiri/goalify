package stores

import (
	"goalify/entities"

	"github.com/google/uuid"
)

type GoalCategoryStore interface {
	CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (*entities.GoalCategory, error)
	GetGoalCategoriesByUserId(userId uuid.UUID) ([]*entities.GoalCategory, error)
	UpdateGoalCategory(categoryId uuid.UUID, goalId uuid.UUID) (*entities.GoalCategory, error)
	GetGoalCategoryById(categoryId uuid.UUID) (*entities.GoalCategory, error)
}
