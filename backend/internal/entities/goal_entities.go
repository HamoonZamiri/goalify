package entities

import (
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"  json:"updated_at"`
	Title       string    `db:"title"       json:"title"`
	Description string    `db:"description" json:"description"`
	// status can be "complete" | "not_complete"
	Status     string    `db:"status"      json:"status"`
	ID         uuid.UUID `db:"id"          json:"id"`
	UserID     uuid.UUID `db:"user_id"     json:"user_id"`
	CategoryID uuid.UUID `db:"category_id" json:"category_id"`
}

type GoalCategory struct {
	CreatedAt time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt time.Time `db:"updated_at"  json:"updated_at"`
	Title     string    `db:"title"       json:"title"`
	Goals     []*Goal   `                 json:"goals"`
	XPPerGoal int       `db:"xp_per_goal" json:"xp_per_goal"`
	ID        uuid.UUID `db:"id"          json:"id"`
	UserID    uuid.UUID `db:"user_id"     json:"user_id"`
}
