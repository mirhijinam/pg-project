package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/mirhijinam/pg-project/internal/config"
)

func MustOpenDB(ctx context.Context, cfg config.DBConfig) *sql.DB {
	// construct the dsn
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		cfg.PgUser, cfg.PgPass, cfg.PgHost, cfg.PgPort, cfg.PgDb,
	)

	// open db
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("failed to open database", "error", err.Error())
		os.Exit(1)
	}

	// check if db is alive
	if err = db.PingContext(ctx); err != nil {
		slog.Error("failed to ping database", "error", err.Error())
		os.Exit(1)
	}

	return db
}
