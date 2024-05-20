package handler

import (
	"fmt"
	"goalify/options"
)

type CreateGoalRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryId  string `json:"category_id"`
}

type CreateGoalCategoryRequest struct {
	Title     string `json:"title"`
	XpPerGoal int    `json:"xp_per_goal"`
}

type UpdateGoalCategoryRequest struct {
	Title     options.Option[string] `json:"title"`
	XpPerGoal options.Option[int]    `json:"xp_per_goal"`
}

const (
	TEXT_MAX_LEN    = 255
	XP_MAX_PER_GOAL = 100
)

func NewGoalCategoryRequest(title string, xpPerGoal int) CreateGoalCategoryRequest {
	return CreateGoalCategoryRequest{title, xpPerGoal}
}

func (r CreateGoalCategoryRequest) Valid() map[string]string {
	problems := make(map[string]string)

	if r.Title == "" {
		problems["title"] = "title is required"
	} else if len(r.Title) > TEXT_MAX_LEN {
		problems["title"] = "title must be less than 255 characters"
	}

	if r.XpPerGoal <= 0 {
		problems["xp_per_goal"] = "xp_per_goal must be greater than 0"
	} else if r.XpPerGoal > XP_MAX_PER_GOAL {
		problems["xp_per_goal"] = fmt.Sprintf("xp_per_goal must be less than %d", XP_MAX_PER_GOAL)
	}
	return problems
}

func (r UpdateGoalCategoryRequest) Valid() map[string]string {
	problems := make(map[string]string)
	if r.Title.IsPresent() && r.Title.ValueOrZero() == "" {
		problems["title"] = "title cannot be empty"
	}

	if r.XpPerGoal.IsPresent() && r.XpPerGoal.ValueOrZero() <= 0 {
		problems["xp_per_goal"] = "xp_per_goal must be greater than 0"
	}

	if r.XpPerGoal.IsPresent() && r.XpPerGoal.ValueOrZero() > XP_MAX_PER_GOAL {
		problems["xp_per_goal"] = fmt.Sprintf("xp_per_goal must be less than %d", XP_MAX_PER_GOAL)
	}

	return problems
}

func (r CreateGoalRequest) Valid() map[string]string {
	problems := make(map[string]string)

	if r.Title == "" {
		problems["title"] = "title is required"
	}
	if len(r.Title) > TEXT_MAX_LEN {
		problems["title"] = "title must be less than 255 characters"
	}

	if r.Description == "" {
		problems["description"] = "description is required"
	}
	if len(r.Description) > TEXT_MAX_LEN {
		problems["description"] = "description must be less than 255 characters"
	}

	if r.CategoryId == "" {
		problems["category_id"] = "category_id is required"
	}

	return problems
}
