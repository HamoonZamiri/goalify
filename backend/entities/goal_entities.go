package entities

import (
	"database/sql"
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

/*
This struct is used to represent the rows returned by the query below
SELECT gc.id, gc.title, gc.xp_per_goal, gc.user_id, gc.created_at, gc.updated_at,
g.id as goal_id, g.title as goal_title, g.description, g.status,
g.created_at as goal_created_at, g.updated_at as goal_updated_at
FROM goal_categories gc LEFT JOIN goals g
ON gc.id = g.category_id`

the struct below represents rows returned by this query specifically
*/
type CategoryWithGoalRow struct {
	// Category Fields
	CategoryId        uuid.UUID `db:"id"`
	CategoryUserId    uuid.UUID `db:"user_id"`
	CategoryCreatedAt time.Time `db:"created_at"`
	CategoryUpdatedAt time.Time `db:"updated_at"`
	CategoryTitle     string    `db:"title"`
	CategoryXpPerGoal int       `db:"xp_per_goal"`

	// Goals Fields
	GoalCreatedAt   sql.NullTime   `db:"goal_created_at"`
	GoalUpdatedAt   sql.NullTime   `db:"goal_updated_at"`
	GoalTitle       sql.NullString `db:"goal_title"`
	GoalDescription sql.NullString `db:"description"`
	GoalStatus      sql.NullString `db:"status"`
	GoalId          sql.NullString `db:"goal_id"`
}

func (c *CategoryWithGoalRow) ToGoalCategory() *GoalCategory {
	return &GoalCategory{
		CreatedAt:   c.CategoryCreatedAt,
		UpdatedAt:   c.CategoryUpdatedAt,
		Title:       c.CategoryTitle,
		Xp_per_goal: c.CategoryXpPerGoal,
		Id:          c.CategoryId,
		UserId:      c.CategoryUserId,
		Goals:       make([]*Goal, 0),
	}
}

func (c *CategoryWithGoalRow) ToGoal() *Goal {
	return &Goal{
		CreatedAt:   c.GoalCreatedAt.Time,
		UpdatedAt:   c.GoalUpdatedAt.Time,
		Title:       c.GoalTitle.String,
		Description: c.GoalDescription.String,
		Status:      c.GoalStatus.String,
		Id:          uuid.MustParse(c.GoalId.String),
		CategoryId:  c.CategoryId,
		UserId:      c.CategoryUserId,
	}
}
