package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/mirhijinam/pg-project/internal/model"
)

// StopCmd handles stopping a command by its id
// @Summary Stop a command
// @Description Stop a command by its id
// @Tags Cmd
// @Produce application/json
// @Param id path int true "id of the command to stop"
// @Success 200 {object} envelope "Command stopped successfully, returns the id of the stopped command"
// @Failure 400 {object} envelope "Bad request if the id is incorrect or conversion error"
// @Failure 404 {object} envelope "Not found if there is no command with the given id"
// @Failure 500 {object} envelope "Internal server error"
// @Router /stop_cmd/{id} [post]
func (h *CommandHandler) StopCmd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inpIdStr, ok := strings.CutPrefix(r.URL.Path, stopUrlPrefix)
		if !ok {
			h.serverErrorResponse(w, r, model.ErrRecordNotFound)
		}

		inpId, err := strconv.Atoi(inpIdStr)
		if err != nil {
			h.serverErrorResponse(w, r, err)
			slog.Error("occasionally caught atoi's error. probability = 0.00001%")
			return
		}
		slog.Info("StopCmd:", slog.Int("id", inpId))

		err = h.CommandService.DeleteCommand(inpId)
		if err != nil {
			h.serverErrorResponse(w, r, err)
			return
		}

		headers := make(http.Header)

		err = writeJSON(w, http.StatusOK, envelope{"CmdId": inpId}, headers)
		if err != nil {
			h.serverErrorResponse(w, r, err)
			return
		}
	}
}
