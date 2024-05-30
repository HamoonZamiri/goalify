package handler

import (
	"goalify/jsonutil"
	"goalify/middleware"
	"goalify/responses"
	"goalify/svcerror"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *GoalHandler) HandleCreateGoal(w http.ResponseWriter, r *http.Request) {
	body, problems, err := jsonutil.DecodeValid[CreateGoalRequest](r)
	if len(problems) > 0 {
		apiError := responses.NewAPIError("invalid request", problems)
		jsonutil.Encode(w, r, http.StatusBadRequest, apiError)
		return
	}

	if err != nil {
		slog.Error("decode valid:", "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "error decoding request body", nil)
		return
	}

	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error("middleware.GetIdFromHeader: ", "err", err)
		responses.SendAPIError(w, r, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		slog.Error("parse userId: ", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing user id", nil)
		return
	}

	parsedCategoryId, err := uuid.Parse(body.CategoryId)
	if err != nil {
		slog.Error("parse categoryId: ", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing category id", nil)
		return
	}

	goal, err := h.goalService.CreateGoal(body.Title, body.Description, parsedUserId, parsedCategoryId)
	if err != nil {
		slog.Error("create goal: ", "err", err)
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	if err := jsonutil.Encode(w, r, http.StatusOK, responses.New(goal, "goal created successfully")); err != nil {
		slog.Error("json encode: ", "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "internal error encoding response", nil)
	}
}
