// Package stores is the repository layer for the loot domain
package stores

import (
	"context"
	"goalify/internal/entities"

	sqlcdb "goalify/internal/db/generated"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type LootStore interface {
	CreateChest(chestType, description string, price int) (*entities.Chest, error)
	GetChestByID(chestID uuid.UUID) (*entities.Chest, error)
	GetAllChests() ([]*entities.Chest, error)
	UpdateChestByID(chestID uuid.UUID, updates map[string]any) (*entities.Chest, error)
	DeleteChestByID(chestID uuid.UUID) error
}

type lootStore struct {
	queries *sqlcdb.Queries
}

// Helper function to convert sqlc Chest to entity Chest
func pgxChestToEntity(c sqlcdb.Chest) *entities.Chest {
	return &entities.Chest{
		ID:          uuid.UUID(c.ID.Bytes),
		Type:        string(c.Type),
		Description: c.Description,
		Price:       int(c.Price),
		CreatedAt:   c.CreatedAt.Time,
		UpdatedAt:   c.UpdatedAt.Time,
	}
}

func NewChestStore(queries *sqlcdb.Queries) LootStore {
	return &lootStore{queries: queries}
}

func (s *lootStore) CreateChest(chestType, description string, price int) (*entities.Chest, error) {
	params := sqlcdb.CreateChestParams{
		Type:        sqlcdb.ChestType(chestType),
		Description: description,
		Price:       int32(price),
	}

	chest, err := s.queries.CreateChest(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxChestToEntity(chest), nil
}

func (s *lootStore) GetChestByID(chestID uuid.UUID) (*entities.Chest, error) {
	chest, err := s.queries.GetChestById(
		context.Background(),
		pgtype.UUID{Bytes: chestID, Valid: true},
	)
	if err != nil {
		return nil, err
	}

	return pgxChestToEntity(chest), nil
}

func (s *lootStore) GetAllChests() ([]*entities.Chest, error) {
	chests, err := s.queries.GetAllChests(context.Background())
	if err != nil {
		return nil, err
	}

	result := make([]*entities.Chest, len(chests))
	for i, c := range chests {
		result[i] = pgxChestToEntity(c)
	}

	return result, nil
}

func (s *lootStore) UpdateChestByID(
	chestID uuid.UUID,
	updates map[string]any,
) (*entities.Chest, error) {
	params := sqlcdb.UpdateChestByIdParams{
		ID: pgtype.UUID{Bytes: chestID, Valid: true},
	}

	// Convert map updates to typed parameters
	if chestType, ok := updates["type"]; ok {
		if typeStr, ok := chestType.(string); ok {
			params.Type = sqlcdb.NullChestType{ChestType: sqlcdb.ChestType(typeStr), Valid: true}
		}
	}
	if description, ok := updates["description"]; ok {
		if descStr, ok := description.(string); ok {
			params.Description = pgtype.Text{String: descStr, Valid: true}
		}
	}
	if price, ok := updates["price"]; ok {
		if priceInt, ok := price.(int); ok {
			params.Price = pgtype.Int4{Int32: int32(priceInt), Valid: true}
		}
	}

	chest, err := s.queries.UpdateChestById(context.Background(), params)
	if err != nil {
		return nil, err
	}

	return pgxChestToEntity(chest), nil
}

func (s *lootStore) DeleteChestByID(chestID uuid.UUID) error {
	return s.queries.DeleteChestById(context.Background(), pgtype.UUID{Bytes: chestID, Valid: true})
}
