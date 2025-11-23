package middleware

import (
	"errors"
	"fmt"
	"goalify/internal/responses"
	"goalify/internal/users/service"
	"log/slog"
	"net/http"
	"strings"
)

func AuthenticatedOnly(userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				authstr := r.Header.Get("Authorization")
				authError := fmt.Errorf(
					"%w: could not authenticate request",
					responses.ErrUnauthorized,
				)
				if authstr == "" {
					responses.SendAPIError(w, r, http.StatusUnauthorized, authError.Error(), nil)
					return
				}

				split := strings.Split(authstr, " ")
				if len(split) != 2 {
					responses.SendAPIError(w, r, http.StatusUnauthorized, authError.Error(), nil)
					return
				}

				token := split[1]
				id, err := userService.VerifyToken(token)
				if errors.Is(err, responses.ErrUnauthorized) {
					responses.SendAPIError(w, r, http.StatusUnauthorized, authError.Error(), nil)
					return
				}
				if err != nil {
					responses.SendAPIError(
						w,
						r,
						http.StatusInternalServerError,
						authError.Error(),
						nil,
					)
					return
				}

				r.Header.Set("user_id", id)
				next.ServeHTTP(w, r)
			},
		)
	}
}

func QueryTokenAuth(userService service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				token := r.URL.Query().Get("token")
				if token == "" {
					responses.SendAPIError(
						w,
						r,
						http.StatusUnauthorized,
						"empty query: unauthorized request",
						nil,
					)
					return
				}

				id, err := userService.VerifyToken(token)
				if err != nil {
					slog.Error("middleware.QueryTokenAuth: service.VerifyToken:", "err", err)
					responses.SendAPIError(
						w,
						r,
						http.StatusUnauthorized,
						"invalid access token: unauthorized request",
						nil,
					)
					return
				}

				r.Header.Set("user_id", id)
				next.ServeHTTP(w, r)
			})
	}
}

func GetIDFromHeader(r *http.Request) (string, error) {
	id := r.Header.Get("user_id")
	if id == "" {
		return "", errors.New("user_id missing in header")
	}
	return id, nil
}
