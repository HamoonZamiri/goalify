package routes

import (
	"goalify/config"
	gh "goalify/goals/handler"
	"goalify/middleware"
	uh "goalify/users/handler"
	"goalify/utils/events"
	"net/http"

	"github.com/rs/cors"
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
	configService config.ConfigService,
	em *events.EventManager,
) http.Handler {
	// // var corsDebug bool
	// if configService.MustGetEnv("ENV") == "dev" {
	// 	corsDebug = true
	// } else {
	// 	corsDebug = false
	// }

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{
			http.MethodPatch,
			http.MethodPost,
			http.MethodPut,
			http.MethodGet,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowCredentials: true,
		Debug:            false,
		AllowedHeaders:   []string{"*"},
	})
	CorsChain := middleware.CreateChain(middleware.Logging, c.Handler)
	AuthChain := middleware.CreateChain(middleware.Logging, c.Handler, middleware.AuthenticatedOnly)
	QueryTokenAuthChain := middleware.CreateChain(middleware.QueryTokenAuth)

	mux.Handle("GET /health", CorsChain(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello\n"))
	})))

	// users domain
	addRoute(mux, "POST", "/api/users/signup", userHandler.HandleSignup, CorsChain)
	addRoute(mux, "POST", "/api/users/login", userHandler.HandleLogin, CorsChain)
	addRoute(mux, "POST", "/api/users/refresh", userHandler.HandleRefresh, CorsChain)
	addRoute(mux, "PUT", "/api/users", userHandler.HandleUpdateUserById, AuthChain)

	// goals domain
	addRoute(mux, "POST", "/api/goals", goalHandler.HandleCreateGoal, AuthChain)
	addRoute(mux, "PUT", "/api/goals/{goalId}", goalHandler.HandleUpdateGoalById, AuthChain)
	addRoute(mux, http.MethodDelete, "/api/goals/{goalId}", goalHandler.HandleDeleteGoalById, AuthChain)

	addRoute(mux, "POST", "/api/goals/categories", goalHandler.HandleCreateGoalCategory, AuthChain)
	addRoute(mux, "GET", "/api/goals/categories", goalHandler.HandleGetGoalCategoriesByUserId, AuthChain)
	addRoute(mux, "GET", "/api/goals/categories/{categoryId}", goalHandler.HandleGetGoalCategoryById, AuthChain)
	addRoute(mux, "PUT", "/api/goals/categories/{categoryId}", goalHandler.HandleUpdateGoalCategoryById, AuthChain)
	addRoute(mux, "DELETE", "/api/goals/categories/{categoryId}", goalHandler.HandleDeleteGoalCategoryById, AuthChain)

	// need options method available on all endpoints for CORS
	addRoute(mux, "OPTIONS", "/api/*", nil, CorsChain)

	// Server Sent Events endpoint
	addRoute(mux, http.MethodOptions, "/api/events", nil, CorsChain)
	addRoute(mux, http.MethodGet, "/api/events", em.SSEHandler, QueryTokenAuthChain)
	return mux
}
