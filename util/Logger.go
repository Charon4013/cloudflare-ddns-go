package util

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger() {

	handler := slog.NewTextHandler(os.Stdout, nil)

	logger = slog.New(handler)
}

func GetLogger() *slog.Logger {
	return logger
}
