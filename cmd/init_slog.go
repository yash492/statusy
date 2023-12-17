package main

import (
	"log/slog"
	"os"
)

func setupSlog() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	slog.SetDefault(logger)
	slog.Info("logger has started")
}
