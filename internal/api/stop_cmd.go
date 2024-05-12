package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type StopCmdResponse struct {
}

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

		err = h.CommandService.DeleteCommand(inpId)
		if err != nil {
			slog.Error("failed to stop a command by such id", "error", err.Error())
			h.serverErrorResponse(w, r, err)
			return
		}

		headers := make(http.Header)

		err = writeJSON(w, http.StatusOK, envelope{"id of stopped command": inpId}, headers)
		if err != nil {
			slog.Error("failed to write json response", "error", err.Error())
			h.serverErrorResponse(w, r, err)
			return
		}
	}
}
