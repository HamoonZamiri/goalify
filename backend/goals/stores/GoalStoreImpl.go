package stores

import (
	"errors"
	"goalify/models"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// type GoalStore interface {
//   CreateGoal(title, description string, userId, categoryId uuid.UUID) error
//   CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) error
//   UpdateGoalStatus(goalId uuid.UUID, status string) error
//   GetGoalsByUserId(userId uuid.UUID) ([]models.Goal, error)
//   GetGoalCategoriesByUserId(userId uuid.UUID) ([]models.GoalCategory, error)
//   UpdateGoalCategory(categoryId uuid.UUID, goalId uuid.UUID) error
//   UpdateGoalTitle(title string, goalId uuid.UUID) error
//   UpdateGoalDescription(description string, goalId uuid.UUID) error
//   GetGoalById(goalId uuid.UUID) (models.Goal, error)
//   GetGoalCategoryById(categoryId uuid.UUID) (models.GoalCategory, error)
// }

type GoalStoreImpl struct {
	db *sqlx.DB
}

func NewGoalStore(db *sqlx.DB) *GoalStoreImpl {
	return &GoalStoreImpl{db: db}
}

func (s *GoalStoreImpl) CreateGoal(title, description string, userId, categoryId uuid.UUID) error {
	_, err := s.db.Exec("INSERT INTO goals (title, description, user_id, category_id) VALUES ($1, $2, $3, $4)", title, description, userId, categoryId)
	if err != nil {
		return err
	}
	return nil
}

func (s *GoalStoreImpl) CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) error {
	_, err := s.db.Exec("INSERT INTO goal_categories (title, xp_per_goal, user_id) VALUES ($1, $2, $3)", title, xpPerGoal, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s *GoalStoreImpl) UpdateGoalStatus(goalId uuid.UUID, status string) error {
	_, err := s.db.Exec("UPDATE goals SET status = $1 WHERE id = $2", status, goalId)
	if err != nil {
		return err
	}
	return nil
}

func (s *GoalStoreImpl) GetGoalsByUserId(userId uuid.UUID) ([]models.Goal, error) {
	var goals []models.Goal

	err := s.db.Select(&goals, "SELECT * FROM goals WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (s *GoalStoreImpl) GetCategoriesByUserId(userId uuid.UUID) ([]models.GoalCategory, error) {
	query := `SELECT gc.id, gc.title, gc.xp_per_goal, gc.user_id, g.id as goal_id, 
  g.title as goal_title, g.description, g.status
  FROM goal_categories gc JOIN goals g ON gc.id = g.category_id WHERE gc.user_id = $1`

	var categories []models.GoalCategory
	categoryMap := make(map[uuid.UUID]models.GoalCategory)
	rows, err := s.db.Queryx(query, userId)
	if err != nil {
		return nil, err
	}

	err = mapGoalCategoryRows(rows, categoryMap)
	if err != nil {
		slog.Error("Error mapping joint goals and categories", "error", err)
		return nil, err
	}

	for _, category := range categoryMap {
		categories = append(categories, category)
	}

	return categories, nil
}

func (s *GoalStoreImpl) GetGoalCategoryById(categoryId uuid.UUID) (models.GoalCategory, error) {
	query := `SELECT gc.id, gc.title, gc.xp_per_goal, gc.user_id, g.id as goal_id, 
  g.title as goal_title, g.description, g.status
  FROM goal_categories gc JOIN goals g ON gc.id = g.category_id WHERE gc.id = $1`

	categoryMap := make(map[uuid.UUID]models.GoalCategory)

	rows, err := s.db.Queryx(query, categoryId)
	if err != nil {
		return models.GoalCategory{}, err
	}

	err = mapGoalCategoryRows(rows, categoryMap)
	if err != nil {
		slog.Error("Error mapping joint goals and categories", "error", err)
		return models.GoalCategory{}, err
	}

	if len(categoryMap) == 0 {
		return models.GoalCategory{}, errors.New("category not found")
	}

	return categoryMap[categoryId], nil
}

func mapGoalCategoryRows(rows *sqlx.Rows, categoryMap map[uuid.UUID]models.GoalCategory) error {
	for rows.Next() {
		var (
			id          uuid.UUID
			title       string
			xp_per_goal int
			user_id     uuid.UUID
			goal_id     uuid.UUID
			goal_title  string
			description string
			status      string
		)
		err := rows.Scan(&id, &title, &xp_per_goal, &user_id, &goal_id, &goal_title, &description, &status)
		if err != nil {
			return err
		}

		goal := models.Goal{
			Id:          goal_id,
			Title:       goal_title,
			Description: description,
			Status:      status,
			UserId:      user_id,
			CategoryId:  id,
		}

		if category, ok := categoryMap[id]; !ok {
			newCategory := models.GoalCategory{
				Id:          id,
				Title:       title,
				Xp_per_goal: xp_per_goal,
				UserId:      user_id,
				Goals:       []models.Goal{goal},
			}
			categoryMap[id] = newCategory
		} else {
			category.Goals = append(category.Goals, goal)
		}
	}
	return nil
}

func (s *GoalStoreImpl) UpdateGoalCategory(categoryId uuid.UUID, goalId uuid.UUID) error {
	_, err := s.db.Exec("UPDATE goals SET category_id = $1 WHERE id = $2", categoryId, goalId)
	if err != nil {
		return err
	}
	return nil
}

func (s *GoalStoreImpl) UpdateGoalTitle(title string, goalId uuid.UUID) error {
	_, err := s.db.Exec("UPDATE goals SET title = $1 WHERE id = $2", title, goalId)
	if err != nil {
		return err
	}
	return nil
}

func (s *GoalStoreImpl) UpdateGoalDescription(description string, goalId uuid.UUID) error {
	_, err := s.db.Exec("UPDATE goals SET description = $1 WHERE id = $2", description, goalId)
	if err != nil {
		return err
	}
	return nil
}

func (s *GoalStoreImpl) GetGoalById(goalId uuid.UUID) (models.Goal, error) {
	var goal models.Goal

	err := s.db.Get(&goal, "SELECT * FROM goals WHERE id = $1", goalId)
	if err != nil {
		return models.Goal{}, err
	}
	return goal, nil
}
