/*
Package middleware includes REST API middleware handlers. For example this package includes middleware to log requests and responses, authenticate users and apply CORS headers.
*/
package middleware

import (
	"goalify/config"
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
	configService := config.GetConfig()

	c := cors.New(cors.Options{
		AllowedOrigins: configService.AllowedOrigins,
		AllowedMethods: []string{
			http.MethodPatch, http.MethodPost, http.MethodPut, http.MethodGet,
			http.MethodDelete, http.MethodOptions,
		},
		AllowCredentials: true,
		Debug:            configService.IsDevelopment(),
		AllowedHeaders:   []string{"*"},
	})

	return MiddleWareChains{
		CorsChain:           CreateChain(Logging, c.Handler),
		AuthChain:           CreateChain(Logging, c.Handler, AuthenticatedOnly(userService)),
		QueryTokenAuthChain: CreateChain(Logging, c.Handler, QueryTokenAuth(userService)),
	}
}
