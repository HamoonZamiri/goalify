package handler

import (
	"goalify/middleware"
	"goalify/users/service"
	"goalify/utils/jsonutil"
	"goalify/utils/responses"
	"goalify/utils/svcerror"
	"log/slog"
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
	if len(problems) > 0 {
		responses.SendAPIError(w, r, http.StatusUnprocessableEntity, "invalid request body", problems)
		return
	}

	if err != nil {
		slog.Error("handler.HandleLogin: jsonutil.DecodeValid:", "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "internal error decoding request body", nil)
		return
	}

	user, err := h.userService.SignUp(decoded.Email, decoded.Password)
	if err != nil {
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	res := responses.New(user, "user created")
	if err := jsonutil.Encode(w, r, http.StatusOK, *res); err != nil {
		slog.Error("json encode: ", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoded, err := jsonutil.Decode[LoginRequest](r)
	if err != nil {
		slog.Error("handler.HandleLogin: jsonutil.DecodeValid:", "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "internal error decoding request body", nil)
		return
	}

	user, err := h.userService.Login(decoded.Email, decoded.Password)
	if err != nil {
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	if err := jsonutil.Encode(w, r, http.StatusOK, responses.New(user, "user logged in successfully")); err != nil {
		slog.Error("handler.HandleLogin: jsonutil.Encode:", "err", err)
		responses.SendAPIError(w, r, http.StatusInternalServerError, "internal error", nil)
		return
	}
}

func (h *UserHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	decoded, err := jsonutil.Decode[RefreshRequest](r)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	user, err := h.userService.Refresh(decoded.UserId, decoded.RefreshToken)
	if err != nil {
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	servResp := responses.New(user, "user refreshed")

	if err := jsonutil.Encode(w, r, http.StatusOK, servResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("json encode: ", "err", err)
		return
	}
}

func (h *UserHandler) HandleUpdateUserById(w http.ResponseWriter, r *http.Request) {
	userId, err := middleware.GetIdFromHeader(r)
	if err != nil {
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	decoded, err := jsonutil.Decode[UpdateRequest](r)
	if err != nil {
		responses.SendAPIError(w, r, http.StatusBadRequest, err.Error(), nil)
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
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	res := responses.New(user, "successfully updated user")
	if err := jsonutil.Encode(w, r, http.StatusOK, res); err != nil {
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
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
		responses.SendAPIError(w, r, svcerror.GetErrorCode(err), err.Error(), nil)
		return
	}

	res := responses.New(level, "level retrieved sucessfully")
	if err := jsonutil.Encode(w, r, http.StatusOK, res); err != nil {
		responses.SendAPIError(w, r, http.StatusInternalServerError, err.Error(), nil)
		return
	}
}
