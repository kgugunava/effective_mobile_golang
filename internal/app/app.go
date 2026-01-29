package app

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kgugunava/effective_mobile_golang/internal/adapters/postgres"
	"github.com/kgugunava/effective_mobile_golang/internal/api"
	"github.com/kgugunava/effective_mobile_golang/internal/api/handlers"
	"github.com/kgugunava/effective_mobile_golang/internal/config"
	"github.com/kgugunava/effective_mobile_golang/internal/service"
)

type App struct {
	Cfg config.Config
	Router *gin.Engine
	Logger *slog.Logger
}

func NewApp(db *pgxpool.Pool, logger *slog.Logger, cfg config.Config) *App {
	app := &App{
		Cfg: cfg,
		Logger: logger,
	}


	// db := postgres.NewPostgres()

    // if err := db.ConnectToPostgresMainDatabase(app.Cfg); err != nil {
    //     panic(err)
    // }

	// if err := db.CreateDatabase(app.Cfg); err != nil {
    //     panic(err)
    // }
    
    // if err := db.ConnectToDatabase(app.Cfg); err != nil {
    //     panic(err)
    // }
    
    // if err := db.CreateDatabaseTables(); err != nil {
    //     panic(err)
    // }
    
    // app.DB = &db

	subscriptionsRepository := postgres.NewSubscriptionRepository(db, logger)

	subscriptionsService := service.NewSubscriptionService(subscriptionsRepository, logger)

	apiSubscriptions := handlers.NewSubscriptionAPI(subscriptionsService, logger)

    
    app.Router = api.NewRouter(*apiSubscriptions)
    
    return app
}