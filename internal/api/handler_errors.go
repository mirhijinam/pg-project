package api

import (
	"log/slog"
	"net/http"
)

func (h *CommandHandler) responseLog(_ *http.Request, err error) {
	slog.Error(err.Error())
}

func (h *CommandHandler) responseCreator(w http.ResponseWriter, r *http.Request, status int, env map[string]interface{}) {
	err := writeJSONBody(w, status, env, nil)
	if err != nil {
		h.responseLog(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *CommandHandler) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	h.responseLog(r, err)
	ans := "answer"
	msg := "error! the server encountered a problem and could not process your request"
	env := envelope{ans: msg}

	h.responseCreator(w, r, http.StatusNotFound, env)
}

func (h *CommandHandler) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ans := "answer"
	msg := "error! " + err.Error()
	env := envelope{ans: msg}

	h.responseCreator(w, r, http.StatusBadRequest, env)
}

func (h *CommandHandler) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	ans := "answer"
	msg := "error! the requested resource could not be found"
	env := envelope{ans: msg}

	h.responseCreator(w, r, http.StatusNotFound, env)
}
