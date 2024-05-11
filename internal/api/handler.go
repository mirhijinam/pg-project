package api

import (
	"context"
	"net/http"

	"github.com/mirhijinam/pg-project/internal/config"
	"github.com/mirhijinam/pg-project/internal/model"
	"github.com/mirhijinam/pg-project/internal/pkg/db"
	"github.com/mirhijinam/pg-project/internal/repository"
	"github.com/mirhijinam/pg-project/internal/service"
)

var (
	stopUrlPrefix    = "/stop_cmd/"
	getListUrlPrefix = "/cmd_list/"
)

type Service interface {
	CreateCommand(*model.Command, bool) error
	DeleteCommand(int) error
	GetCommand(context.Context, int) (model.Command, error)
	GetCommandList(context.Context) ([]model.Command, error)
}

type CommandHandler struct {
	CommandService Service
}

func New(ctx context.Context, dbCfg config.DBConfig) *http.ServeMux {
	maxWorkers := config.GetLoggerConfig().MaxCount
	dbConn := db.MustOpenDB(ctx, dbCfg)

	ch := &CommandHandler{
		service.New(
			repository.New(dbConn),
			service.NewPool(maxWorkers),
		),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /create_cmd", ch.CreateCmd())
	mux.HandleFunc("POST /stop_cmd", ch.StopCmd())
	mux.HandleFunc("GET /cmd_list/", ch.GetCmdList())
	mux.HandleFunc("GET /cmd_list/{id}", ch.GetCmd())

	return mux
}

type Handler interface {
	CreateCmd() http.HandlerFunc
	StopCmd() http.HandlerFunc
	GetCmd() http.HandlerFunc
	GetCmdList() http.HandlerFunc
}
