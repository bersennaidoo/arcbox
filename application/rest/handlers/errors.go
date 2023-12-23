package handlers

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (h *SnipHandler) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	h.log.Error(trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *SnipHandler) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (h *SnipHandler) notFound(w http.ResponseWriter) {
	h.clientError(w, http.StatusNotFound)
}
