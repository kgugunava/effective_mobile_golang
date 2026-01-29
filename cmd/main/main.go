package main

import (
	"log/slog"
	"os"

	"github.com/kgugunava/effective_mobile_golang/internal/app"
)


func main() {
	logger := initLogger()
	app := app.NewApp(logger)
	app.Router.Run(app.Cfg.ServerAddress)
}

func initLogger() *slog.Logger {
	isDev := os.Getenv("LOGGER_MODE") == "development"

	var logger *slog.Logger
	if isDev {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	slog.SetDefault(logger)

	return logger
}