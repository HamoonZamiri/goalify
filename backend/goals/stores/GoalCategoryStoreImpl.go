package stores

import (
	"database/sql"
	"errors"
	"fmt"
	"goalify/db"
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
  LEFT JOIN goals g 
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
  FROM goal_categories gc LEFT JOIN goals g 
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

func (s *GoalCategoryStoreImpl) UpdateGoalCategoryById(categoryId uuid.UUID, updates map[string]any) (*entities.GoalCategory, error) {
	query, args := db.BuildUpdateQuery("goal_categories", updates, categoryId)
	var gc entities.GoalCategory

	err := s.db.QueryRowx(fmt.Sprintf("%s RETURNING *", query), args...).StructScan(&gc)
	if err != nil {
		return nil, fmt.Errorf("queryrowx: %w", err)
	}
	return &gc, nil
}

func mapGoalCategoryRows(rows *sqlx.Rows, categoryMap map[uuid.UUID]*entities.GoalCategory) error {
	for rows.Next() {
		var gc entities.GoalCategory
		var goalId, goalTitle, goalDescription, goalStatus sql.NullString

		err := rows.Scan(&gc.Id, &gc.Title, &gc.Xp_per_goal, &gc.UserId, &goalId, &goalTitle, &goalDescription, &goalStatus)
		if err != nil {
			return err
		}

		if _, ok := categoryMap[gc.Id]; !ok {
			gc.Goals = []*entities.Goal{}
			categoryMap[gc.Id] = &gc
		}

		if goalId.Valid {
			goal := entities.Goal{
				Id:          uuid.MustParse(goalId.String),
				Title:       goalTitle.String,
				Description: goalDescription.String,
				Status:      goalStatus.String,
			}
			categoryMap[gc.Id].Goals = append(categoryMap[gc.Id].Goals, &goal)
		}
	}
	return nil
}
