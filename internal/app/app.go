package app

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/kgugunava/effective_mobile_golang/internal/adapters/postgres"
	"github.com/kgugunava/effective_mobile_golang/internal/api"
	"github.com/kgugunava/effective_mobile_golang/internal/api/handlers"
	"github.com/kgugunava/effective_mobile_golang/internal/config"
	"github.com/kgugunava/effective_mobile_golang/internal/service"
)

type App struct {
	Cfg config.Config
	Router *gin.Engine
	DB *postgres.Postgres
	Logger *slog.Logger
}

func NewApp(logger *slog.Logger) *App {
	app := &App{
		Cfg: config.NewConfig(),
	}
	app.Cfg.InitConfig()

	db := postgres.NewPostgres()

    if err := db.ConnectToPostgresMainDatabase(app.Cfg); err != nil {
        panic(err)
    }

	if err := db.CreateDatabase(app.Cfg); err != nil {
        panic(err)
    }
    
    if err := db.ConnectToDatabase(app.Cfg); err != nil {
        panic(err)
    }
    
    if err := db.CreateDatabaseTables(); err != nil {
        panic(err)
    }
    
    app.DB = &db

	subscriptionsRepository := postgres.NewSubscriptionRepository(app.DB.Pool, app.Logger)

	subscriptionsService := service.NewSubscriptionService(subscriptionsRepository, app.Logger)

	apiSubscriptions := handlers.NewSubscriptionAPI(subscriptionsService, app.Logger)

    
    app.Router = api.NewRouter(*apiSubscriptions)
    
    return app
}