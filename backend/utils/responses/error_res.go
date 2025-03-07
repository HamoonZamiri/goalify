package responses

import (
	"errors"
	"fmt"
	"goalify/utils/jsonutil"
	"log/slog"
	"net/http"
)

type APIError struct {
	Errors  map[string]string `json:"errors,omitempty"`
	Message string            `json:"message"`
}

var (
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInternalServer = errors.New("internal server error")
	ErrNotFound       = errors.New("not found")
)

func GetErrorCode(err error) int {
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest
	}
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, ErrUnauthorized) {
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}

func NewAPIError(message string, errors map[string]string) *APIError {
	return &APIError{
		Message: message,
		Errors:  errors,
	}
}

func SendAPIError(w http.ResponseWriter, r *http.Request, status int, err string, errors map[string]string) {
	jsonutil.Encode(w, r, status, NewAPIError(err, errors))
}

func HandleDecodeError(w http.ResponseWriter, r *http.Request, problems map[string]string, err error) {
	// unprocessable entity -> 422
	var res error
	if len(problems) > 0 {
		res = fmt.Errorf("%w: %w", ErrBadRequest, err)
		SendAPIError(w, r, http.StatusUnprocessableEntity, res.Error(), problems)
		return
	}

	// error decoding JSON -> 400
	res = fmt.Errorf("%w: %w", ErrBadRequest, err)
	SendAPIError(w, r, http.StatusBadRequest, res.Error(), nil)
	slog.Error("HandleDecodeError: jsonutil.DecodeValid: ", "err", res)
}

func SendInternalServerError(w http.ResponseWriter, r *http.Request) {
	SendAPIError(w, r, http.StatusInternalServerError, ErrInternalServer.Error(), nil)
}
