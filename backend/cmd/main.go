package main

import (
	"goalify/cmd/app"
	"log/slog"
	"os"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("run: ", "err", err)
		os.Exit(1)
	}
}
