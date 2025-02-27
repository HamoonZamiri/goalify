package middleware

import (
	"goalify/users/service"
	"net/http"

	"github.com/rs/cors"
)

type (
	Middleware       func(http.Handler) http.Handler
	MiddleWareChains struct {
		CorsChain           Middleware
		AuthChain           Middleware
		QueryTokenAuthChain Middleware
	}
)

func CreateChain(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			mw := mws[i]
			next = mw(next)
		}
		return next
	}
}

func SetupMiddleware(userService service.UserService) MiddleWareChains {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{
			http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodGet,
			http.MethodDelete, http.MethodOptions,
		},
		AllowCredentials: true,
		Debug:            false,
		AllowedHeaders:   []string{"*"},
	})

	return MiddleWareChains{
		CorsChain:           CreateChain(Logging, c.Handler),
		AuthChain:           CreateChain(Logging, c.Handler, AuthenticatedOnly(userService)),
		QueryTokenAuthChain: CreateChain(Logging, c.Handler, QueryTokenAuth(userService)),
	}
}
