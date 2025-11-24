// Package jsonutil is used for pairing decoding json with validation
package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Validator interface {
	Valid() (problems map[string]string)
}

func Encode[T any](w http.ResponseWriter, _ *http.Request, status int, data T) error {
	// credit:
	// https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/#handle-decodingencoding-in-one-place
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		var typeErr *json.UnmarshalTypeError
		var syntaxErr *json.SyntaxError
		switch {
		case errors.As(err, &typeErr):
			return v, fmt.Errorf(
				"invalid type for field '%s', expected %s",
				typeErr.Field,
				typeErr.Type,
			)
		case errors.As(err, &syntaxErr):
			return v, fmt.Errorf("invalid JSON syntax at position %d", syntaxErr.Offset)
		default:
			return v, fmt.Errorf("decode json: %w", err)
		}
	}
	return v, nil
}

func DecodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	v, err := Decode[T](r)
	if err != nil {
		return v, nil, err
	}

	if problems := v.Valid(); len(problems) > 0 {
		return v, problems, fmt.Errorf("request validation issue")
	}

	return v, nil, nil
}
