package handler

import (
	"fmt"
	"goalify/internal/goals/service"
	"goalify/pkg/options"
	"goalify/pkg/stacktrace"

	"github.com/google/uuid"
)

type (
	GoalHandler struct {
		goalService service.GoalService
		traceLogger stacktrace.TraceLogger
	}
	CreateGoalRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
	}
	CreateGoalCategoryRequest struct {
		Title     string `json:"title"`
		XpPerGoal int    `json:"xp_per_goal"`
	}
	UpdateGoalCategoryRequest struct {
		Title     options.Option[string] `json:"title"`
		XpPerGoal options.Option[int]    `json:"xp_per_goal"`
	}
	UpdateGoalRequest struct {
		Title       options.Option[string] `json:"title"`
		Description options.Option[string] `json:"description"`
		CategoryID  options.Option[string] `json:"category_id"`
		Status      options.Option[string] `json:"status"`
	}
	DeleteGoalRequest struct {
		GoalID string `json:"goal_id"`
	}
)

func NewGoalHandler(
	goalService service.GoalService,
	traceLogger stacktrace.TraceLogger,
) *GoalHandler {
	return &GoalHandler{goalService, traceLogger}
}

const (
	TextMaxLen   = 255
	XPMaxPerGoal = 100
)

func NewGoalCategoryRequest(title string, xpPerGoal int) CreateGoalCategoryRequest {
	return CreateGoalCategoryRequest{title, xpPerGoal}
}

func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func (r CreateGoalCategoryRequest) Valid() map[string]string {
	problems := make(map[string]string)

	if r.Title == "" {
		problems["title"] = "title is required"
	} else if len(r.Title) > TextMaxLen {
		problems["title"] = "title must be less than 255 characters"
	}

	if r.XpPerGoal <= 0 {
		problems["xp_per_goal"] = "xp per goal must be greater than 0"
	} else if r.XpPerGoal > XPMaxPerGoal {
		problems["xp_per_goal"] = fmt.Sprintf("xp per goal must be less than %d", XPMaxPerGoal)
	}
	return problems
}

func (r UpdateGoalCategoryRequest) Valid() map[string]string {
	problems := make(map[string]string)
	if r.Title.IsPresent() && r.Title.ValueOrZero() == "" {
		problems["title"] = "title cannot be empty"
	}

	if r.XpPerGoal.IsPresent() && r.XpPerGoal.ValueOrZero() <= 0 {
		problems["xp_per_goal"] = "xp per goal must be greater than 0"
	}

	if r.XpPerGoal.IsPresent() && r.XpPerGoal.ValueOrZero() > XPMaxPerGoal {
		problems["xp_per_goal"] = fmt.Sprintf("xp per goal must be less than %d", XPMaxPerGoal)
	}

	return problems
}

func (r CreateGoalRequest) Valid() map[string]string {
	problems := make(map[string]string)

	if r.Title == "" {
		problems["title"] = "title is required"
	}
	if len(r.Title) > TextMaxLen {
		problems["title"] = "title must be less than 255 characters"
	}

	if r.Description == "" {
		problems["description"] = "description is required"
	}
	if len(r.Description) > TextMaxLen {
		problems["description"] = "description must be less than 255 characters"
	}

	if r.CategoryID == "" {
		problems["category_id"] = "category id is required"
	}

	return problems
}

func (r UpdateGoalRequest) Valid() map[string]string {
	problems := make(map[string]string)

	if r.Title.IsPresent() && r.Title.ValueOrZero() == "" {
		problems["title"] = "title cannot be empty"
	}

	if r.Title.IsPresent() && len(r.Title.ValueOrZero()) > TextMaxLen {
		problems["title"] = "title must be less than 255 characters"
	}

	if r.Description.IsPresent() && len(r.Description.ValueOrZero()) > TextMaxLen {
		problems["description"] = "description must be less than 255 characters"
	}

	if r.Description.IsPresent() && r.Description.ValueOrZero() == "" {
		problems["description"] = "description cannot be empty"
	}

	if r.CategoryID.IsPresent() && r.CategoryID.ValueOrZero() == "" {
		problems["category_id"] = "category id cannot be empty"
	}

	if r.CategoryID.IsPresent() && !isValidUUID(r.CategoryID.ValueOrZero()) {
		problems["category_id"] = "category id must be a valid UUID"
	}

	if r.Status.IsPresent() && r.Status.ValueOrZero() != "complete" &&
		r.Status.ValueOrZero() != "not_complete" {
		problems["status"] = "status must be either 'complete' or 'not_complete'"
	}
	return problems
}
