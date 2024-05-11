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

func main() {
	ctx := context.Background()

	dbCfg, err := config.GetDBConfig()
	if err != nil {
		slog.Error("failed to parse db config", "error", err.Error())
		os.Exit(1)
	}

	mux := api.New(ctx, dbCfg)

	srvCfg := config.GetServerConfig()

	err = run(mux, srvCfg)
	if err != nil {
		slog.Error("failed to run the server", "error", err.Error())
		os.Exit(2)
	}
}

func run(mux *http.ServeMux, srvCfg config.ServerConfig) error {
	srv := http.Server{
		Addr:    srvCfg.HTTPPort,
		Handler: mux,
	}

	serveChan := make(chan error, 1)
	go func() {
		serveChan <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		slog.Info("Shutting down the server")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	case err := <-serveChan:
		return err
	}
}
