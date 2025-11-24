// Package handler is the API layer request response handling package for the users domain
package handler

import (
	"goalify/internal/middleware"
	"goalify/internal/responses"
	"goalify/internal/users/service"
	"goalify/pkg/jsonutil"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	decoded, problems, err := jsonutil.DecodeValid[SignupRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}

	user, err := h.userService.SignUp(decoded.Email, decoded.Password)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusCreated, user)
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoded, problems, err := jsonutil.DecodeValid[LoginRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}

	user, err := h.userService.Login(decoded.Email, decoded.Password)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, user)
}

func (h *UserHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	decoded, problems, err := jsonutil.DecodeValid[RefreshRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}

	user, err := h.userService.Refresh(decoded.UserID, decoded.RefreshToken)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, user)
}

func (h *UserHandler) HandleUpdateUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetIDFromHeader(r)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	decoded, problems, err := jsonutil.DecodeValid[UpdateRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing auth header", nil)
		return
	}

	updates := make(map[string]any)
	if decoded.Xp.IsPresent() {
		updates["xp"] = decoded.Xp.ValueOrZero()
	}
	if decoded.LevelID.IsPresent() {
		updates["level_id"] = decoded.LevelID.ValueOrZero()
	}
	if decoded.CashAvailable.IsPresent() {
		updates["cash_available"] = decoded.CashAvailable.ValueOrZero()
	}

	if len(updates) == 0 {
		responses.SendAPIError(w, r, http.StatusBadRequest, "bad request: no updates provided", nil)
		return
	}

	user, err := h.userService.UpdateUserByID(parsedUserID, updates)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, user)
}

func (h *UserHandler) GetLevelByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("levelId")
	if id == "" {
		responses.SendAPIError(w, r, http.StatusBadRequest, "url param requires id", nil)
		return
	}

	castedID, err := strconv.Atoi(id)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusBadRequest, "url param id must be an integer", nil)
		return
	}

	level, err := h.userService.GetLevelByID(castedID)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, level)
}
