package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/mirhijinam/pg-project/internal/api"
	"github.com/mirhijinam/pg-project/internal/config"
	"github.com/mirhijinam/pg-project/internal/pkg/db"
	"github.com/mirhijinam/pg-project/internal/repository"
	"github.com/mirhijinam/pg-project/internal/service"
)

func main() {
	dbCfg, err := config.GetDBConfig()
	if err != nil {
		slog.Error("failed to parse db config", "error", err.Error())
		os.Exit(1)
	}
	ctx := context.Background()
	dbConn := db.MustOpenDB(ctx, dbCfg)

	if err != nil {
		slog.Error("failed to init db", "error", err.Error())
		os.Exit(1)
	}

	maxWorkers := config.GetLoggerConfig().MaxCount
	pool := service.NewPool(maxWorkers)

	srvCfg := config.GetServerConfig()

	cr := repository.New(dbConn)
	cs := service.New(cr, pool)
	ch := api.New(cs)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /create_cmd", ch.CreateCmd())
	mux.HandleFunc("DELETE /delete_cmd", ch.DeleteCmd())
	mux.HandleFunc("GET /cmd_list", ch.GetCmdList())
	mux.HandleFunc("GET /cmd_list/{id}", ch.GetCmd())

	srv := http.Server{
		Addr:    srvCfg.HTTPPort,
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		slog.Error("failed to start server", "error", err.Error())
		os.Exit(2)
	}
}
