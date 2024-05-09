package stores

import (
	"goalify/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GoalStoreImpl struct {
	db *sqlx.DB
}

func NewGoalStore(db *sqlx.DB) *GoalStoreImpl {
	return &GoalStoreImpl{db: db}
}

func (s *GoalStoreImpl) CreateGoal(title, description string, userId, categoryId uuid.UUID) (*entities.Goal, error) {
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

func (s *GoalStoreImpl) UpdateGoalStatus(goalId uuid.UUID, status string) (*entities.Goal, error) {
	query := `UPDATE goals SET status = $1 WHERE id = $2 RETURNING *`

	var goal entities.Goal
	err := s.db.QueryRowx(query, status, goalId).StructScan(&goal)
	if err != nil {
		return nil, err
	}

	return &goal, nil
}

func (s *GoalStoreImpl) GetGoalsByUserId(userId uuid.UUID) ([]*entities.Goal, error) {
	var goals []*entities.Goal

	rows, err := s.db.Queryx("SELECT * FROM goals WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var goal entities.Goal
		err = rows.StructScan(&goal)
		if err != nil {
			return nil, err
		}
		goals = append(goals, &goal)
	}
	return goals, nil
}

func (s *GoalStoreImpl) UpdateGoalTitle(title string, goalId uuid.UUID) (*entities.Goal, error) {
	query := `UPDATE goals SET title = $1 WHERE id = $2 RETURNING *`
	var goal entities.Goal
	err := s.db.QueryRowx(query, title, goalId).StructScan(&goal)
	if err != nil {
		return nil, err
	}
	return &goal, err
}

func (s *GoalStoreImpl) UpdateGoalDescription(description string, goalId uuid.UUID) (*entities.Goal, error) {
	query := `UPDATE goals SET description = $1 WHERE id = $2 RETURNING *`

	var goal entities.Goal
	err := s.db.QueryRowx(query, description, goalId).StructScan(&goal)
	if err != nil {
		return nil, err
	}

	return &goal, nil
}

func (s *GoalStoreImpl) GetGoalById(goalId uuid.UUID) (*entities.Goal, error) {
	var goal entities.Goal

	err := s.db.Get(&goal, "SELECT * FROM goals WHERE id = $1", goalId)
	if err != nil {
		return nil, err
	}
	return &goal, nil
}
