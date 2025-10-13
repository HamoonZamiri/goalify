package handler

import (
	"fmt"
	"goalify/internal/entities"
	"goalify/internal/middleware"
	"goalify/pkg/jsonutil"
	"goalify/internal/responses"
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
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusUnauthorized, "user is not authenticated", nil)
		return
	}

	parsedUUID, err := uuid.Parse(userId)
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

func (h *GoalHandler) HandleGetGoalCategoriesByUserId(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleGetGoalCategoriesByUserId")
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	cats, err := h.goalService.GetGoalCategoriesByUserId(parsedUUID)
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

func (h *GoalHandler) HandleGetGoalCategoryById(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleGetGoalCategoryById")
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendInternalServerError(w, r)
		return
	}

	categoryId := r.PathValue("categoryId")
	parsedCategoryId, err := uuid.Parse(categoryId)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusBadRequest, "bad request: invalid category id", nil)
		return
	}

	cat, err := h.goalService.GetGoalCategoryById(parsedCategoryId, parsedUUID)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, cat)
}

func (h *GoalHandler) HandleUpdateGoalCategoryById(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleUpdateGoalCategoryById")
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	categoryId := r.PathValue("categoryId")
	parsedCategoryId, err := uuid.Parse(categoryId)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusBadRequest, "bad request: invalid category id", nil)
		return
	}

	body, problems, err := jsonutil.DecodeValid[UpdateGoalCategoryRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}

	updates := make(map[string]any)
	if body.Title.IsPresent() {
		updates["title"] = body.Title.ValueOrZero()
	}

	if body.XpPerGoal.IsPresent() {
		updates["xp_per_goal"] = body.XpPerGoal.ValueOrZero()
	}

	if len(updates) == 0 {
		noUpdatesError := fmt.Errorf("%w: no fields given to update", responses.ErrBadRequest)
		responses.SendAPIError(w, r, http.StatusBadRequest, noUpdatesError.Error(), nil)
		return
	}

	cat, err := h.goalService.UpdateGoalCategoryById(parsedCategoryId, updates, parsedUUID)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, cat)
}

func (h *GoalHandler) HandleDeleteGoalCategoryById(w http.ResponseWriter, r *http.Request) {
	funcStr := h.traceLogger.GetTrace("handler.HandleDeleteGoalCategoryById")
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	categoryId := r.PathValue("categoryId")
	parsedCategoryId, err := uuid.Parse(categoryId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing /{categoryId} param", nil)
		return
	}

	err = h.goalService.DeleteGoalCategoryById(parsedCategoryId, parsedUUID)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: goalService.DeleteGoalCategoryById:", funcStr), "err", err)
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	res := map[string]any{"id": categoryId, "deleted": true}
	responses.SendResponse[map[string]any](w, r, http.StatusOK, res)
}
