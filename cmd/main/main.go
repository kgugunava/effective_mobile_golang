package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/kgugunava/effective_mobile_golang/internal/app"
	"github.com/kgugunava/effective_mobile_golang/internal/config"
)

func main() {
	logger := initLogger()
	cfg := config.NewConfig()
	cfg.InitConfig()

	dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
		cfg.SslMode,
	)

	runMigrations(logger, dbURL)

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	application := app.NewApp(db, logger, cfg)
	logger.Info("starting server", slog.String("address", cfg.ServerAddress))
	application.Router.Run(cfg.ServerAddress)
}

func runMigrations(logger *slog.Logger, dbURL string) {
	logger.Info("starting database migrations", slog.String("db_url", dbURL))

	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		logger.Error("failed to create migrate instance", 
			slog.String("db_url", dbURL),
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Error("failed to apply migrations",
			slog.String("db_url", dbURL),
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	version, _, _ := m.Version()
	logger.Info("migrations applied successfully",
		slog.Int64("version", int64(version)),
	)
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