package main

import (
	"context"
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
	"path/filepath"
	"runtime"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

func addRoute(mux *http.ServeMux, method, path string, handler http.HandlerFunc, mwChain middleware.Middleware) {
	mux.Handle(method+" "+path, mwChain(http.HandlerFunc(handler)))
}

func NewServer(userHandler *uh.UserHandler, goalHandler *gh.GoalHandler,
	configService *config.ConfigService, em *events.EventManager, userService usrSrv.UserService,
) http.Handler {
	mux := http.NewServeMux()
	mw := middleware.SetupMiddleware(userService)
	routes.AddRoutes(mux, userHandler, goalHandler, em, mw)
	return mux
}

func Run() error {
	// global slog default logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// instantiate config service
	var err error
	var dbInstance *sqlx.DB
	configService := config.NewConfigService(options.None[string]())
	currEnv := configService.MustGetEnv(config.ENV)
	if currEnv == "test" {
		dbInstance, err = db.NewWithConnString(configService.MustGetEnv(config.TEST_DB_CONN_STRING))
	} else {
		dbInstance, err = db.New(
			configService.MustGetEnv(config.DB_NAME),
			configService.MustGetEnv(config.DB_USER),
			configService.MustGetEnv(config.DB_PASSWORD))
	}

	// using goose run migrations from db/migrations
	if currEnv == "dev" || currEnv == "test" {
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		migrationDir := filepath.Join(basepath, "./db/migrations")
		err = goose.UpContext(context.Background(), dbInstance.DB, migrationDir)
		if err != nil {
			panic(err)
		}
	}

	if dbInstance == nil {
		panic("db instance is nil")
	}
	if err != nil {
		panic(err)
	}

	// logs for stack trace implementing stacktrace.TraceLogger
	goalDomainLogger := stacktrace.NewDomainStackTraceLogger("Goals")

	eventManager := events.NewEventManager()

	userStore := us.NewUserStore(dbInstance)
	userService := usrSrv.NewUserService(userStore, eventManager)
	userHandler := uh.NewUserHandler(userService)

	goalStore := gs.NewGoalStore(dbInstance)
	goalCategoryStore := gs.NewGoalCategoryStore(dbInstance)
	goalService := gSrv.NewGoalService(
		goalStore,
		goalCategoryStore,
		goalDomainLogger,
		eventManager,
	)
	goalHandler := gh.NewGoalHandler(goalService, goalDomainLogger)

	srv := NewServer(userHandler, goalHandler, configService, eventManager, userService)
	port := configService.MustGetEnv(config.PORT)
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: srv,
	}

	err = nil
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
