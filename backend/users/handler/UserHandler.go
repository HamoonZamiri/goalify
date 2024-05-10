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
	decoded, problems, err := jsonutil.DecodeValid[SignupRequest](r)
	if len(problems) > 0 {
		apiErr := responses.NewAPIError(err.Error(), problems)
		jsonutil.Encode(w, r, http.StatusBadRequest, apiErr)
		return
	}

	user, err := h.userService.SignUp(decoded.Email, decoded.Password)
	if err != nil {
		status := getErrorCode(err)
		apiErr := responses.NewAPIError(err.Error(), nil)
		err := jsonutil.Encode(w, r, status, apiErr)
		if err != nil {
			slog.Error("json encode: ", "err", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	res := responses.New(user, "user created")
	if err := jsonutil.Encode(w, r, http.StatusOK, *res); err != nil {
		slog.Error("json encode: ", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoded, err := jsonutil.Decode[LoginRequest](r)
	if err != nil {
		slog.Error("json decode: ", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.userService.Login(decoded.Email, decoded.Password)
	if err != nil {
		http.Error(w, err.Error(), getErrorCode(err))
		return
	}

	if err := jsonutil.Encode(w, r, http.StatusOK, user); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		slog.Error("json encode: ", "err", err)
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
		slog.Error("json encode: ", "err", err)
		return
	}
}
