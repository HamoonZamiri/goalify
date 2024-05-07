package handler

import (
	"goalify/jsonutil"
	"goalify/responses"
	"goalify/users/service"
	"log/slog"
	"net/http"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest = SignupRequest

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) HandleSignup(w http.ResponseWriter, r *http.Request) {
	decoded, err := jsonutil.Decode[SignupRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.SignUp(decoded.Email, decoded.Password)
	if err != nil {
		slog.Debug("error: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := responses.New(user, "user created")
	if err := jsonutil.Encode(w, r, http.StatusOK, *res); err != nil {
		slog.Error("json encode: %w", err)
		return
	}
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoded, err := jsonutil.Decode[LoginRequest](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(decoded.Email, decoded.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if err := jsonutil.Encode(w, r, http.StatusOK, user); err != nil {
		slog.Error("json encode: %w", err)
		return
	}
}
