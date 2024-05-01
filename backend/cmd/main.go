package main

import (
	"log/slog"
	"net/http"
	"os"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello\n"))
	})

	return mux
}

func run() error {
	srv := NewServer()
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
	if err := run(); err != nil {
		slog.Error("run: %v", err)
		os.Exit(1)
	}
}
