package stores

import (
	"goalify/entities"

	"github.com/jmoiron/sqlx"
)

type ChestStore interface {
	CreateChest(chestType, description string, price int) (*entities.Chest, error)
}

type chestStore struct {
	db *sqlx.DB
}

func NewChestStore(db *sqlx.DB) ChestStore {
	return &chestStore{db: db}
}

func (s *chestStore) CreateChest(chestType, description string, price int) (*entities.Chest, error) {
	query := `INSERT INTO chests (type, description, price) VALUES ($1, $2, $3) RETURNING *`

	var chest entities.Chest
	err := s.db.QueryRowx(query, chestType, description, price).StructScan(&chest)
	if err != nil {
		return nil, err
	}
	return &chest, nil
}
