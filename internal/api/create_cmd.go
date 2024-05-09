package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mirhijinam/pg-project/internal/model"
)

type createCmdRequest struct {
	CmdRaw    string `json:"cmd_raw"`
	IsLongCmd bool   `json:"is_long_cmd"`
}

func (h *CommandHandler) CreateCmd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inp := createCmdRequest{}
		err := readJSONBody(w, r, &inp)
		if err != nil {
			h.badRequestResponse(w, r, err)
			return
		}

		if len(inp.CmdRaw) == 0 {
			h.badRequestResponse(w, r, errors.New("empty command text"))
			return
		}

		name, isSudo := getSudoNameCommand(inp.CmdRaw)

		if !isSudo {
			cmd := model.Command{
				Name:      name,
				Raw:       inp.CmdRaw,
				CreatedAt: time.Now(),
			}

			err = h.CommandService.CreateCommand(&cmd, inp.IsLongCmd)
			if err != nil {
				h.serverErrorResponse(w, r, err)
				return
			}

			headers := make(http.Header)
			headers.Set("Location", fmt.Sprintf("/cmd_list/%d", cmd.Id))

			err = writeJSONBody(w, http.StatusCreated, envelope{"cmd": cmd}, headers)
			if err != nil {
				h.serverErrorResponse(w, r, err)
				return
			}
		}
	}
}
