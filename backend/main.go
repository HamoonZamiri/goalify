package main

import (
	"goalify/config"
	"goalify/db"
	gh "goalify/goals/handler"
	gSrv "goalify/goals/service"
	gs "goalify/goals/stores"
	"goalify/middleware"
	"goalify/routes"
	uh "goalify/users/handler"
	usrSrv "goalify/users/service"
	us "goalify/users/stores"
	"goalify/utils/events"
	"goalify/utils/options"
	"goalify/utils/stacktrace"
	"log/slog"
	"net/http"
	"os"
)

func addRoute(mux *http.ServeMux, method, path string, handler http.HandlerFunc, mwChain middleware.Middleware) {
	mux.Handle(method+" "+path, mwChain(http.HandlerFunc(handler)))
}

func NewServer(userHandler *uh.UserHandler, goalHandler *gh.GoalHandler,
	configService *config.ConfigService,
) http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(mux, userHandler, goalHandler, *configService)
	return mux
}

func Run() error {
	// instantiate config service
	configService := config.NewConfigService(options.None[string]())

	var dbName string
	if configService.MustGetEnv("ENV") == "test" {
		dbName = configService.MustGetEnv("TEST_DB_NAME")
	} else {
		dbName = configService.MustGetEnv("DB_NAME")
	}
	db, _ := db.New(dbName, configService.MustGetEnv("DB_PASSWORD"),
		configService.MustGetEnv("DB_NAME"))

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

	srv := NewServer(userHandler, goalHandler, configService)
	port := configService.MustGetEnv("PORT")
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: srv,
	}

	var err error = nil
	slog.Info("Listening on " + port)

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
