package svcerror

import "errors"

var (
	ErrBadRequest     = errors.New("bad request")
	ErrInternalServer = errors.New("internal server error")
	ErrNotFound       = errors.New("not found")
)
