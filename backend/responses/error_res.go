package responses

import (
	"goalify/jsonutil"
	"net/http"
)

type APIError struct {
	Errors  map[string]string `json:"errors"`
	Message string            `json:"message"`
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
