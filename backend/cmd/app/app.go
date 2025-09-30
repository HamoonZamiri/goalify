package app

import (
	"context"
	"goalify/config"
	"goalify/db"
	sqlcdb "goalify/db/generated"
	gh "goalify/goals/handler"
	gSrv "goalify/goals/service"
	gs "goalify/goals/stores"
	"goalify/middleware"
	"goalify/routes"
	uh "goalify/users/handler"
	usrSrv "goalify/users/service"
	us "goalify/users/stores"
	"goalify/utils/events"
	"goalify/utils/stacktrace"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func NewServer(userHandler *uh.UserHandler, goalHandler *gh.GoalHandler,
	em *events.EventManager, userService usrSrv.UserService,
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
	var pgxPool *pgxpool.Pool
	configService := config.GetConfig()
	currEnv := configService.Env

	var connStr string
	if configService.Env == config.LocalTest {
		connStr = configService.GetTestDBConnectionString()
	} else {
		connStr = configService.GetDBConnectionString()
	}
	pgxPool, err = db.NewPgxPoolWithConnString(context.Background(), connStr)
	if err != nil {
		panic(err)
	}
	if pgxPool == nil {
		panic("pgx pool is nil")
	}

	// using goose run migrations from db/migrations
	if currEnv == config.LocalTest {
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		migrationDir := filepath.Join(basepath, "../../db/migrations")
		err = goose.UpContext(context.Background(), stdlib.OpenDBFromPool(pgxPool), migrationDir)
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic(err)
	}

	// logs for stack trace implementing stacktrace.TraceLogger
	goalDomainLogger := stacktrace.NewDomainStackTraceLogger("Goals")

	eventManager := events.NewEventManager()

	// Create sqlc queries
	queries := sqlcdb.New(pgxPool)

	userStore := us.NewUserStore(queries)
	userService := usrSrv.NewUserService(userStore, eventManager)
	userHandler := uh.NewUserHandler(userService)

	goalStore := gs.NewGoalStore(queries)
	goalCategoryStore := gs.NewGoalCategoryStore(queries)
	goalService := gSrv.NewGoalService(
		goalStore,
		goalCategoryStore,
		goalDomainLogger,
		eventManager,
	)
	goalHandler := gh.NewGoalHandler(goalService, goalDomainLogger)

	srv := NewServer(userHandler, goalHandler, eventManager, userService)
	port := configService.Port
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

