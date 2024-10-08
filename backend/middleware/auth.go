package middleware

import (
	"errors"
	"goalify/users/service"
	"goalify/utils/responses"
	"log/slog"
	"net/http"
	"strings"
)

func AuthenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authstr := r.Header.Get("Authorization")
			if authstr == "" {
				responses.SendAPIError(w, r, http.StatusUnauthorized, "empty header: unauthorized request", nil)
				return
			}

			split := strings.Split(authstr, " ")
			if len(split) != 2 {
				responses.SendAPIError(w, r, http.StatusUnauthorized, "malformed header: unauthorized request", nil)
				return
			}

			token := split[1]
			id, err := service.VerifyToken(token)
			if err != nil {
				slog.Error("middleware.AuthenticatedOnly: service.VerifyToken:", "err", err)
				responses.SendAPIError(w, r, http.StatusUnauthorized, "invalid access token: unauthorized request", nil)
				return
			}

			r.Header.Set("user_id", id)
			next.ServeHTTP(w, r)
		},
	)
}

func QueryTokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get("token")
			if token == "" {
				responses.SendAPIError(w, r, http.StatusUnauthorized, "empty query: unauthorized request", nil)
				return
			}

			id, err := service.VerifyToken(token)
			if err != nil {
				slog.Error("middleware.QueryTokenAuth: service.VerifyToken:", "err", err)
				responses.SendAPIError(w, r, http.StatusUnauthorized, "invalid access token: unauthorized request", nil)
				return
			}

			r.Header.Set("user_id", id)
			next.ServeHTTP(w, r)
		})
}

func GetIdFromHeader(r *http.Request) (string, error) {
	id := r.Header.Get("user_id")
	if id == "" {
		return "", errors.New("user_id missing in header")
	}
	return id, nil
}
