package middleware

import (
	"errors"
	"goalify/users/service"
	"log/slog"
	"net/http"
	"strings"
)

func AuthenticatedOnly(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authstr := r.Header.Get("Authorization")
			slog.Debug("header", "Authorization", authstr)
			if authstr == "" {
				http.Error(w, "empty header unauthorized request", http.StatusUnauthorized)
				return
			}

			split := strings.Split(authstr, " ")
			if len(split) != 2 {
				http.Error(w, "malformed header: unauthorized request", http.StatusUnauthorized)
				return
			}

			token := split[1]
			id, err := service.VerifyToken(token)
			if err != nil {
				slog.Error("verify token: ", "err", err)
				http.Error(w, "invalid access token: unauthorized request", http.StatusUnauthorized)
				return
			}

			r.Header.Set("user_id", id)
			next(w, r)
		},
	)
}

func GetIdFromHeader(r *http.Request) (string, error) {
	id := r.Header.Get("user_id")
	if id == "" {
		return "", errors.New("user_id missing in header")
	}
	return id, nil
}
