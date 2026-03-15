package _postgres

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"assignment-03/pkg/modules"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"time"
	"os"
)

type Dialect struct {
	DB *sqlx.DB
}

func NewPGDialect(ctx context.Context, cfg *modules.PostgreConfig) *Dialect {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	var db *sqlx.DB
	var err error

	for i := 1; i <= 20; i++ {
		db, err = sqlx.ConnectContext(ctx, "postgres", dsn)
		if err == nil {
			err = db.PingContext(ctx)
		}
		if err == nil {
			break
		}
		fmt.Printf("Database not ready yet (attempt %d/20): %v\n", i, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		panic(err)
	}

	if os.Getenv("DISABLE_MIGRATIONS") != "true" {
		autoMigrate(cfg)
	}
	return &Dialect{DB: db}
}

func autoMigrate(cfg *modules.PostgreConfig) {
	sourceURL := "file://database/migrations"
	databaseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		panic(err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
}