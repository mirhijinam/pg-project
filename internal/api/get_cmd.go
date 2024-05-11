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

var (
	urlPrefix = "/cmd_list/"
)

func (h *CommandHandler) GetCmd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		inpIdStr, ok := strings.CutPrefix(r.URL.Path, urlPrefix)
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
