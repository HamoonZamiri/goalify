package svcerror

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest     = errors.New("bad request")
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
	return http.StatusInternalServerError
}
