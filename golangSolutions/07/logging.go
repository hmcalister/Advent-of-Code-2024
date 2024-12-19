package main

import (
	"io"
	"log/slog"
	"os"
)

const (
	LOG_FILE_PATH = "log"
)

func SetLogging(debugFlag bool) *os.File {
	logFileHandler, err := os.Create(LOG_FILE_PATH)
	if err != nil {
		slog.Error("cannot open log file", "error", err)
		os.Exit(1)
	}

	var slogHandler slog.Handler
	if debugFlag {
		multiwriter := io.MultiWriter(os.Stdout, logFileHandler)
		slogHandler = slog.NewTextHandler(multiwriter, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		slogHandler = slog.NewJSONHandler(logFileHandler, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	slog.SetDefault(slog.New(
		slogHandler,
	))

	return logFileHandler
}
