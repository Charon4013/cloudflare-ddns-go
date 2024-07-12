package util

import (
	"log/slog"
	"os"
)

var logger *slog.Logger
var logFile *os.File

func InitLogger() {

	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	handler := slog.NewTextHandler(logFile, nil)

	logger = slog.New(handler)
}

func GetLogger() *slog.Logger {
	return logger
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
