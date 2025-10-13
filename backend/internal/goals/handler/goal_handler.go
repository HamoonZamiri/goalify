package handler

import (
	"fmt"
	"goalify/internal/middleware"
	"goalify/pkg/jsonutil"
	"goalify/internal/responses"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *GoalHandler) HandleCreateGoal(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleCreateGoal")
	body, problems, err := jsonutil.DecodeValid[CreateGoalRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}

	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing user id", nil)
		return
	}

	parsedCategoryId, err := uuid.Parse(body.CategoryId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing category id", nil)
		return
	}

	goal, err := h.goalService.CreateGoal(body.Title, body.Description, parsedUserId, parsedCategoryId)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusCreated, goal)
}

func (h *GoalHandler) HandleUpdateGoalById(w http.ResponseWriter, r *http.Request) {
	body, problems, err := jsonutil.DecodeValid[UpdateGoalRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}

	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error("HandleUpdateGoalById: middleware.GetIdFromHeader: ", "err", err)
		responses.SendInternalServerError(w, r)
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		slog.Error("HandleUpdateGoalById: uuid.Parse(userId):", "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	parsedGoalId, err := uuid.Parse(r.PathValue("goalId"))
	if err != nil {
		slog.Error("HandleUpdateGoalById: uuid.Parse(goalId): ", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "invalid goal id", nil)
		return
	}

	updates := make(map[string]interface{})
	if body.Title.IsPresent() {
		updates["title"] = body.Title.ValueOrZero()
	}
	if body.Description.IsPresent() {
		updates["description"] = body.Description.ValueOrZero()
	}
	if body.CategoryId.IsPresent() {
		updates["category_id"] = body.CategoryId.ValueOrZero
	}
	if body.Status.IsPresent() {
		updates["status"] = body.Status.ValueOrZero()
	}
	if len(updates) == 0 {
		responses.SendAPIError(w, r, http.StatusBadRequest, "no updates provided", nil)
		return
	}

	updatedGoal, err := h.goalService.UpdateGoalById(parsedGoalId, updates, parsedUserId)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, updatedGoal)
}

func (gh *GoalHandler) HandleDeleteGoalById(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error("HandleDeleteGoalById: middleware.GetIdFromHeader: ", "err", err)
		responses.SendInternalServerError(w, r)
		return
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error("HandleDeleteGoalById: uuid.Parse:", "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	goalId := r.PathValue("goalId")
	goalUUID, err := uuid.Parse(goalId)
	if err != nil {
		slog.Error("HandleDeleteGoalById: uuid.Parse:", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "bad request: invalid goal id", nil)
		return
	}

	err = gh.goalService.DeleteGoalById(goalUUID, userUUID)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	res := map[string]any{"id": goalId, "deleted": true}
	responses.SendResponse(w, r, http.StatusOK, res)
}
