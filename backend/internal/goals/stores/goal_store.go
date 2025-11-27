package stores

import (
	"context"
	"database/sql"
	"goalify/internal/entities"

	db "goalify/internal/db"
	sqlcdb "goalify/internal/db/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type GoalStore interface {
	CreateGoal(title, description string, userID, categoryID uuid.UUID) (*entities.Goal, error)
	UpdateGoalStatus(goalID, userID uuid.UUID, status string) (*entities.Goal, error)
	GetGoalsByUserID(userID uuid.UUID) ([]*entities.Goal, error)
	GetGoalByID(goalID, userID uuid.UUID) (*entities.Goal, error)
	UpdateGoalByID(goalID, userID uuid.UUID, updates map[string]any) (*entities.Goal, error)
	DeleteGoalByID(goalID, userID uuid.UUID) error
}

type goalStore struct {
	queries *sqlcdb.Queries
}

// Helper function to convert sqlc Goal to entity Goal
func pgxGoalToEntity(g sqlcdb.Goal) *entities.Goal {
	return &entities.Goal{
		ID:          uuid.UUID(g.ID.Bytes),
		Title:       g.Title,
		Description: g.Description.String,
		UserID:      uuid.UUID(g.UserID.Bytes),
		CategoryID:  uuid.UUID(g.CategoryID.Bytes),
		Status:      string(g.Status.GoalStatus),
		CreatedAt:   g.CreatedAt.Time,
		UpdatedAt:   g.UpdatedAt.Time,
	}
}

func NewGoalStore(queries *sqlcdb.Queries) GoalStore {
	return &goalStore{
		queries: queries,
	}
}

func (s *goalStore) CreateGoal(
	title, description string,
	userID, categoryID uuid.UUID,
) (*entities.Goal, error) {
	params := sqlcdb.CreateGoalParams{
		Title:       title,
		Description: pgtype.Text{String: description, Valid: true},
		UserID:      pgtype.UUID{Bytes: userID, Valid: true},
		CategoryID:  pgtype.UUID{Bytes: categoryID, Valid: true},
	}

	goal, err := s.queries.CreateGoal(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxGoalToEntity(goal), nil
}

func (s *goalStore) UpdateGoalStatus(
	goalID, userID uuid.UUID,
	status string,
) (*entities.Goal, error) {
	params := sqlcdb.UpdateGoalStatusParams{
		Status: sqlcdb.NullGoalStatus{GoalStatus: sqlcdb.GoalStatus(status), Valid: true},
		ID:     db.ToPgxUUID(goalID),
		UserID: db.ToPgxUUID(userID),
	}

	goal, err := s.queries.UpdateGoalStatus(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxGoalToEntity(goal), nil
}

func (s *goalStore) GetGoalsByUserID(userID uuid.UUID) ([]*entities.Goal, error) {
	goals, err := s.queries.GetGoalsByUserId(
		context.Background(),
		pgtype.UUID{Bytes: userID, Valid: true},
	)
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Goal, len(goals))
	for i, g := range goals {
		result[i] = pgxGoalToEntity(g)
	}

	return result, nil
}

func (s *goalStore) GetGoalByID(goalID, userID uuid.UUID) (*entities.Goal, error) {
	goal, err := s.queries.GetGoalById(
		context.Background(),
		sqlcdb.GetGoalByIdParams{
			ID:     db.ToPgxUUID(goalID),
			UserID: db.ToPgxUUID(userID),
		},
	)
	if err != nil {
		return nil, err
	}

	return pgxGoalToEntity(goal), nil
}

func anyToString(v any) (str string, ok bool) {
	if v == nil {
		return "", false
	}
	str, ok = v.(string)
	return str, ok
}

func (s *goalStore) UpdateGoalByID(
	goalID uuid.UUID,
	userID uuid.UUID,
	updates map[string]any,
) (*entities.Goal, error) {
	params := sqlcdb.UpdateGoalByIdParams{
		ID:     db.ToPgxUUID(goalID),
		UserID: db.ToPgxUUID(userID),
	}

	// Convert map updates to typed parameters
	if title, ok := anyToString(updates["title"]); ok {
		params.Title = pgtype.Text{String: title, Valid: true}
	}

	if description, ok := updates["description"]; ok {
		if descStr, ok := description.(string); ok {
			params.Description = pgtype.Text{String: descStr, Valid: true}
		}
	}
	if status, ok := updates["status"]; ok {
		if statusStr, ok := status.(string); ok {
			params.Status = sqlcdb.NullGoalStatus{
				GoalStatus: sqlcdb.GoalStatus(statusStr),
				Valid:      true,
			}
		}
	}
	if categoryID, ok := updates["category_id"]; ok {
		if catUUID, ok := categoryID.(uuid.UUID); ok {
			params.CategoryID = db.ToPgxUUID(catUUID)
		}
	}

	goal, err := s.queries.UpdateGoalById(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxGoalToEntity(goal), nil
}

func (s *goalStore) DeleteGoalByID(goalID, userID uuid.UUID) error {
	rows, err := s.queries.DeleteGoalById(context.Background(), sqlcdb.DeleteGoalByIdParams{
		ID:     db.ToPgxUUID(goalID),
		UserID: db.ToPgxUUID(userID),
	})
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
