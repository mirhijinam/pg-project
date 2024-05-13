package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mirhijinam/pg-project/internal/model"
)

type GetCmdResponse struct {
}

// GetCmd handles getting of a command by id
// @Description Getting a command by id
// @Tags Cmd
// @Produce application/json
// @Param id path int true "id of the command to get"
// @Success 200 {object} envelope "Command was gotten successfully"
// @Failure 400 {object} envelope "Bad request if the id is incorrect"
// @Failure 404 {object} envelope "Not found if there is no command with the given id"
// @Failure 500 {object} envelope "Internal server error"
// @Router /cmd_list/{id} [get]
func (h *CommandHandler) GetCmd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(contextTimeoutSec)*time.Second)
		defer cancel()

		inpIdStr, ok := strings.CutPrefix(r.URL.Path, prefixGetListUrl)
		if !ok {
			err := errors.New("id retrieval from the url")
			slog.Error("failed to parse url", "error", err.Error())
			h.badRequestResponse(w, r, err)
			return
		}

		inpId, err := strconv.Atoi(inpIdStr)
		if err != nil {
			slog.Error("failed to define an id from the url", "error", err.Error())
			h.serverErrorResponse(w, r, err)
			return
		}

		cmd, err := h.CommandService.GetCommand(ctx, inpId)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrRecordNotFound):
				slog.Error("failed to found a command with such id", "error", err.Error())
				h.notFoundResponse(w, r)
			default:
				slog.Error("failed to get a command by such id", "error", err.Error())
				h.serverErrorResponse(w, r, err)
			}
			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"requested command": cmd}, nil)
		if err != nil {
			slog.Error("failed to write json response", "error", err.Error())
			h.serverErrorResponse(w, r, err)
			return
		}
	}
}
