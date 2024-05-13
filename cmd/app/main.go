package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mirhijinam/pg-project/internal/api"
	"github.com/mirhijinam/pg-project/internal/config"
)

// @title pg-project
// @version 1.0
// @description API to run commands. A command is a bash script. Allow to run commands in parallel

// @host localhost:7070
// @BasePath /
func main() {
	ctx := context.Background()

	dbCfg, err := config.GetDBConfig()
	if err != nil {
		slog.Error("failed to parse db config", "error", err.Error())
		os.Exit(1)
	}

	srvCfg, err := config.GetServerConfig()
	if err != nil {
		slog.Error("failed to parse server config", "error", err.Error())
		os.Exit(1)
	}

	handler := api.New(ctx, dbCfg)

	err = run(*handler, srvCfg)
	if err != nil {
		slog.Error("failed to run the server", "error", err.Error())
		os.Exit(2)
	}
}

func run(handler http.Handler, srvCfg config.ServerConfig) error {
	srv := http.Server{
		Handler:      handler,
		Addr:         srvCfg.HTTPPort,
		WriteTimeout: srvCfg.Timeout,
		ReadTimeout:  srvCfg.Timeout,
		IdleTimeout:  srvCfg.IdleTimeout,
	}

	serveChan := make(chan error, 1)
	go func() {
		serveChan <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	case err := <-serveChan:
		return err
	}
}
