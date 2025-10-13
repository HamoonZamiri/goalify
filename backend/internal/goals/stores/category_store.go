package stores

import (
	"context"
	"database/sql"
	sqlcdb "goalify/internal/db/generated"
	"goalify/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
		queries *sqlcdb.Queries
	}
)

// Helper function to convert sqlc GoalCategory to entity GoalCategory
func pgxGoalCategoryToEntity(gc sqlcdb.GoalCategory) *entities.GoalCategory {
	return &entities.GoalCategory{
		Id:          uuid.UUID(gc.ID.Bytes),
		Title:       gc.Title,
		Xp_per_goal: int(gc.XpPerGoal),
		UserId:      uuid.UUID(gc.UserID.Bytes),
		CreatedAt:   gc.CreatedAt.Time,
		UpdatedAt:   gc.UpdatedAt.Time,
		Goals:       []*entities.Goal{}, // Initialize empty slice
	}
}

// Helper to map JOIN query rows to GoalCategories with nested Goals
func mapGoalCategoryWithGoalsRows(rows []sqlcdb.GetGoalCategoriesWithGoalsByUserIdRow) []*entities.GoalCategory {
	categoryMap := make(map[uuid.UUID]*entities.GoalCategory)
	categorySlice := make([]*entities.GoalCategory, 0)

	for _, row := range rows {
		categoryId := uuid.UUID(row.ID.Bytes)

		// Create category if it doesn't exist in map
		if _, ok := categoryMap[categoryId]; !ok {
			gc := &entities.GoalCategory{
				Id:          categoryId,
				Title:       row.Title,
				Xp_per_goal: int(row.XpPerGoal),
				UserId:      uuid.UUID(row.UserID.Bytes),
				CreatedAt:   row.CreatedAt.Time,
				UpdatedAt:   row.UpdatedAt.Time,
				Goals:       []*entities.Goal{},
			}
			categoryMap[categoryId] = gc
			categorySlice = append(categorySlice, gc)
		}

		// Add goal if it exists (LEFT JOIN may have null goals)
		if row.GoalID.Valid {
			goal := &entities.Goal{
				Id:          uuid.UUID(row.GoalID.Bytes),
				Title:       row.GoalTitle.String,
				Description: row.Description.String,
				Status:      string(row.Status.GoalStatus),
				CategoryId:  categoryId,
				UserId:      uuid.UUID(row.UserID.Bytes),
				CreatedAt:   row.GoalCreatedAt.Time,
				UpdatedAt:   row.GoalUpdatedAt.Time,
			}
			categoryMap[categoryId].Goals = append(categoryMap[categoryId].Goals, goal)
		}
	}

	return categorySlice
}

// Helper for single category with goals
func mapGoalCategoryWithGoalsSingleRow(rows []sqlcdb.GetGoalCategoryWithGoalsByIdRow) (*entities.GoalCategory, error) {
	if len(rows) == 0 {
		return nil, sql.ErrNoRows
	}

	firstRow := rows[0]
	gc := &entities.GoalCategory{
		Id:          uuid.UUID(firstRow.ID.Bytes),
		Title:       firstRow.Title,
		Xp_per_goal: int(firstRow.XpPerGoal),
		UserId:      uuid.UUID(firstRow.UserID.Bytes),
		CreatedAt:   firstRow.CreatedAt.Time,
		UpdatedAt:   firstRow.UpdatedAt.Time,
		Goals:       []*entities.Goal{},
	}

	for _, row := range rows {
		if row.GoalID.Valid {
			goal := &entities.Goal{
				Id:          uuid.UUID(row.GoalID.Bytes),
				Title:       row.GoalTitle.String,
				Description: row.Description.String,
				Status:      string(row.Status.GoalStatus),
				CategoryId:  uuid.UUID(row.ID.Bytes),
				UserId:      uuid.UUID(row.UserID.Bytes),
				CreatedAt:   row.GoalCreatedAt.Time,
				UpdatedAt:   row.GoalUpdatedAt.Time,
			}
			gc.Goals = append(gc.Goals, goal)
		}
	}

	return gc, nil
}

func NewGoalCategoryStore(queries *sqlcdb.Queries) GoalCategoryStore {
	return &goalCategoryStore{
		queries: queries,
	}
}


func (s *goalCategoryStore) CreateGoalCategory(title string, xpPerGoal int, userId uuid.UUID) (*entities.GoalCategory, error) {
	params := sqlcdb.CreateGoalCategoryParams{
		Title:     title,
		XpPerGoal: int32(xpPerGoal),
		UserID:    pgtype.UUID{Bytes: userId, Valid: true},
	}

	gc, err := s.queries.CreateGoalCategory(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxGoalCategoryToEntity(gc), nil
}

func (s *goalCategoryStore) GetGoalCategoriesByUserId(userId uuid.UUID) ([]*entities.GoalCategory, error) {
	rows, err := s.queries.GetGoalCategoriesWithGoalsByUserId(context.Background(), pgtype.UUID{Bytes: userId, Valid: true})
	if err != nil {
		return nil, err
	}

	return mapGoalCategoryWithGoalsRows(rows), nil
}

func (s *goalCategoryStore) GetGoalCategoryById(categoryId uuid.UUID) (*entities.GoalCategory, error) {
	rows, err := s.queries.GetGoalCategoryWithGoalsById(context.Background(), pgtype.UUID{Bytes: categoryId, Valid: true})
	if err != nil {
		return nil, err
	}

	return mapGoalCategoryWithGoalsSingleRow(rows)
}

func (s *goalCategoryStore) UpdateGoalCategoryById(categoryId uuid.UUID, updates map[string]any) (*entities.GoalCategory, error) {
	params := sqlcdb.UpdateGoalCategoryByIdParams{
		ID: pgtype.UUID{Bytes: categoryId, Valid: true},
	}

	// Convert map updates to typed parameters
	if title, ok := updates["title"]; ok {
		if titleStr, ok := title.(string); ok {
			params.Title = pgtype.Text{String: titleStr, Valid: true}
		}
	}
	if xpPerGoal, ok := updates["xp_per_goal"]; ok {
		if xpInt, ok := xpPerGoal.(int); ok {
			params.XpPerGoal = pgtype.Int4{Int32: int32(xpInt), Valid: true}
		}
	}

	gc, err := s.queries.UpdateGoalCategoryById(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxGoalCategoryToEntity(gc), nil
}

func (s *goalCategoryStore) DeleteGoalCategoryById(categoryId uuid.UUID) error {
	return s.queries.DeleteGoalCategoryById(context.Background(), pgtype.UUID{Bytes: categoryId, Valid: true})
}

