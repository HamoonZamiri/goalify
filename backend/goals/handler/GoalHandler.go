package handler

import (
	"goalify/goals/service"
	"net/http"
)

type GoalHandler struct {
	goalService service.GoalService
}

func NewGoalHandler(goalService service.GoalService) *GoalHandler {
	return &GoalHandler{goalService}
}

func (h *GoalHandler) HandleCreateGoal(w http.ResponseWriter, r *http.Request) {
}
