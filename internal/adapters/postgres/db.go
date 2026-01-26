package postgres

import (
	"context"
	"log"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

    "github.com/kgugunava/effective_mobile_golang/internal/config"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres() Postgres{
	return Postgres{
		Pool: &pgxpool.Pool{},
	}
}

func (p *Postgres) ConnectToPostgresMainDatabase(cfg config.Config) error { // для подключения к бд постгреса для создания нашей бд
    dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, "postgres", cfg.SslMode)
    
    fmt.Println("Connecting to DB with:", dbUrl)
    
    newPostgresPool, err := pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        log.Fatal("Eror while connecting to postgres database\n", err)
        return err
    }
    p.Pool = newPostgresPool
    return nil
}

func (p *Postgres) ConnectToDatabase(cfg config.Config) error { // для подключения к нужной бд
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.SslMode)
    
    newPostgresPool, err := pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        log.Fatal("Error while connecting to database\n", err)
        return err
    }
    p.Pool = newPostgresPool
    return nil
} 

func (p *Postgres) CreateDatabase(cfg config.Config) error {
	var dbExists bool
    err := p.Pool.QueryRow(context.Background(), 
        "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DbName).Scan(&dbExists)
    if err != nil {
        log.Fatal(err)
    }
    
    if !dbExists {
        _, err := p.Pool.Exec(context.Background(), 
            fmt.Sprintf("CREATE DATABASE %s", cfg.DbName))
        if err != nil {
            log.Fatal(err)
            return err
        }

        p.Pool.Close()
    
        dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", 
            cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.SslMode)
        
        p.Pool, err = pgxpool.New(context.Background(), dbUrl)
        if err != nil {
            log.Fatal("Eror while creating database\n", err)
            return err
        }
    }
    return nil
}

func (p *Postgres) CreateDatabaseTables() error {
    _, err := p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS subscriptions (
            subscription_id UUID PRIMARY KEY,
            service_name VARCHAR NOT NULL,
			price INTEGER NOT NULL CHECK (price > 0),
			user_id UUID NOT NULL,
			start_date DATE,
			end_date DATE
        );
    `)
	if err != nil {
		log.Fatal("Error creating teams table\n", err)
		return err
	}

	return nil
}