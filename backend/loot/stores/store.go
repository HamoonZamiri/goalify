package stores

import (
	"fmt"
	"goalify/entities"

	"github.com/jmoiron/sqlx"
)

type LootStore interface {
	CreateChest(chestType, description string, price int) (*entities.Chest, error)
}

type lootStore struct {
	db *sqlx.DB
}

func NewChestStore(db *sqlx.DB) LootStore {
	return &lootStore{db: db}
}

func (s *lootStore) CreateChest(chestType, description string, price int) (*entities.Chest, error) {
	query := fmt.Sprintf(`INSERT INTO %s (type, description, price) VALUES ($1, $2, $3) RETURNING *`, CHEST_TABLE)

	var chest entities.Chest
	err := s.db.QueryRowx(query, chestType, description, price).StructScan(&chest)
	if err != nil {
		return nil, err
	}
	return &chest, nil
}
