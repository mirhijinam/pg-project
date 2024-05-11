package api

import (
	"context"
	"net/http"
)

func (h *CommandHandler) GetCmdList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		cmdList, err := h.CommandService.GetCommandList(ctx)
		if err != nil {
			h.serverErrorResponse(w, r, err)
			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"cmdList": cmdList}, nil)
		if err != nil {
			h.serverErrorResponse(w, r, err)
		}
	}
}
