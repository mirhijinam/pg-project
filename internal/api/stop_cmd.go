package api

import (
	"net/http"
)

type deleteCmdRequest struct {
	CmdId int `json:"cmd_id"`
}

func (h *CommandHandler) StopCmd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		inp := deleteCmdRequest{}
		err := readJSONBody(w, r, &inp)
		if err != nil {
			h.badRequestResponse(w, r, err)
			return
		}
		err = h.CommandService.DeleteCommand(inp.CmdId)
		if err != nil {
			h.serverErrorResponse(w, r, err)
			return
		}

		headers := make(http.Header)

		err = writeJSONBody(w, http.StatusOK, envelope{"CmdId": inp.CmdId}, headers)
		if err != nil {
			h.serverErrorResponse(w, r, err)
			return
		}
	}
}
