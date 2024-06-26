package entities

import (
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	// status can be "complete" | "not_complete"
	Status     string    `db:"status" json:"status"`
	Id         uuid.UUID `db:"id" json:"id"`
	UserId     uuid.UUID `db:"user_id" json:"user_id"`
	CategoryId uuid.UUID `db:"category_id" json:"category_id"`
}

type GoalCategory struct {
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Title       string    `db:"title" json:"title"`
	Goals       []*Goal   `json:"goals"`
	Xp_per_goal int       `db:"xp_per_goal" json:"xp_per_goal"`
	Id          uuid.UUID `db:"id" json:"id"`
	UserId      uuid.UUID `db:"user_id" json:"user_id"`
}

// we need default constructors to ensure no nil fields are returned in json
// we had a problem where category.goals was nil, and where the slices are nil in json returns,
func NewGoal() *Goal {
	return &Goal{}
}

func NewGoalCategory() *GoalCategory {
	return &GoalCategory{Goals: []*Goal{}}
}
