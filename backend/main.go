package main

import (
	"goalify/db"
	gh "goalify/goals/handler"
	gSrv "goalify/goals/service"
	gs "goalify/goals/stores"
	"goalify/middleware"
	uh "goalify/users/handler"
	usrSrv "goalify/users/service"
	us "goalify/users/stores"
	"goalify/utils/events"
	"goalify/utils/stacktrace"
	"log/slog"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func addRoute(mux *http.ServeMux, method, path string, handler http.HandlerFunc, mwChain middleware.Middleware) {
	mux.Handle(method+" "+path, mwChain(http.HandlerFunc(handler)))
}

func NewServer(userHandler *uh.UserHandler, goalHandler *gh.GoalHandler) http.Handler {
	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		Debug:            true,
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
	CorsChain := middleware.CreateChain(c.Handler)
	AuthChain := middleware.CreateChain(c.Handler, middleware.AuthenticatedOnly)

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

	addRoute(mux, "POST", "/api/goals/categories", goalHandler.HandleCreateGoalCategory, AuthChain)
	addRoute(mux, "GET", "/api/goals/categories", goalHandler.HandleGetGoalCategoriesByUserId, AuthChain)
	addRoute(mux, "GET", "/api/goals/categories/{categoryId}", goalHandler.HandleGetGoalCategoryById, AuthChain)
	addRoute(mux, "PUT", "/api/goals/categories/{categoryId}", goalHandler.HandleUpdateGoalCategoryById, AuthChain)
	addRoute(mux, "DELETE", "/api/goals/categories/{categoryId}", goalHandler.HandleDeleteGoalCategoryById, AuthChain)

	return mux
}

func Run() error {
	db, _ := db.New("goalify")

	// logs for stack trace implementing stacktrace.TraceLogger
	goalDomainLogger := stacktrace.NewDomainStackTraceLogger("Goals")

	eventManager := events.NewEventManager()

	userStore := us.NewUserStore(db)
	userService := usrSrv.NewUserService(userStore, eventManager)
	userHandler := uh.NewUserHandler(userService)

	goalStore := gs.NewGoalStore(db)
	goalCategoryStore := gs.NewGoalCategoryStore(db)
	goalService := gSrv.NewGoalService(goalStore, goalCategoryStore,
		goalDomainLogger, eventManager)
	goalHandler := gh.NewGoalHandler(goalService, goalDomainLogger)

	srv := NewServer(userHandler, goalHandler)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: srv,
	}

	var err error = nil
	slog.Info("Listening on 8080")

	if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("ListenAndServe: ", "err", err)
	}

	return err
}

func main() {
	if err := Run(); err != nil {
		slog.Error("run: ", "err", err)
		os.Exit(1)
	}
}
