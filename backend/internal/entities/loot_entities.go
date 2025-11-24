package entities

import (
	"goalify/pkg/options"
	"time"

	"github.com/google/uuid"
)

type Chest struct {
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"  json:"updated_at"`
	Type        string    `db:"type"        json:"type"`
	Description string    `db:"description" json:"description"`
	Price       int       `db:"price"       json:"price"`
	ID          uuid.UUID `db:"id"          json:"id"`
}

type ChestItem struct {
	CreatedAt time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt time.Time           `db:"updated_at" json:"updated_at"`
	ImageURL  string              `db:"image_url"  json:"image_url"`
	Rarity    string              `db:"rarity"     json:"rarity"`
	Price     options.Option[int] `db:"price"      json:"price"`
	ID        uuid.UUID           `db:"id"         json:"id"`
}

type ChestItemDropRate struct {
	ItemID   uuid.UUID `db:"item_id"   json:"item_id"`
	ChestID  uuid.UUID `db:"chest_id"  json:"chest_id"`
	DropRate float32   `db:"drop_rate" json:"drop_rate"`
}

type UserChest struct {
	UserID        uuid.UUID `db:"user_id"        json:"user_id"`
	ChestID       uuid.UUID `db:"chest_id"       json:"chest_id"`
	QuantityOwned int       `db:"quantity_owned" json:"quantity_owned"`
}

type UserItem struct {
	Status string    `db:"status"  json:"status"`
	UserID uuid.UUID `db:"user_id" json:"user_id"`
	ItemID uuid.UUID `db:"item_id" json:"item_id"`
}
