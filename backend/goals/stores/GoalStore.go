package stores

import (
	"goalify/models"

	"github.com/google/uuid"
)

type GoalStore interface {
	CreateGoal(title, description string, userId, categoryId uuid.UUID) error
	CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) error
	UpdateGoalStatus(goalId uuid.UUID, status string) error
	GetGoalsByUserId(userId uuid.UUID) ([]models.Goal, error)
	GetGoalCategoriesByUserId(userId uuid.UUID) ([]models.GoalCategory, error)
	UpdateGoalCategory(categoryId uuid.UUID, goalId uuid.UUID) error
	UpdateGoalTitle(title string, goalId uuid.UUID) error
	UpdateGoalDescription(description string, goalId uuid.UUID) error
	GetGoalById(goalId uuid.UUID) (models.Goal, error)
	GetGoalCategoryById(categoryId uuid.UUID) (models.GoalCategory, error)
}
