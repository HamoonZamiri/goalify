package handler

type CreateGoalRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryId  string `json:"category_id"`
}

type CreateGoalCategoryRequest struct {
	Title     string `json:"title"`
	XpPerGoal int    `json:"xp_per_goal"`
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
		problems["xp_per_goal"] = "xp_per_goal must be less than 100"
	}
	return problems
}
