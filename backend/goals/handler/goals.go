package handler

import (
	"fmt"
	"goalify/middleware"
	"goalify/utils/jsonutil"
	"goalify/utils/responses"
	"goalify/utils/svcerror"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *GoalHandler) HandleCreateGoal(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleCreateGoal")
	body, problems, err := jsonutil.DecodeValid[CreateGoalRequest](r)
	if len(problems) > 0 {
		apiError := responses.NewAPIError("invalid request", problems)
		jsonutil.Encode(w, r, http.StatusBadRequest, apiError)
		return
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: jsonutil.DecodeValid[CreateGoalRequest]:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "error decoding request body", nil)
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
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	if err := jsonutil.Encode(w, r, http.StatusOK, responses.New(goal, "goal created successfully")); err != nil {
		slog.Error("HandleCreateGoal: jsonutil.Encode(goal):", "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "internal error encoding response", nil)
	}
}

func (h *GoalHandler) HandleUpdateGoalById(w http.ResponseWriter, r *http.Request) {
	body, problems, err := jsonutil.DecodeValid[UpdateGoalRequest](r)
	if len(problems) > 0 {
		apiError := responses.NewAPIError("invalid request", problems)
		jsonutil.Encode(w, r, http.StatusBadRequest, apiError)
		return
	}

	if err != nil {
		slog.Error("HandleUpdateGoalById: jsonutil.DecodeValid[UpdateGoalRequest]():", "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "error decoding request body", nil)
		return
	}

	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error("HandleUpdateGoalById: middleware.GetIdFromHeader: ", "err", err)
		responses.SendAPIError(w, r, http.StatusUnauthorized, err.Error(), nil)
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		slog.Error("HandleUpdateGoalById: uuid.Parse(userId):", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing user id", nil)
		return
	}

	parsedGoalId, err := uuid.Parse(r.PathValue("goalId"))
	if err != nil {
		slog.Error("HandleUpdateGoalById: uuid.Parse(goalId): ", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing goal id", nil)
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
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	if err := jsonutil.Encode(w, r, http.StatusOK, responses.New(updatedGoal, "goal updated successfully")); err != nil {
		slog.Error("HandleUpdateGoalById: jsonutil.Encode(updatedGoal):", "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "internal error encoding response", nil)
	}
}

func (gh *GoalHandler) HandleDeleteGoalById(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error("HandleDeleteGoalById: middleware.GetIdFromHeader: ", "err", err)
		responses.SendAPIError(w, r, http.StatusUnauthorized, err.Error(), nil)
		return
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error("HandleDeleteGoalById: uuid.Parse:", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing user id", nil)
		return
	}

	goalId := r.PathValue("goalId")
	goalUUID, err := uuid.Parse(goalId)
	if err != nil {
		slog.Error("HandleDeleteGoalById: uuid.Parse:", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing goal id", nil)
		return
	}

	err = gh.goalService.DeleteGoalById(goalUUID, userUUID)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusInternalServerError, "error deleting goal", nil)
		return
	}

	res := map[string]string{"data": "null", "message": "goal deleted successfully"}
	if err := jsonutil.Encode(w, r, http.StatusOK, res); err != nil {
		slog.Error("HandleDeleteGoalById: jsonutil.Encode(nil):", "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "internal error encoding response", nil)
	}
}
