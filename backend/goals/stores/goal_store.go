package stores

import (
	"fmt"
	"goalify/db"
	"goalify/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GoalStore interface {
	CreateGoal(title, description string, userId, categoryId uuid.UUID) (*entities.Goal, error)
	UpdateGoalStatus(goalId uuid.UUID, status string) (*entities.Goal, error)
	GetGoalsByUserId(userId uuid.UUID) ([]*entities.Goal, error)
	GetGoalById(goalId uuid.UUID) (*entities.Goal, error)
	UpdateGoalById(goalId uuid.UUID, updates map[string]interface{}) (*entities.Goal, error)
	DeleteGoalById(goalId uuid.UUID) error
}

type goalStore struct {
	db *sqlx.DB
}

func NewGoalStore(db *sqlx.DB) GoalStore {
	return &goalStore{db: db}
}

func (s *goalStore) CreateGoal(title, description string, userId, categoryId uuid.UUID) (*entities.Goal, error) {
	query := `INSERT INTO goals (title, description, user_id, category_id)
  VALUES ($1, $2, $3, $4)
  RETURNING *`

	var goal entities.Goal
	err := s.db.QueryRowx(query, title, description, userId, categoryId).StructScan(&goal)
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (s *goalStore) UpdateGoalStatus(goalId uuid.UUID, status string) (*entities.Goal, error) {
	query := `UPDATE goals SET status = $1 WHERE id = $2 RETURNING *`

	var goal entities.Goal
	err := s.db.QueryRowx(query, status, goalId).StructScan(&goal)
	if err != nil {
		return nil, err
	}

	return &goal, nil
}

func (s *goalStore) GetGoalsByUserId(userId uuid.UUID) ([]*entities.Goal, error) {
	goals := make([]*entities.Goal, 0)

	err := s.db.Select(&goals, "SELECT * FROM goals WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	return goals, nil
}

func (s *goalStore) GetGoalById(goalId uuid.UUID) (*entities.Goal, error) {
	var goal entities.Goal

	err := s.db.Get(&goal, "SELECT * FROM goals WHERE id = $1", goalId)
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (s *goalStore) UpdateGoalById(goalId uuid.UUID, updates map[string]interface{}) (*entities.Goal, error) {
	query, args := db.BuildUpdateQuery("goals", updates, goalId)

	var goal entities.Goal
	err := s.db.QueryRowx(fmt.Sprintf("%s RETURNING *", query), args...).StructScan(&goal)
	if err != nil {
		return nil, err
	}
	return &goal, nil
}

func (s *goalStore) DeleteGoalById(goalId uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM goals WHERE id = $1", goalId)
	return err
}
