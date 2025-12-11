package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"
)

type application struct {
	config Config
	logger *zap.Logger
	// models data.Models
	wg sync.WaitGroup
}

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := zap.Must(zap.NewProduction())
	if cfg.Env != "production" {
		logger = zap.Must(zap.NewDevelopment())
	}
	defer logger.Sync()

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal("Failed to open database connection", zap.Error(err))

	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
	}

	err = app.serve()
	if err != nil {
		logger.Fatal("Server has error occured", zap.Error(err))
	}
}

func openDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Database.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConn)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConn)
	db.SetConnMaxIdleTime(cfg.Database.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
