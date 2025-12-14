package main

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/ucok-man/tcsa/internal/data"
	"github.com/ucok-man/tcsa/internal/tlog"
)

const VERSION = "1.0.0"

type application struct {
	config Config
	logger *tlog.Logger
	models data.Models
	wg     sync.WaitGroup
}

func main() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger := tlog.Must(tlog.NewProduction())
	if cfg.Env != "production" {
		logger = tlog.Must(tlog.NewDevelopment())
	}
	defer logger.Sync()

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatalj(tlog.JSON{"message": "failed connecting to database", "err": err})
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	if err != nil {
		logger.Fatalj(tlog.JSON{"message": "server has error occured", "error": err})
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
