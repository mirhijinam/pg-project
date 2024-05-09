package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

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
		log.Fatalf("failed to open database: %v", err)
	}

	// check if db is alive
	if err = db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	fmt.Println("success")
	return db
}
