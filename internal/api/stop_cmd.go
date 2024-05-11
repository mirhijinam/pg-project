package api

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/mirhijinam/pg-project/internal/model"
)

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
