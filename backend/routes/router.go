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
	QueryTokenAuthChain := middleware.CreateChain(c.Handler, middleware.QueryTokenAuth)

	mux.Handle("GET /health", CorsChain(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello\n"))
	})))

	// users domain
	addRoute(mux, http.MethodPost, "/api/users/signup", userHandler.HandleSignup, CorsChain)
	addRoute(mux, http.MethodPost, "/api/users/login", userHandler.HandleLogin, CorsChain)
	addRoute(mux, http.MethodPost, "/api/users/refresh", userHandler.HandleRefresh, CorsChain)
	addRoute(mux, http.MethodPut, "/api/users", userHandler.HandleUpdateUserById, AuthChain)
	addRoute(mux, http.MethodGet, "/api/levels/{levelId}", userHandler.GetLevelById, AuthChain)

	// goals domain
	addRoute(mux, http.MethodPost, "/api/goals", goalHandler.HandleCreateGoal, AuthChain)
	addRoute(mux, http.MethodPut, "/api/goals/{goalId}", goalHandler.HandleUpdateGoalById, AuthChain)
	addRoute(mux, http.MethodDelete, "/api/goals/{goalId}", goalHandler.HandleDeleteGoalById, AuthChain)

	addRoute(mux, http.MethodPost, "/api/goals/categories", goalHandler.HandleCreateGoalCategory, AuthChain)
	addRoute(mux, http.MethodGet, "/api/goals/categories", goalHandler.HandleGetGoalCategoriesByUserId, AuthChain)
	addRoute(mux, http.MethodGet, "/api/goals/categories/{categoryId}", goalHandler.HandleGetGoalCategoryById, AuthChain)
	addRoute(mux, http.MethodPut, "/api/goals/categories/{categoryId}", goalHandler.HandleUpdateGoalCategoryById, AuthChain)
	addRoute(mux, http.MethodDelete, "/api/goals/categories/{categoryId}", goalHandler.HandleDeleteGoalCategoryById, AuthChain)

	// need options method available on all endpoints for CORS
	addRoute(mux, http.MethodOptions, "/api/{endpoints...}", nil, middleware.CreateChain(c.Handler))

	// Server Sent Events endpoint
	addRoute(mux, http.MethodOptions, "/api/events", nil, CorsChain)
	addRoute(mux, http.MethodGet, "/api/events", em.SSEHandler, QueryTokenAuthChain)

	// websocket endpoint
	addRoute(mux, http.MethodGet, "/api/ws", em.WebSocketHandler, QueryTokenAuthChain)
	return mux
}
