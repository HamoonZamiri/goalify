package main

import (
	"context"
	"fmt"
	"goalify/cmd/app"
	"os"
)

func main() {
	ctx := context.Background()
	if err := app.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
