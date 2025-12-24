// Package app is the main entry point for the REST server
package app

import (
	"context"
	"fmt"
	"goalify/internal/config"
	"goalify/internal/db"
	"goalify/internal/events"
	"goalify/internal/middleware"
	"goalify/internal/routes"
	"goalify/pkg/stacktrace"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	sqlcdb "goalify/internal/db/generated"

	gh "goalify/internal/goals/handler"
	gSrv "goalify/internal/goals/service"
	gs "goalify/internal/goals/stores"

	uh "goalify/internal/users/handler"
	usrSrv "goalify/internal/users/service"
	us "goalify/internal/users/stores"
)

func NewServer(userHandler *uh.UserHandler, goalHandler *gh.GoalHandler,
	em *events.EventManager, userService usrSrv.UserService,
) http.Handler {
	mux := http.NewServeMux()
	mw := middleware.SetupMiddleware(userService)
	routes.AddRoutes(mux, userHandler, goalHandler, em, mw)
	return mux
}

func Run(ctx context.Context) error {
	// Setup signal handling
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

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
	pgxPool, err = db.NewPgxPoolWithConnString(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to create pgx pool: %w", err)
	}
	if pgxPool == nil {
		return fmt.Errorf("pgx pool is nil")
	}

	// using goose run migrations from db/migrations
	if currEnv == config.LocalTest {
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		migrationDir := filepath.Join(basepath, "../../internal/db/migrations")
		err = goose.UpContext(ctx, stdlib.OpenDBFromPool(pgxPool), migrationDir)
		if err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
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

	// Start server in goroutine
	go func() {
		slog.Info("Listening on " + port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	// Wait for interrupt signal and gracefully shutdown
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		slog.Info("Shutdown signal received, gracefully shutting down...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}
