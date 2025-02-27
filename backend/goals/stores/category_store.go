package stores

import (
	"database/sql"
	"fmt"
	"goalify/db"
	"goalify/entities"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	GoalCategoryStore interface {
		CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (*entities.GoalCategory, error)
		GetGoalCategoriesByUserId(userId uuid.UUID) ([]*entities.GoalCategory, error)
		GetGoalCategoryById(categoryId uuid.UUID) (*entities.GoalCategory, error)
		UpdateGoalCategoryById(categoryId uuid.UUID, updates map[string]any) (*entities.GoalCategory, error)
		DeleteGoalCategoryById(categoryId uuid.UUID) error
	}
	goalCategoryStore struct {
		db *sqlx.DB
	}
)

func NewGoalCategoryStore(db *sqlx.DB) GoalCategoryStore {
	return &goalCategoryStore{db: db}
}

const category_join_goals_query string = `SELECT gc.id, gc.title, gc.xp_per_goal, gc.user_id, gc.created_at, gc.updated_at, g.id as goal_id, 
  g.title as goal_title, g.description, g.status, g.created_at as goal_created_at, g.updated_at as goal_updated_at
  FROM goal_categories gc LEFT JOIN goals g 
  ON gc.id = g.category_id`

func (s *goalCategoryStore) CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (*entities.GoalCategory, error) {
	query := `INSERT INTO goal_categories (title, xp_per_goal, user_id) 
  VALUES ($1, $2, $3) RETURNING *`
	gc := entities.NewGoalCategory()
	err := s.db.QueryRowx(query, title, xpPerGoal, userId).StructScan(gc)
	if err != nil {
		return nil, err
	}
	return gc, err
}

func (s *goalCategoryStore) GetGoalCategoriesByUserId(userId uuid.UUID) ([]*entities.GoalCategory, error) {
	query := category_join_goals_query + ` WHERE gc.user_id = $1 ORDER BY gc.created_at`

	rows, err := s.db.Queryx(query, userId)
	if err != nil {
		return nil, err
	}

	categorySlice, err := mapGoalCategoryRows(rows)
	if err != nil {
		slog.Error("error mapping joint goals and categories", "error", err)
		return nil, err
	}

	return categorySlice, nil
}

func (s *goalCategoryStore) GetGoalCategoryById(categoryId uuid.UUID) (*entities.GoalCategory, error) {
	query := category_join_goals_query + ` WHERE gc.id = $1`

	rows, err := s.db.Queryx(query, categoryId)
	if err != nil {
		return nil, err
	}

	categorySlice, err := mapGoalCategoryRows(rows)
	if err != nil {
		return nil, err
	}

	if len(categorySlice) == 0 {
		return nil, sql.ErrNoRows
	}

	return categorySlice[0], nil
}

func (s *goalCategoryStore) UpdateGoalCategoryById(categoryId uuid.UUID, updates map[string]any) (*entities.GoalCategory, error) {
	query, args := db.BuildUpdateQuery("goal_categories", updates, categoryId)
	gc := entities.NewGoalCategory()

	err := s.db.QueryRowx(fmt.Sprintf("%s RETURNING *", query), args...).StructScan(gc)
	if err != nil {
		return nil, fmt.Errorf("queryrowx: %w", err)
	}
	return gc, nil
}

func (s *goalCategoryStore) DeleteGoalCategoryById(categoryId uuid.UUID) error {
	query := `DELETE FROM goal_categories WHERE id = $1`
	_, err := s.db.Exec(query, categoryId)
	if err != nil {
		return err
	}
	return nil
}

func mapGoalCategoryRows(rows *sqlx.Rows) ([]*entities.GoalCategory, error) {
	var results []entities.CategoryWithGoalRow
	err := sqlx.StructScan(rows, &results)
	categoryMap := make(map[uuid.UUID]*entities.GoalCategory)
	categorySlice := make([]*entities.GoalCategory, 0)

	if err != nil {
		return categorySlice, err
	}

	for _, result := range results {
		gc := result.ToGoalCategory()
		if _, ok := categoryMap[gc.Id]; !ok {
			categoryMap[gc.Id] = gc
			categorySlice = append(categorySlice, gc)
		}

		if result.GoalId.Valid {
			goal := result.ToGoal()
			categoryMap[gc.Id].Goals = append(categoryMap[gc.Id].Goals, goal)
		}
	}

	return categorySlice, nil
}
