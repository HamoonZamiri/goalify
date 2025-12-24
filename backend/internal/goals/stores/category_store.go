// Package stores is the repository layer package for goal categories
package stores

import (
	"context"
	"database/sql"
	"goalify/internal/entities"
	"goalify/pkg/options"

	db "goalify/internal/db"
	sqlcdb "goalify/internal/db/generated"

	"github.com/google/uuid"
)

type UpdateGoalCategoryParams struct {
	Title     options.Option[string]
	XpPerGoal options.Option[int]
}

type (
	GoalCategoryStore interface {
		CreateGoalCategory(
			title string,
			xpPerGoal int,
			userID uuid.UUID,
		) (*entities.GoalCategory, error)
		GetGoalCategoriesByUserID(userID uuid.UUID) ([]*entities.GoalCategory, error)
		GetGoalCategoryByID(categoryID, userID uuid.UUID) (*entities.GoalCategory, error)
		UpdateGoalCategoryByID(
			categoryID, userID uuid.UUID,
			params UpdateGoalCategoryParams,
		) (*entities.GoalCategory, error)
		DeleteGoalCategoryByID(categoryID, userID uuid.UUID) error
	}
	goalCategoryStore struct {
		queries *sqlcdb.Queries
	}
)

// Helper function to convert sqlc GoalCategory to entity GoalCategory
func pgxGoalCategoryToEntity(gc sqlcdb.GoalCategory) *entities.GoalCategory {
	return &entities.GoalCategory{
		ID:        uuid.UUID(gc.ID.Bytes),
		Title:     gc.Title,
		XPPerGoal: int(gc.XpPerGoal),
		UserID:    uuid.UUID(gc.UserID.Bytes),
		CreatedAt: gc.CreatedAt.Time,
		UpdatedAt: gc.UpdatedAt.Time,
		Goals:     []*entities.Goal{}, // Initialize empty slice
	}
}

// Helper to map JOIN query rows to GoalCategories with nested Goals
func mapGoalCategoryWithGoalsRows(
	rows []sqlcdb.GetGoalCategoriesWithGoalsByUserIdRow,
) []*entities.GoalCategory {
	categoryMap := make(map[uuid.UUID]*entities.GoalCategory)
	categorySlice := make([]*entities.GoalCategory, 0)

	for _, row := range rows {
		categoryID := uuid.UUID(row.ID.Bytes)

		// Create category if it doesn't exist in map
		if _, ok := categoryMap[categoryID]; !ok {
			gc := &entities.GoalCategory{
				ID:        categoryID,
				Title:     row.Title,
				XPPerGoal: int(row.XpPerGoal),
				UserID:    uuid.UUID(row.UserID.Bytes),
				CreatedAt: row.CreatedAt.Time,
				UpdatedAt: row.UpdatedAt.Time,
				Goals:     []*entities.Goal{},
			}
			categoryMap[categoryID] = gc
			categorySlice = append(categorySlice, gc)
		}

		// Add goal if it exists (LEFT JOIN may have null goals)
		if row.GoalID.Valid {
			goal := &entities.Goal{
				ID:          uuid.UUID(row.GoalID.Bytes),
				Title:       row.GoalTitle.String,
				Description: row.Description.String,
				Status:      string(row.Status.GoalStatus),
				CategoryID:  categoryID,
				UserID:      uuid.UUID(row.UserID.Bytes),
				CreatedAt:   row.GoalCreatedAt.Time,
				UpdatedAt:   row.GoalUpdatedAt.Time,
			}
			categoryMap[categoryID].Goals = append(categoryMap[categoryID].Goals, goal)
		}
	}

	return categorySlice
}

// Helper for single category with goals
func mapGoalCategoryWithGoalsSingleRow(
	rows []sqlcdb.GetGoalCategoryWithGoalsByIdRow,
) (*entities.GoalCategory, error) {
	if len(rows) == 0 {
		return nil, sql.ErrNoRows
	}

	firstRow := rows[0]
	gc := &entities.GoalCategory{
		ID:        uuid.UUID(firstRow.ID.Bytes),
		Title:     firstRow.Title,
		XPPerGoal: int(firstRow.XpPerGoal),
		UserID:    uuid.UUID(firstRow.UserID.Bytes),
		CreatedAt: firstRow.CreatedAt.Time,
		UpdatedAt: firstRow.UpdatedAt.Time,
		Goals:     []*entities.Goal{},
	}

	for _, row := range rows {
		if row.GoalID.Valid {
			goal := &entities.Goal{
				ID:          uuid.UUID(row.GoalID.Bytes),
				Title:       row.GoalTitle.String,
				Description: row.Description.String,
				Status:      string(row.Status.GoalStatus),
				CategoryID:  uuid.UUID(row.ID.Bytes),
				UserID:      uuid.UUID(row.UserID.Bytes),
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

func (s *goalCategoryStore) CreateGoalCategory(
	title string,
	xpPerGoal int,
	userID uuid.UUID,
) (*entities.GoalCategory, error) {
	params := sqlcdb.CreateGoalCategoryParams{
		Title:     title,
		XpPerGoal: int32(xpPerGoal),
		UserID:    db.UUIDToPgxUUID(userID),
	}

	gc, err := s.queries.CreateGoalCategory(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxGoalCategoryToEntity(gc), nil
}

func (s *goalCategoryStore) GetGoalCategoriesByUserID(
	userID uuid.UUID,
) ([]*entities.GoalCategory, error) {
	rows, err := s.queries.GetGoalCategoriesWithGoalsByUserId(
		context.Background(),
		db.UUIDToPgxUUID(userID),
	)
	if err != nil {
		return nil, err
	}

	return mapGoalCategoryWithGoalsRows(rows), nil
}

func (s *goalCategoryStore) GetGoalCategoryByID(
	categoryID, userID uuid.UUID,
) (*entities.GoalCategory, error) {
	rows, err := s.queries.GetGoalCategoryWithGoalsById(
		context.Background(),
		sqlcdb.GetGoalCategoryWithGoalsByIdParams{
			ID:     db.UUIDToPgxUUID(categoryID),
			UserID: db.UUIDToPgxUUID(userID),
		})
	if err != nil {
		return nil, err
	}

	return mapGoalCategoryWithGoalsSingleRow(rows)
}

func (s *goalCategoryStore) UpdateGoalCategoryByID(
	categoryID, userID uuid.UUID,
	params UpdateGoalCategoryParams,
) (*entities.GoalCategory, error) {
	sqlcParams := sqlcdb.UpdateGoalCategoryByIdParams{
		ID:     db.UUIDToPgxUUID(categoryID),
		UserID: db.UUIDToPgxUUID(userID),
	}

	sqlcParams.Title = db.OptionStringToPgxText(params.Title)

	xpInt4, err := db.OptionIntToPgxInt4(params.XpPerGoal)
	if err != nil {
		return nil, err
	}
	sqlcParams.XpPerGoal = xpInt4

	gc, err := s.queries.UpdateGoalCategoryById(context.Background(), sqlcParams)
	if err != nil {
		return nil, err
	}

	return pgxGoalCategoryToEntity(gc), nil
}

func (s *goalCategoryStore) DeleteGoalCategoryByID(categoryID, userID uuid.UUID) error {
	rows, err := s.queries.DeleteGoalCategoryById(
		context.Background(),
		sqlcdb.DeleteGoalCategoryByIdParams{
			ID:     db.UUIDToPgxUUID(categoryID),
			UserID: db.UUIDToPgxUUID(userID),
		})
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
