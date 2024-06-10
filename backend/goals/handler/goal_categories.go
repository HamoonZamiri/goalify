package handler

import (
	"fmt"
	"goalify/jsonutil"
	"goalify/middleware"
	"goalify/responses"
	"goalify/svcerror"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *GoalHandler) HandleCreateGoalCategory(w http.ResponseWriter, r *http.Request) {
	funcStr := "Goal Categories: handler.HandleCreateGoalCategory"
	body, problems, err := jsonutil.DecodeValid[CreateGoalCategoryRequest](r)
	if len(problems) > 0 {
		apiError := responses.NewAPIError("invalid request", problems)
		jsonutil.Encode(w, r, http.StatusBadRequest, apiError)
		return
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: jsontutil.DecodeValid:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "error decoding request body", nil)
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
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	res := responses.New(category, "goal category created successfully")
	if err := jsonutil.Encode(w, r, http.StatusOK, res); err != nil {
		slog.Error(fmt.Sprintf("%s: jsonutil.Encode:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "internal error encoding response", nil)
	}
}

func (h *GoalHandler) HandleGetGoalCategoriesByUserId(w http.ResponseWriter, r *http.Request) {
	funcStr := "Goal Categories: handler.HandleGetGoalCategoriesByUserId"
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "error parsing user id", nil)
		return
	}

	cats, err := h.goalService.GetGoalCategoriesByUserId(parsedUUID)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	res := responses.New(cats, "goal categories retrieved successfully")
	if err := jsonutil.Encode(w, r, http.StatusOK, res); err != nil {
		slog.Error(fmt.Sprintf("%s: jsonutil.Encode:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "internal error encoding response", nil)
	}
}

func (h *GoalHandler) HandleGetGoalCategoryById(w http.ResponseWriter, r *http.Request) {
	funcStr := "Goal Categories: handler.HandleGetGoalCategoryById"
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "error parsing user id", nil)
		return
	}

	categoryId := r.PathValue("categoryId")
	parsedCategoryId, err := uuid.Parse(categoryId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "error parsing category id", nil)
		return
	}

	cat, err := h.goalService.GetGoalCategoryById(parsedCategoryId, parsedUUID)
	if err != nil {
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	res := responses.New(cat, "goal category retrieved successfully")
	if err := jsonutil.Encode(w, r, http.StatusOK, res); err != nil {
		slog.Error(fmt.Sprintf("%s: jsonutil.Encode:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "internal error encoding response", nil)
	}
}

func (h *GoalHandler) HandleUpdateGoalCategoryById(w http.ResponseWriter, r *http.Request) {
	funcStr := "Goal Categories: handler.HandleUpdateGoalCategoryById"
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
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	body, problems, err := jsonutil.DecodeValid[UpdateGoalCategoryRequest](r)
	if len(problems) > 0 {
		apiError := responses.NewAPIError("invalid request", problems)
		jsonutil.Encode(w, r, http.StatusBadRequest, apiError)
		return
	}

	if err != nil {
		slog.Error(fmt.Sprintf("%s: jsonutil.DecodeValid:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
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
		responses.SendAPIError(w, r, http.StatusBadRequest, "no fields given to update", nil)
		return
	}

	cat, err := h.goalService.UpdateGoalCategoryById(parsedCategoryId, updates, parsedUUID)
	if err != nil {
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
	}

	if err := jsonutil.Encode(w, r, http.StatusOK, responses.New(cat, "goal category updated successfully")); err != nil {
		slog.Error(fmt.Sprintf("%s: jsonutil.Encode:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
	}
}

func (h *GoalHandler) HandleDeleteGoalCategoryById(w http.ResponseWriter, r *http.Request) {
	funcStr := "Goal Categories: handler.HandleDeleteGoalCategoryById"
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: middleware.GetIdFromHeader:", funcStr), "err", err)
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	parsedUUID, err := uuid.Parse(userId)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: uuid.Parse:", funcStr), "err", err)
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
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
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	if err := jsonutil.Encode(w, r, http.StatusOK, responses.New[string]("null", "goal category deleted successfully")); err != nil {
		slog.Error(fmt.Sprintf("%s: jsonutil.Encode:", funcStr), "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
	}
}
