package main

import (
	"context"
	"goalify/config"
	gh "goalify/goals/handler"
	gSrv "goalify/goals/service"
	gs "goalify/goals/stores"
	"goalify/middleware"
	"goalify/routes"
	"goalify/testsetup"
	uh "goalify/users/handler"
	"goalify/users/service"
	usrSrv "goalify/users/service"
	us "goalify/users/stores"
	"goalify/utils/events"
	"goalify/utils/options"
	"goalify/utils/stacktrace"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
)

func addRoute(mux *http.ServeMux, method, path string, handler http.HandlerFunc, mwChain middleware.Middleware) {
	mux.Handle(method+" "+path, mwChain(http.HandlerFunc(handler)))
}

func NewServer(userHandler *uh.UserHandler, goalHandler *gh.GoalHandler,
	configService *config.ConfigService, em *events.EventManager, userService service.UserService,
) http.Handler {
	mux := http.NewServeMux()
	mw := middleware.SetupMiddleware(userService)
	routes.AddRoutes(mux, userHandler, goalHandler, em, mw)
	return mux
}

func Run() error {
	// instantiate config service
	var err error
	var dbInstance *sqlx.DB
	configService := config.NewConfigService(options.None[string]())
	dbInstance, err = testsetup.GetDbInstance()

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
	// global slog default logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	container, err := testsetup.GetPgContainer()
	if err != nil {
		slog.Error("Failed to create container: ", "err", err)
		os.Exit(1)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	done := make(chan bool)

	go func() {
		sig := <-sigChan
		slog.Info("Received signal: ", "sig", sig)
		cancel()

		slog.Info("Cleaning up resources")
		if err := container.Terminate(context.Background()); err != nil {
			slog.Error("Failed to terminate container: ", "err", err)
		}

		close(done)
	}()

	go func() {
		if err := Run(); err != nil {
			slog.Error("run: ", "err", err)
			os.Exit(1)
		}
	}()

	select {
	case <-done:
		slog.Info("Cleanup complete")
		os.Exit(0)
		// case <-time.After(10 * time.Second):
		// 	slog.Error("timeout")
		// 	os.Exit(1)
	}
}
