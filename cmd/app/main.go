package main

import (
	"log/slog"
	"os"

	"application-design/internal"
)

func main() {
	if err := internal.Run(); err != nil {
		slog.Error("Error while running the server", "err", err)
		os.Exit(1)
	}
}
