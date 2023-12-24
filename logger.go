package main

import (
	"log/slog"
	"os"
)

func initLogger() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	slog.SetDefault(logger)
	slog.Info("logger has started")
}
