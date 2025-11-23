// Package routes is the routing layer for the API defining routes and the methods available
package routes

import (
	"goalify/internal/events"
	"goalify/internal/middleware"
	"net/http"

	gh "goalify/internal/goals/handler"

	uh "goalify/internal/users/handler"
)

func addRoute(mux *http.ServeMux, method, path string,
	handler http.HandlerFunc, mwChain middleware.Middleware,
) {
	mux.Handle(method+" "+path, mwChain(http.HandlerFunc(handler)))
}

func AddRoutes(
	mux *http.ServeMux,
	userHandler *uh.UserHandler,
	goalHandler *gh.GoalHandler,
	em *events.EventManager,
	mw middleware.MiddleWareChains,
) http.Handler {
	mux.Handle(
		"GET /health",
		mw.CorsChain(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write([]byte("Hello\n"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})),
	)

	// users domain
	addRoute(mux, http.MethodPost, "/api/users/signup", userHandler.HandleSignup, mw.CorsChain)
	addRoute(mux, http.MethodPost, "/api/users/login", userHandler.HandleLogin, mw.CorsChain)
	addRoute(mux, http.MethodPost, "/api/users/refresh", userHandler.HandleRefresh, mw.CorsChain)
	addRoute(mux, http.MethodPut, "/api/users", userHandler.HandleUpdateUserByID, mw.AuthChain)
	addRoute(mux, http.MethodGet, "/api/levels/{levelId}", userHandler.GetLevelByID, mw.AuthChain)

	// goals domain
	addRoute(mux, http.MethodPost, "/api/goals", goalHandler.HandleCreateGoal, mw.AuthChain)
	addRoute(
		mux,
		http.MethodPut,
		"/api/goals/{goalId}",
		goalHandler.HandleUpdateGoalByID,
		mw.AuthChain,
	)
	addRoute(
		mux,
		http.MethodDelete,
		"/api/goals/{goalId}",
		goalHandler.HandleDeleteGoalByID,
		mw.AuthChain,
	)

	addRoute(
		mux,
		http.MethodPost,
		"/api/goals/categories",
		goalHandler.HandleCreateGoalCategory,
		mw.AuthChain,
	)
	addRoute(
		mux,
		http.MethodGet,
		"/api/goals/categories",
		goalHandler.HandleGetGoalCategoriesByUserID,
		mw.AuthChain,
	)
	addRoute(
		mux,
		http.MethodGet,
		"/api/goals/categories/{categoryId}",
		goalHandler.HandleGetGoalCategoryByID,
		mw.AuthChain,
	)
	addRoute(
		mux,
		http.MethodPut,
		"/api/goals/categories/{categoryId}",
		goalHandler.HandleUpdateGoalCategoryByID,
		mw.AuthChain,
	)
	addRoute(
		mux,
		http.MethodDelete,
		"/api/goals/categories/{categoryId}",
		goalHandler.HandleDeleteGoalCategoryByID,
		mw.AuthChain,
	)

	// need options method available on all endpoints for CORS
	mux.Handle(
		"OPTIONS /api/",
		mw.CorsChain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		})),
	)

	// Server Sent Events endpoint
	addRoute(mux, http.MethodOptions, "/api/events", nil, mw.CorsChain)
	addRoute(mux, http.MethodGet, "/api/events", em.SSEHandler, mw.QueryTokenAuthChain)

	// websocket endpoint
	addRoute(mux, http.MethodGet, "/api/ws", em.WebSocketHandler, mw.QueryTokenAuthChain)
	return mux
}
