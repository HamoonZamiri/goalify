package service

import (
	"goalify/entities"

	"github.com/google/uuid"
)

type GoalService interface {
	// goals
	CreateGoal(title, description string, userId, categoryId uuid.UUID) (entities.Goal, error)
	UpdateGoalStatus(goalId uuid.UUID, status string) (entities.Goal, error)
	GetGoalsByUserId(userId uuid.UUID) ([]entities.Goal, error)
	UpdateGoalTitle(title string, goalId uuid.UUID) (entities.Goal, error)
	UpdateGoalDescription(description string, goalId uuid.UUID) (entities.Goal, error)
	GetGoalById(goalId uuid.UUID) (entities.Goal, error)

	// categories
	CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (entities.GoalCategory, error)
	GetGoalCategoriesByUserId(userId uuid.UUID) ([]entities.GoalCategory, error)
	UpdateGoalCategory(categoryId uuid.UUID, goalId uuid.UUID) (entities.GoalCategory, error)
	GetGoalCategoryById(categoryId uuid.UUID) (entities.GoalCategory, error)
}
