package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/mirhijinam/pg-project/internal/model"
)

type Service interface {
	CreateCommand(*model.Command, bool) error
	ExecCommand(context.Context, model.Command) error
}

type CommandHandler struct {
	CommandService Service
	Logger         *slog.Logger
}

func New(s Service) *CommandHandler {
	return &CommandHandler{
		CommandService: s,
	}

}

type Handler interface {
	GetCmd() http.HandlerFunc
	GetCmdList() http.HandlerFunc
	ExecCmd() http.HandlerFunc
	StopCmd() http.HandlerFunc
}
