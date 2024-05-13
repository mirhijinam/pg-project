package api

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/mirhijinam/pg-project/internal/config"
	"github.com/mirhijinam/pg-project/internal/model"
	"github.com/mirhijinam/pg-project/internal/pkg/db"
	"github.com/mirhijinam/pg-project/internal/pkg/logger"
	"github.com/mirhijinam/pg-project/internal/repository"
	"github.com/mirhijinam/pg-project/internal/service"
)

var (
	contextTimeoutSec   = 600
	countWorkersDefault = 1
	prefixStopUrl       = "/stop_cmd/"
	prefixGetListUrl    = "/cmd_list/"
)

type Handler interface {
	CreateCmd() http.HandlerFunc
	StopCmd() http.HandlerFunc
	GetCmd() http.HandlerFunc
	GetCmdList() http.HandlerFunc
}

type Service interface {
	CreateCommand(context.Context, *model.Command, bool) error
	StopCommand(context.Context, int) error
	GetCommand(context.Context, int) (model.Command, error)
	GetCommandList(context.Context) ([]model.Command, error)
}

type CommandHandler struct {
	CommandService Service
}

func New(ctx context.Context, dbCfg config.DBConfig) *http.Handler {
	countWorkersMax, err := strconv.Atoi(os.Getenv("MAXCOUNT"))
	if err != nil {
		countWorkersMax = countWorkersDefault
	}

	dbConn := db.MustOpenDB(ctx, dbCfg)

	ch := &CommandHandler{
		CommandService: service.New(
			repository.New(dbConn),
			service.NewPool(countWorkersMax),
		),
	}

	log := logger.New(logger.SetupLogger(os.Getenv("LOGENV")))

	mux := http.NewServeMux()
	mux.HandleFunc("POST /create_cmd", ch.CreateCmd())
	mux.HandleFunc("POST /stop_cmd/{id}", ch.StopCmd())
	mux.HandleFunc("GET /cmd_list/", ch.GetCmdList())
	mux.HandleFunc("GET /cmd_list/{id}", ch.GetCmd())

	logMux := log(mux)

	return &logMux
}
