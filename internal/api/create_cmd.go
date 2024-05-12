package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/mirhijinam/pg-project/internal/model"
)

type createCmdRequest struct {
	CmdRaw    string `json:"cmd_raw"`
	IsLongCmd bool   `json:"is_long_cmd"`
}

// CreateCmd handles the creation and execution of a command
// @Description Create and execute a command
// @Tags Cmd
// @Produce application/json
// @Param request body createCmdRequest true "Request Create Command"
// @Header 201 {string} Locationswag init "/cmd_list/{id}" "Location of the created command"
// @Success 201 {object} envelope "Command created successfully"
// @Failure 400 {object} envelope "Bad request if the input data is incorrect"
// @Failure 403 {object} envelope "Forbidden if attempting a sudo command without admin rights"
// @Failure 500 {object} envelope "Internal server error if command creation fails"
// @Router /create_cmd [post]
func (h *CommandHandler) CreateCmd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inp := createCmdRequest{}

		err := readJSONBody(w, r, &inp)
		if err != nil {
			slog.Error("failed to read json request", "error", err.Error())
			h.badRequestResponse(w, r, err)
			return
		}

		if len(inp.CmdRaw) == 0 {
			err := errors.New("empty command text")
			slog.Error("failed to create the provided command", "error", err.Error())
			h.badRequestResponse(w, r, err)
			return
		}

		name, isSudo := getSudoNameCommand(inp.CmdRaw)
		admin := isAdmin(r.Header.Get("token"))
		cmd := model.Command{
			Name:      name,
			Raw:       inp.CmdRaw,
			CreatedAt: time.Now(),
		}

		if isSudo && !admin {

			h.forbiddenAccessResponse(w, r)
			return
		} else {
			err = h.CommandService.CreateCommand(&cmd, inp.IsLongCmd)
			if err != nil {
				slog.Error("failed to create the provided command", "error", err.Error())
				h.serverErrorResponse(w, r, err)
				return
			}

			headers := make(http.Header)
			headers.Set("Location", fmt.Sprintf("/cmd_list/%d", cmd.Id))

			err = writeJSON(w, http.StatusCreated, envelope{"created command": cmd}, headers)
			if err != nil {
				slog.Error("failed to write json response", "error", err.Error())
				h.serverErrorResponse(w, r, err)
				return
			}
		}
	}
}
