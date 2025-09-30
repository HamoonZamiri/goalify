package stores

import (
	"context"
	sqlcdb "goalify/db/generated"
	"goalify/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	db      *sqlx.DB
	queries *sqlcdb.Queries
}

// Helper function to convert sqlc Goal to entity Goal
func pgxGoalToEntity(g sqlcdb.Goal) *entities.Goal {
	return &entities.Goal{
		Id:          uuid.UUID(g.ID.Bytes),
		Title:       g.Title,
		Description: g.Description.String,
		UserId:      uuid.UUID(g.UserID.Bytes),
		CategoryId:  uuid.UUID(g.CategoryID.Bytes),
		Status:      string(g.Status.GoalStatus),
		CreatedAt:   g.CreatedAt.Time,
		UpdatedAt:   g.UpdatedAt.Time,
	}
}

func NewGoalStore(db *sqlx.DB, queries *sqlcdb.Queries) GoalStore {
	return &goalStore{
		db:      db,
		queries: queries,
	}
}

func (s *goalStore) CreateGoal(title, description string, userId, categoryId uuid.UUID) (*entities.Goal, error) {
	params := sqlcdb.CreateGoalParams{
		Title:       title,
		Description: pgtype.Text{String: description, Valid: true},
		UserID:      pgtype.UUID{Bytes: userId, Valid: true},
		CategoryID:  pgtype.UUID{Bytes: categoryId, Valid: true},
	}

	goal, err := s.queries.CreateGoal(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxGoalToEntity(goal), nil
}

func (s *goalStore) UpdateGoalStatus(goalId uuid.UUID, status string) (*entities.Goal, error) {
	params := sqlcdb.UpdateGoalStatusParams{
		Status: sqlcdb.NullGoalStatus{GoalStatus: sqlcdb.GoalStatus(status), Valid: true},
		ID:     pgtype.UUID{Bytes: goalId, Valid: true},
	}

	goal, err := s.queries.UpdateGoalStatus(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxGoalToEntity(goal), nil
}

func (s *goalStore) GetGoalsByUserId(userId uuid.UUID) ([]*entities.Goal, error) {
	goals, err := s.queries.GetGoalsByUserId(context.Background(), pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Goal, len(goals))
	for i, g := range goals {
		result[i] = pgxGoalToEntity(g)
	}

	return result, nil
}

func (s *goalStore) GetGoalById(goalId uuid.UUID) (*entities.Goal, error) {
	goal, err := s.queries.GetGoalById(context.Background(), pgtype.UUID{Bytes: goalId, Valid: true})
	if err != nil {
		return nil, err
	}

	return pgxGoalToEntity(goal), nil
}

func (s *goalStore) UpdateGoalById(goalId uuid.UUID, updates map[string]interface{}) (*entities.Goal, error) {
	params := sqlcdb.UpdateGoalByIdParams{
		ID: pgtype.UUID{Bytes: goalId, Valid: true},
	}

	// Convert map updates to typed parameters
	if title, ok := updates["title"]; ok {
		if titleStr, ok := title.(string); ok {
			params.Title = pgtype.Text{String: titleStr, Valid: true}
		}
	}
	if description, ok := updates["description"]; ok {
		if descStr, ok := description.(string); ok {
			params.Description = pgtype.Text{String: descStr, Valid: true}
		}
	}
	if status, ok := updates["status"]; ok {
		if statusStr, ok := status.(string); ok {
			params.Status = sqlcdb.NullGoalStatus{GoalStatus: sqlcdb.GoalStatus(statusStr), Valid: true}
		}
	}
	if categoryId, ok := updates["category_id"]; ok {
		if catUUID, ok := categoryId.(uuid.UUID); ok {
			params.CategoryID = pgtype.UUID{Bytes: catUUID, Valid: true}
		}
	}

	goal, err := s.queries.UpdateGoalById(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxGoalToEntity(goal), nil
}

func (s *goalStore) DeleteGoalById(goalId uuid.UUID) error {
	return s.queries.DeleteGoalById(context.Background(), pgtype.UUID{Bytes: goalId, Valid: true})
}
