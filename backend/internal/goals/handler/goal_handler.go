package handler

import (
	"fmt"
	"goalify/internal/goals/stores"
	"goalify/internal/middleware"
	"goalify/internal/responses"
	"goalify/pkg/jsonutil"
	"goalify/pkg/options"
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

	userID, err := middleware.GetIDFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing user id", nil)
		return
	}

	parsedCategoryID, err := uuid.Parse(body.CategoryID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing category id", nil)
		return
	}

	goal, err := h.goalService.CreateGoal(
		body.Title,
		body.Description,
		parsedUserID,
		parsedCategoryID,
	)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusCreated, goal)
}

func (h *GoalHandler) HandleUpdateGoalByID(w http.ResponseWriter, r *http.Request) {
	body, problems, err := jsonutil.DecodeValid[UpdateGoalRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}

	userID, err := middleware.GetIDFromHeader(r)
	if err != nil {
		slog.Error("HandleUpdateGoalById: middleware.GetIdFromHeader: ", "err", err)
		responses.SendInternalServerError(w, r)
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error("HandleUpdateGoalById: uuid.Parse(userId):", "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	parsedGoalID, err := uuid.Parse(r.PathValue("goalId"))
	if err != nil {
		slog.Error("HandleUpdateGoalById: uuid.Parse(goalId): ", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "invalid goal id", nil)
		return
	}

	var params stores.UpdateGoalParams
	params.Title = body.Title
	params.Description = body.Description
	params.Status = body.Status

	if body.CategoryID.IsPresent() {
		var parsedCategoryID uuid.UUID
		parsedCategoryID, err = uuid.Parse(body.CategoryID.ValueOrZero())
		if err != nil {
			responses.SendAPIError(w, r, http.StatusBadRequest, "invalid category id", nil)
			return
		}
		params.CategoryID = options.Some(parsedCategoryID)
	}

	if !params.Title.IsPresent() && !params.Description.IsPresent() &&
		!params.CategoryID.IsPresent() && !params.Status.IsPresent() {
		responses.SendAPIError(w, r, http.StatusBadRequest, "no updates provided", nil)
		return
	}

	updatedGoal, err := h.goalService.UpdateGoalByID(parsedGoalID, params, parsedUserID)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, updatedGoal)
}

func (h *GoalHandler) HandleDeleteGoalByID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetIDFromHeader(r)
	if err != nil {
		slog.Error("HandleDeleteGoalById: middleware.GetIdFromHeader: ", "err", err)
		responses.SendInternalServerError(w, r)
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error("HandleDeleteGoalById: uuid.Parse:", "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	goalID := r.PathValue("goalId")
	goalUUID, err := uuid.Parse(goalID)
	if err != nil {
		slog.Error("HandleDeleteGoalById: uuid.Parse:", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "bad request: invalid goal id", nil)
		return
	}

	err = h.goalService.DeleteGoalByID(goalUUID, userUUID)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	res := map[string]any{"id": goalID, "deleted": true}
	responses.SendResponse(w, r, http.StatusOK, res)
}
