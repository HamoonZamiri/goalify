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

	"github.com/jmoiron/sqlx"
)

func addRoute(mux *http.ServeMux, method, path string, handler http.HandlerFunc, mwChain middleware.Middleware) {
	mux.Handle(method+" "+path, mwChain(http.HandlerFunc(handler)))
}

func NewServer(userHandler *uh.UserHandler, goalHandler *gh.GoalHandler,
	configService *config.ConfigService, em *events.EventManager,
) http.Handler {
	mux := http.NewServeMux()
	routes.AddRoutes(mux, userHandler, goalHandler, *configService, em)
	return mux
}

func Run() error {
	// instantiate config service
	configService := config.NewConfigService(options.None[string]())
	var dbInstance *sqlx.DB

	if configService.MustGetEnv("ENV") == "test" {
		dbInstance, _ = db.NewWithConnString(configService.MustGetEnv(config.TEST_DB_CONN_STRING))
	} else {
		dbInstance, _ = db.New(configService.MustGetEnv(config.DB_NAME), configService.MustGetEnv(config.DB_PASSWORD),
			configService.MustGetEnv(config.DB_NAME))
	}
	if dbInstance == nil {
		panic("db instance is nil")
	}

	// logs for stack trace implementing stacktrace.TraceLogger
	goalDomainLogger := stacktrace.NewDomainStackTraceLogger("Goals")

	eventManager := events.NewEventManager()

	userStore := us.NewUserStore(dbInstance)
	userService := usrSrv.NewUserService(userStore, eventManager)
	userHandler := uh.NewUserHandler(userService)

	goalStore := gs.NewGoalStore(dbInstance)
	goalCategoryStore := gs.NewGoalCategoryStore(dbInstance)
	goalService := gSrv.NewGoalService(goalStore, goalCategoryStore,
		goalDomainLogger, eventManager)
	goalHandler := gh.NewGoalHandler(goalService, goalDomainLogger)

	srv := NewServer(userHandler, goalHandler, configService, eventManager)
	port := configService.MustGetEnv(config.PORT)
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
