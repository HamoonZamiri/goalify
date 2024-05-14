package stores

import (
	"errors"
	"goalify/entities"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GoalCategoryStoreImpl struct {
	db *sqlx.DB
}

func NewGoalCategoryStore(db *sqlx.DB) *GoalCategoryStoreImpl {
	return &GoalCategoryStoreImpl{db: db}
}

func (s *GoalCategoryStoreImpl) UpdateGoalCategory(categoryId uuid.UUID, goalId uuid.UUID) (*entities.GoalCategory, error) {
	var goalCategory entities.GoalCategory
	err := s.db.QueryRowx("UPDATE goals SET category_id = $1 WHERE id = $2 RETURNING *", categoryId, goalId).StructScan(&goalCategory)
	if err != nil {
		return nil, err
	}
	return &goalCategory, err
}

func (s *GoalCategoryStoreImpl) CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (*entities.GoalCategory, error) {
	query := `INSERT INTO goal_categories (title, xp_per_goal, user_id) 
  VALUES ($1, $2, $3) RETURNING *`
	var goalCategory entities.GoalCategory
	err := s.db.QueryRowx(query, title, xpPerGoal, userId).StructScan(&goalCategory)
	if err != nil {
		return nil, err
	}
	return &goalCategory, err
}

func (s *GoalCategoryStoreImpl) GetGoalCategoriesByUserId(userId uuid.UUID) ([]*entities.GoalCategory, error) {
	query := `SELECT gc.id, gc.title, gc.xp_per_goal, gc.user_id, g.id as goal_id, 
  g.title as goal_title, g.description, g.status
  FROM goal_categories gc 
  JOIN goals g 
  ON gc.id = g.category_id 
  WHERE gc.user_id = $1`

	var categories []*entities.GoalCategory
	categoryMap := make(map[uuid.UUID]*entities.GoalCategory)
	rows, err := s.db.Queryx(query, userId)
	if err != nil {
		return nil, err
	}

	err = mapGoalCategoryRows(rows, categoryMap)
	if err != nil {
		slog.Error("error mapping joint goals and categories", "error", err)
		return nil, err
	}

	for _, category := range categoryMap {
		categories = append(categories, category)
	}

	return categories, nil
}

func (s *GoalCategoryStoreImpl) GetGoalCategoryById(categoryId uuid.UUID) (*entities.GoalCategory, error) {
	query := `SELECT gc.id, gc.title, gc.xp_per_goal, gc.user_id, g.id as goal_id, 
  g.title as goal_title, g.description, g.status
  FROM goal_categories gc JOIN goals g 
  ON gc.id = g.category_id 
  WHERE gc.id = $1`

	categoryMap := make(map[uuid.UUID]*entities.GoalCategory)

	rows, err := s.db.Queryx(query, categoryId)
	if err != nil {
		return nil, err
	}

	err = mapGoalCategoryRows(rows, categoryMap)
	if err != nil {
		return nil, err
	}

	if len(categoryMap) == 0 {
		return nil, errors.New("category not found")
	}

	return categoryMap[categoryId], nil
}

func mapGoalCategoryRows(rows *sqlx.Rows, categoryMap map[uuid.UUID]*entities.GoalCategory) error {
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

		goal := entities.Goal{
			Id:          goal_id,
			Title:       goal_title,
			Description: description,
			Status:      status,
			UserId:      user_id,
			CategoryId:  id,
		}

		if category, ok := categoryMap[id]; !ok {
			newCategory := entities.GoalCategory{
				Id:          id,
				Title:       title,
				Xp_per_goal: xp_per_goal,
				UserId:      user_id,
				Goals:       []*entities.Goal{&goal},
			}
			categoryMap[id] = &newCategory
		} else {
			category.Goals = append(category.Goals, &goal)
		}
	}
	return nil
}
