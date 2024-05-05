package main

import (
	"goalify/db"
	"goalify/users/handler"
	"goalify/users/service"
	"goalify/users/stores"
	"log/slog"
	"net/http"
	"os"
)

func NewServer(userHandler *handler.UserHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello\n"))
	})

	// users domain
	mux.HandleFunc("POST /api/users/signup", userHandler.HandleSignup)

	return mux
}

func Run() error {
	db, _ := db.New("goalify")

	var userStore stores.UserStore = stores.NewUserStore(db)
	var userService service.UserService = service.NewUserService(userStore)
	userHandler := handler.NewUserHandler(userService)

	srv := NewServer(userHandler)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: srv,
	}

	var err error = nil
	slog.Info("Listening on 8080")

	if err = httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("ListenAndServe: %v", err)
	}

	return err
}

func main() {
	if err := Run(); err != nil {
		slog.Error("run: %v", err)
		os.Exit(1)
	}
}
