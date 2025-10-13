package entities

import (
	"goalify/pkg/options"
	"time"

	"github.com/google/uuid"
)

type Chest struct {
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Id          uuid.UUID `db:"id" json:"id"`
	Type        string    `db:"type" json:"type"` // bronze | silver | gold
	Description string    `db:"description" json:"description"`
	Price       int       `db:"price" json:"price"`
}

type ChestItem struct {
	CreatedAt time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt time.Time           `db:"updated_at" json:"updated_at"`
	Id        uuid.UUID           `db:"id" json:"id"`
	ImageUrl  string              `db:"image_url" json:"image_url"`
	Rarity    string              `db:"rarity" json:"rarity"` // common | rare | epic | legendary
	Price     options.Option[int] `db:"price" json:"price"`
}

type ChestItemDropRate struct {
	ItemId   uuid.UUID `db:"item_id" json:"item_id"`
	ChestId  uuid.UUID `db:"chest_id" json:"chest_id"`
	DropRate float32   `db:"drop_rate" json:"drop_rate"`
}

type UserChest struct {
	UserId        uuid.UUID `db:"user_id" json:"user_id"`
	ChestId       uuid.UUID `db:"chest_id" json:"chest_id"`
	QuantityOwned int       `db:"quantity_owned" json:"quantity_owned"`
}

type UserItem struct {
	UserId uuid.UUID `db:"user_id" json:"user_id"`
	ItemId uuid.UUID `db:"item_id" json:"item_id"`
	Status string    `db:"status" json:"status"` // equipped | unequipped
}
