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
	"log/slog"
	"net/http"
	"os"
)

func NewServer(userHandler *uh.UserHandler, goalHandler *gh.GoalHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello\n"))
	})

	// users domain
	mux.HandleFunc("POST /api/users/signup", userHandler.HandleSignup)
	mux.HandleFunc("POST /api/users/login", userHandler.HandleLogin)
	mux.HandleFunc("POST /api/users/refresh", userHandler.HandleRefresh)
	mux.Handle(http.MethodPut+" /api/users", middleware.AuthenticatedOnly(userHandler.HandleUpdateUserById))

	// goals domain
	mux.Handle("POST /api/goals/create", middleware.AuthenticatedOnly(goalHandler.HandleCreateGoal))
	mux.Handle("POST /api/goals/categories", middleware.AuthenticatedOnly(goalHandler.HandleCreateGoalCategory))
	mux.Handle("GET /api/goals/categories", middleware.AuthenticatedOnly(goalHandler.HandleGetGoalCategoriesByUserId))
	mux.Handle("GET /api/goals/categories/{categoryId}", middleware.AuthenticatedOnly(goalHandler.HandleGetGoalCategoryById))
	mux.Handle("PUT /api/goals/categories/{categoryId}", middleware.AuthenticatedOnly(goalHandler.HandleUpdateGoalCategoryById))
	mux.Handle("DELETE /api/goals/categories/{categoryId}", middleware.AuthenticatedOnly(goalHandler.HandleDeleteGoalCategoryById))

	mux.Handle("POST /api/goals", middleware.AuthenticatedOnly(goalHandler.HandleCreateGoal))
	mux.Handle("PUT /api/goals/{goalId}", middleware.AuthenticatedOnly(goalHandler.HandleUpdateGoalById))
	return mux
}

func Run() error {
	db, _ := db.New("goalify")

	userStore := us.NewUserStore(db)
	userService := usrSrv.NewUserService(userStore)
	userHandler := uh.NewUserHandler(userService)

	goalStore := gs.NewGoalStore(db)
	goalCategoryStore := gs.NewGoalCategoryStore(db)
	goalService := gSrv.NewGoalService(goalStore, goalCategoryStore)
	goalHandler := gh.NewGoalHandler(goalService)

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
