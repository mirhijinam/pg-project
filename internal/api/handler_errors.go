package api

import (
	"fmt"
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

func (h *CommandHandler) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	ans := "answer"
	msg := fmt.Sprintf("error! the %s method is not supported for this resource", r.Method)
	env := envelope{ans: msg}

	h.responseCreator(w, r, http.StatusMethodNotAllowed, env)
}

func (h *CommandHandler) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	ans := "answer"
	msg := "error! unable to update the record due to an edit conflict, please try again"
	env := envelope{ans: msg}

	h.responseCreator(w, r, http.StatusConflict, env)
}
