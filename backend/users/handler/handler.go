package handler

import (
	"goalify/middleware"
	"goalify/users/service"
	"goalify/utils/jsonutil"
	"goalify/utils/responses"
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

	user, err := h.userService.Refresh(decoded.UserId, decoded.RefreshToken)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, user)
}

func (h *UserHandler) HandleUpdateUserById(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	decoded, problems, err := jsonutil.DecodeValid[UpdateRequest](r)
	if err != nil {
		responses.HandleDecodeError(w, r, problems, err)
		return
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusBadRequest, "error parsing auth header", nil)
		return
	}

	updates := make(map[string]interface{})
	if decoded.Xp.IsPresent() {
		updates["xp"] = decoded.Xp.ValueOrZero()
	}
	if decoded.LevelId.IsPresent() {
		updates["level_id"] = decoded.LevelId.ValueOrZero()
	}
	if decoded.CashAvailable.IsPresent() {
		updates["cash_available"] = decoded.CashAvailable.ValueOrZero()
	}

	if len(updates) == 0 {
		responses.SendAPIError(w, r, http.StatusBadRequest, "no updates provided", nil)
		return
	}

	user, err := h.userService.UpdateUserById(parsedUserId, updates)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, user)
}

func (h *UserHandler) GetLevelById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("levelId")
	if id == "" {
		responses.SendAPIError(w, r, http.StatusBadRequest, "url param requires id", nil)
		return
	}

	castedId, err := strconv.Atoi(id)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusBadRequest, "url param id must be an integer", nil)
		return
	}

	level, err := h.userService.GetLevelById(castedId)
	if err != nil {
		responses.SendAPIError(w, r, responses.GetErrorCode(err), err.Error(), nil)
		return
	}

	responses.SendResponse(w, r, http.StatusOK, level)
}
