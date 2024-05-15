package handler

import (
	"goalify/goals/service"
	"goalify/jsonutil"
	"goalify/middleware"
	"goalify/responses"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type GoalHandler struct {
	goalService service.GoalService
}

func NewGoalHandler(goalService service.GoalService) *GoalHandler {
	return &GoalHandler{goalService}
}

func (h *GoalHandler) HandleCreateGoal(w http.ResponseWriter, r *http.Request) {
}

func (h *GoalHandler) HandleCreateGoalCategory(w http.ResponseWriter, r *http.Request) {
	body, problems, err := jsonutil.DecodeValid[CreateGoalCategoryRequest](r)
	if len(problems) > 0 {
		apiError := responses.NewAPIError("invalid request", problems)
		jsonutil.Encode(w, r, http.StatusBadRequest, apiError)
	}

	if err != nil {
		slog.Error("decode valid:", "err", err)
		apiError := responses.NewAPIError("error decoding request body", nil)
		jsonutil.Encode(w, r, http.StatusInternalServerError, apiError)
	}

	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error("create goal category auth: ", "err", err)
		apiError := responses.NewAPIError("user is not authenticated", nil)
		jsonutil.Encode(w, r, http.StatusUnauthorized, apiError)
	}

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error("parse user id: ", "err", err)
		apiError := responses.NewAPIError("error parsing user id", nil)
		jsonutil.Encode(w, r, http.StatusInternalServerError, apiError)
	}

	category, err := h.goalService.CreateGoalCategory(body.Title, body.XpPerGoal, parsedUUID)
	if err != nil {
		slog.Error("create goal category: ", "err", err)
		apiError := responses.NewAPIError(err.Error(), nil)
		jsonutil.Encode(w, r, http.StatusInternalServerError, apiError)
	}

	res := responses.New(category, "goal category created successfully")
	if err := jsonutil.Encode(w, r, http.StatusOK, res); err != nil {
		slog.Error("json encode: ", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
