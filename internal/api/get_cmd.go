package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/mirhijinam/pg-project/internal/model"
)

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
		ctx := context.Background()
		inpIdStr, ok := strings.CutPrefix(r.URL.Path, getListUrlPrefix)
		if !ok {
			h.serverErrorResponse(w, r, model.ErrRecordNotFound)
		}

		inpId, err := strconv.Atoi(inpIdStr)
		if err != nil {
			h.serverErrorResponse(w, r, err)
			slog.Error("occasionally caught atoi's error. probability = 0.00001%")
			return
		}
		slog.Info("GetCmd:", slog.Int("id", inpId))

		cmd, err := h.CommandService.GetCommand(ctx, inpId)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrRecordNotFound):
				h.notFoundResponse(w, r)
			default:
				h.serverErrorResponse(w, r, err)
			}
			return
		}

		err = writeJSON(w, http.StatusOK, envelope{"cmd": cmd}, nil)
		if err != nil {
			h.serverErrorResponse(w, r, err)
		}
	}
}
