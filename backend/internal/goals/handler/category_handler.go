// Package handler is the API request/response handling for goals and goal categories
package handler

import (
	"fmt"
	"goalify/internal/entities"
	"goalify/internal/goals/stores"
	"goalify/internal/middleware"
	"goalify/internal/responses"
	"goalify/pkg/jsonutil"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *GoalHandler) HandleCreateGoalCategory(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleCreateGoalCategory")
	body, problems, err := jsonutil.DecodeValid[CreateGoalCategoryRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}
	userID, err := middleware.GetIDFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusUnauthorized, "user is not authenticated", nil)
		return
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "error parsing user id", nil)
		return
	}

	category, err := h.goalService.CreateGoalCategory(body.Title, body.XpPerGoal, parsedUUID)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusCreated, category)
}

func (h *GoalHandler) HandleGetGoalCategoriesByUserID(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleGetGoalCategoriesByUserId")
	userID, err := middleware.GetIDFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	cats, err := h.goalService.GetGoalCategoriesByUserID(parsedUUID)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res := responses.ServerResponse[[]*entities.GoalCategory]{
		Object: responses.ObjectList,
		Data:   cats,
	}
	responses.SendResponse(w, r, http.StatusOK, res)
}

func (h *GoalHandler) HandleGetGoalCategoryByID(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleGetGoalCategoryById")
	userID, err := middleware.GetIDFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	categoryID := r.PathValue("categoryId")
	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusBadRequest, "bad request: invalid category id", nil)
		return
	}

	cat, err := h.goalService.GetGoalCategoryByID(parsedCategoryID, parsedUUID)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, cat)
}

func (h *GoalHandler) HandleUpdateGoalCategoryByID(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleUpdateGoalCategoryById")
	userID, err := middleware.GetIDFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	categoryID := r.PathValue("categoryId")
	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusBadRequest, "bad request: invalid category id", nil)
		return
	}

	body, problems, err := jsonutil.DecodeValid[UpdateGoalCategoryRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}

	params := stores.UpdateGoalCategoryParams{
		Title:     body.Title,
		XpPerGoal: body.XpPerGoal,
	}

	if !params.Title.IsPresent() && !params.XpPerGoal.IsPresent() {
		noUpdatesError := fmt.Errorf("%w: no fields given to update", responses.ErrBadRequest)
		responses.SendAPIError(w, r, http.StatusBadRequest, noUpdatesError.Error(), nil)
		return
	}

	cat, err := h.goalService.UpdateGoalCategoryByID(parsedCategoryID, params, parsedUUID)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, cat)
}

func (h *GoalHandler) HandleDeleteGoalCategoryByID(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleDeleteGoalCategoryById")
	userID, err := middleware.GetIDFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	categoryID := r.PathValue("categoryId")
	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(
			w,
			r,
			http.StatusBadRequest,
			"error parsing /{categoryId} param",
			nil,
		)
		return
	}

	err = h.goalService.DeleteGoalCategoryByID(parsedCategoryID, parsedUUID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: goalService.DeleteGoalCategoryById:", funcStr), "err", err)
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	res := map[string]any{"id": categoryID, "deleted": true}
	responses.SendResponse(w, r, http.StatusOK, res)
}

func (h *GoalHandler) HandleResetGoalsByCategoryID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetIDFromHeader(r)
	if err != nil {
		slog.Error("HandleResetGoalsByCategoryID: middleware.GetIdFromHeader: ", "err", err)
		responses.SendInternalServerError(w, r)
		return
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		slog.Error("HandleResetGoalsByCategoryID: uuid.Parse:", "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	categoryID := r.PathValue("categoryID")
	parsedCategoryID, err := uuid.Parse(categoryID)
	if err != nil {
		slog.Error("HandleResetGoalsByCategoryID: uuid.Parse:", "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "bad request: invalid category id", nil)
		return
	}

	err = h.goalService.ResetGoalsByCategoryID(parsedCategoryID, userUUID)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusNoContent, map[string]any{})
}
