package jsonutil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Validator interface {
	Valid() (problems map[string]string)
}

// credit: https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/#handle-decodingencoding-in-one-place
func Encode[T any](w http.ResponseWriter, _ *http.Request, status int, data T) error {
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
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func DecodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	v, err := Decode[T](r)
	if err != nil {
		return v, nil, err
	}

	if problems := v.Valid(); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}

	return v, nil, nil
}
