package api

import (
	"context"
	"net/http"

	"github.com/mirhijinam/pg-project/internal/model"
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

func New(s Service) *CommandHandler {
	return &CommandHandler{
		CommandService: s,
	}

}

type Handler interface {
	CreateCmd() http.HandlerFunc
	StopCmd() http.HandlerFunc
	GetCmd() http.HandlerFunc
	GetCmdList() http.HandlerFunc
}
