package api

import (
	"context"
	"log/slog"
	"net/http"
)

// GetCmdList handles getting of all commands.
// @Description Getting a list of all commands.
// @Tags Cmd
// @Produce application/json
// @Success 200 {object} envelope "List of commands were gotten successfully"
// @Failure 500 {object} envelope "Internal server error"
// @Router /cmd_list [get]
func (h *CommandHandler) GetCmdList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		cmdList, err := h.CommandService.GetCommandList(ctx)
		if err != nil {
			slog.Error("failed to get a command list", "error", err.Error())
			h.serverErrorResponse(w, r, err)
			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"list of executed commands": cmdList}, nil)
		if err != nil {
			slog.Error("failed to write json response", "error", err.Error())
			h.serverErrorResponse(w, r, err)
			return
		}
	}
}
