package handler

import (
	"errors"
	"goalify/jsonutil"
	"goalify/responses"
	"goalify/svcerror"
	"goalify/users/service"
	"log/slog"
	"net/http"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest = SignupRequest

type RefreshRequest struct {
	UserId       string `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func getErrorCode(err error) int {
	if errors.Is(err, svcerror.ErrBadRequest) {
		return http.StatusBadRequest
	}
	if errors.Is(err, svcerror.ErrNotFound) {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}

func (h *UserHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	decoded, err := jsonutil.Decode[SignupRequest](r)
	if err != nil {
		http.Error(w, err.Error(), getErrorCode(err))
		return
	}

	user, err := h.userService.SignUp(decoded.Email, decoded.Password)
	if err != nil {
		http.Error(w, err.Error(), getErrorCode(err))
		return
	}

	res := responses.New(user, "user created")
	if err := jsonutil.Encode(w, r, http.StatusOK, *res); err != nil {
		slog.Error("json encode: %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoded, err := jsonutil.Decode[LoginRequest](r)
	if err != nil {
		slog.Error("json decode: %w", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(decoded.Email, decoded.Password)
	if err != nil {
		http.Error(w, err.Error(), getErrorCode(err))
		return
	}

	if err := jsonutil.Encode(w, r, http.StatusOK, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("json encode: %w", err)
		return
	}
}

func (h *UserHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	decoded, err := jsonutil.Decode[RefreshRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.Refresh(decoded.UserId, decoded.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), getErrorCode(err))
		return
	}

	if err := jsonutil.Encode(w, r, http.StatusOK, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("json encode: %w", err)
		return
	}
}
